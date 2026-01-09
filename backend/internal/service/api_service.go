package service

import (
	"backend/internal/cleaners"
	"backend/internal/detector"
	"backend/internal/models"
	"context"
	"errors"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// maxPathsToCollect limits the number of file paths collected during the reviewing phase
const maxPathsToCollect = 500

// workers determines the concurrency level of operations, based on the CPU count
var workers = runtime.NumCPU()

// LoadCleanerMap transforms the flat list of cleaners into a nested map structure.
//
// It loads all definitions via cleaners_util and organizes them for O(1) lookup
// during the analysis phase.
// Returns a map keyed by [CleanerID][OptionID] containing the list of Actions.
func LoadCleanerMap(ctx context.Context) (map[string]map[string][]models.Action, error) {
	allCleaners, err := cleaners.LoadAllCleaners(ctx)
	if err != nil {
		slog.Error("Error loading all cleaners: %v\n", err)
		return nil, err
	}

	cleanerMap := make(map[string]map[string][]models.Action)
	for _, cleaner := range allCleaners {
		cleanerMap[cleaner.ID] = make(map[string][]models.Action)
		for _, option := range cleaner.Options {
			cleanerMap[cleaner.ID][option.ID] = option.Actions
		}
	}
	return cleanerMap, nil
}

// AnalyzeRequests serves as the entry point for processing a batch of cleanup requests.
//
// It orchestrates the analysis by:
// 1. Validating that the requested CleanerID and OptionID exist.
// 2. Spinning up concurrent workers (limited by the 'workers' global) to process requests.
// 3. Aggregating the results (Size, FileCount) into a single response.
func AnalyzeRequests(ctx context.Context,requests []models.CleanRequest,
	cleanerMap map[string]map[string][]models.Action) (*models.AnalyzeResponse, error) {
	response := &models.AnalyzeResponse{
		Items: make([]models.AnalyzeItem, 0),
	}

	semaphore := make(chan struct{}, workers)
	var wg sync.WaitGroup
	resultsChan := make(chan models.AnalyzeItem, len(requests))

	for _, request := range requests {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		actions, ok := cleanerMap[request.CleanerID][request.OptionID]
		if !ok {
			continue
		}

		wg.Add(1)
		select {
		case semaphore <- struct{}{}:
			go func(request models.CleanRequest, actions []models.Action) {
				defer wg.Done() // decrease the counter when the goroutine completes
				defer func() { <-semaphore }() // clear the semaphore slot when done

				item, err := AnalyzeActions(ctx, request, actions)
				if err != nil {
					return
				}

				select {
				case resultsChan <- item:
				case <-ctx.Done():
					return
				}
			} (request, actions)
		case <-ctx.Done():
			wg.Done()
			return nil, ctx.Err()
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(semaphore)
	}()

	for item := range resultsChan {
		response.Items = append(response.Items, item)
		response.TotalSize += item.Size
		response.TotalFiles += item.FileCount
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return response, nil
}

// AnalyzeActions processes the specific actions (paths/globs) associated with a single cleaner option.
//
// It checks OS compatibility for each action and executes them concurrently using
// the same worker-pool pattern as AnalyzeRequests.
func AnalyzeActions(ctx context.Context, request models.CleanRequest, actions []models.Action) (models.AnalyzeItem, error) {
	var size uint64 = 0
	var fileCount uint64 = 0
	var foundPaths []string

	semaphore := make(chan struct{}, workers)
	var wg sync.WaitGroup
	resultChan := make(chan models.ActionResult, len(actions))

	for _, action := range actions {
		if ctx.Err() != nil {
			return models.AnalyzeItem{}, ctx.Err()
		}

		if !detector.IsOSSupported(action.OS) {
			continue
		}

		wg.Add(1)
		select {
		case semaphore <- struct{}{}:
			go func(action models.Action) {
				defer wg.Done()
				defer func() { <-semaphore }()

				actionSize, actionCount, actionPaths := ProcessAction(ctx, action)

				result := models.ActionResult{
					Size:      actionSize,
					FileCount: actionCount,
					Paths:     actionPaths,
				}

				select {
				case resultChan <- result:
				case <-ctx.Done():
					return
				}
			}(action)
		case <-ctx.Done():
			wg.Done()
			return models.AnalyzeItem{}, ctx.Err()
		}
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(semaphore)
	}()

	for result := range resultChan {
		size += result.Size
		fileCount += result.FileCount
		foundPaths = append(foundPaths, result.Paths...)
	}

	if ctx.Err() != nil {
		return models.AnalyzeItem{}, ctx.Err()
	}

	return models.AnalyzeItem{
		CleanerID: request.CleanerID,
		OptionID:  request.OptionID,
		Size:      size,
		FileCount: fileCount,
		Paths:     foundPaths,
	}, nil
}

// ProcessAction acts as a router to determine the correct file discovery strategy.
//
// It expands environment variables in paths (e.g., %APPDATA%) and selects between:
// - Globbing (if "*" is present or explicitly set)
// - Recursive walking ("walk.files")
// - Single file verification
func ProcessAction(ctx context.Context, action models.Action) (uint64, uint64, []string) {
	if ctx.Err() != nil {
		return 0, 0, nil
	}

	searchPath := detector.ExpandPath(action.Path)

	if action.Search == "glob" || strings.Contains(searchPath, "*") {
		return ProcessGlobAction(ctx, searchPath)
	} else if action.Search == "walk.files" {
		return ProcessWalkAction(ctx, searchPath)
	}

	return ProcessFileAction(searchPath)
}

// ProcessGlobAction handles file discovery using standard filesystem glob patterns.
//
// It performs a concurrent stat() on all matches found by filepath.Glob.
// Thread-safe: Uses a mutex to safely update the size, count, and path slice.
func ProcessGlobAction(ctx context.Context, searchPath string) (uint64, uint64, []string) {
	var size uint64 = 0
	var fileCount uint64 = 0
	var paths []string

	matches, err := filepath.Glob(searchPath)
	if err != nil {
		log.Printf("Error in glob %s: %v\n", searchPath, err)
		return 0, 0, nil
	}

	slog.Info("Processing glob", "path", searchPath, "matches", len(matches))

	var wg sync.WaitGroup
	var mutex sync.Mutex
	semaphore := make(chan struct{}, workers)

	cancelled := false
	for _, match := range matches {
		select {
		case <-ctx.Done():
			cancelled = true
		default:
		}

		if cancelled {
			break
		}

		wg.Add(1)

		select {
		case semaphore <- struct{}{}:
			go func(match string) {
				defer wg.Done()
				defer func() { <-semaphore }()

				if ctx.Err() != nil {
					return
				}

				info, err := os.Stat(match)
				if err != nil || info.IsDir() {
					return
				}

				mutex.Lock()
				size += uint64(info.Size())
				fileCount++
				if len(paths) < maxPathsToCollect {
					paths = append(paths, match)
				}
				mutex.Unlock()

			}(match)
		case <-ctx.Done():
			wg.Done()
			break
		}
	}

	wg.Wait()
	close(semaphore)

	return size, fileCount, paths
}

// ProcessWalkAction handles recursive directory traversal.
//
// It employs a producer-consumer pattern:
// - CollectFilePaths (Producer): Walks the dir and pushes paths to a channel.
// - ProcessFileWorker (Consumers): 'workers' amount of goroutines read from the channel and stat files.
func ProcessWalkAction(ctx context.Context, searchPath string) (uint64, uint64, []string) {
	var size, fileCount uint64
	var paths []string

	var wg sync.WaitGroup
	var mutex sync.Mutex
	fileChan := make(chan string, 100)

	// start worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go ProcessFileWorker(ctx, fileChan, &size, &fileCount, &paths, &mutex, &wg)
	}

	CollectFilePaths(ctx, searchPath, fileChan)

	close(fileChan)
	wg.Wait()

	return size, fileCount, paths
}

// ProcessFileWorker is the consumer for ProcessWalkAction.
// It reads file paths from a channel, calculates their size, and aggregates results.
func ProcessFileWorker(ctx context.Context, fileChan chan string, size *uint64, fileCount *uint64,
	paths *[]string, mutex *sync.Mutex, wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-fileChan:
			if !ok {
				return
			}

			info, err := os.Stat(path)
			if err != nil || info.IsDir() {
				continue
			}

			mutex.Lock()
			*size += uint64(info.Size())
			*fileCount++
			if len(*paths) < maxPathsToCollect {
				*paths = append(*paths, path)
			}
			mutex.Unlock()
		}
	}
}

// CollectFilePaths is the producer for ProcessWalkAction.
// It walks the directory tree and sends valid file paths to the fileChan.
func CollectFilePaths(ctx context.Context, searchPath string, fileChan chan string) {
	err := filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		select {
		case fileChan <- path:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	if err != nil && !errors.Is(ctx.Err(), err) {
		slog.Error("Error walking directory %s: %v\n", searchPath, err)
	}
}

// ProcessFileAction handles the simplest case: verifying a single specific file path.
func ProcessFileAction(searchPath string) (uint64, uint64, []string) {
	info, err := os.Stat(searchPath)
	if err != nil || info.IsDir() {
		return 0, 0, nil
	}

	return uint64(info.Size()), 1, []string{searchPath}
}

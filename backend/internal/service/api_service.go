package service

import (
	"backend/src/cleaners_util"
	"backend/src/detector"
	"backend/src/structures"
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
var workers = runtime.NumCPU() * 2


// LoadCleanerMap transforms the flat list of cleaners into a nested map structure.
//
// It loads all definitions via cleaners_util and organizes them for O(1) lookup
// during the analysis phase.
// Returns a map keyed by [CleanerID][OptionID] containing the list of Actions.
func LoadCleanerMap() (map[string]map[string][]structures.Action, error) {
	allCleaners, err := cleaners_util.LoadAllCleaners()
	if err != nil {
		slog.Error("Error loading all cleaners: %v\n", err)
		return nil, err
	}

	cleanerMap := make(map[string]map[string][]structures.Action)
	for _, cleaner := range allCleaners {
		cleanerMap[cleaner.ID] = make(map[string][]structures.Action)
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
func AnalyzeRequests(requests []structures.CleanRequest,
	cleanerMap map[string]map[string][]structures.Action) *structures.AnalyzeResponse {
	response := &structures.AnalyzeResponse{
		Items: make([]structures.AnalyzeItem, 0),
	}

	semaphore := make(chan struct{}, workers)
	var wg sync.WaitGroup
	resultsChan := make(chan structures.AnalyzeItem, len(requests))

	for _, request := range requests {
		actions, ok := cleanerMap[request.CleanerID][request.OptionID]
		if !ok {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{} // block if max workers reached

		go func(request structures.CleanRequest, actions []structures.Action) {
			defer wg.Done() // decrease the counter when the goroutine completes
			defer func() { <-semaphore }() // clear the semaphore slot when done

			item := AnalyzeActions(request, actions)
			resultsChan	<- item
		} (request, actions)
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

	return response
}

// AnalyzeActions processes the specific actions (paths/globs) associated with a single cleaner option.
//
// It checks OS compatibility for each action and executes them concurrently using
// the same worker-pool pattern as AnalyzeRequests.
func AnalyzeActions(request structures.CleanRequest, actions []structures.Action) structures.AnalyzeItem {
	var size uint64 = 0
	var fileCount uint64 = 0
	var foundPaths []string

	semaphore := make(chan struct{}, workers)
	var wg sync.WaitGroup
	resultChan := make(chan structures.ActionResult, len(actions))

	for _, action := range actions {
		if !detector.IsOSSupported(action.OS) {
			continue
		}

		wg.Add(1)
		semaphore <- struct{}{}

		go func(action structures.Action) {
			defer wg.Done()
			defer func() { <-semaphore }()

			actionSize, actionCount, actionPaths := ProcessAction(action, 0)

			resultChan <- structures.ActionResult{
				Size:      actionSize,
				FileCount: actionCount,
				Paths:     actionPaths,
			}
		}(action)
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

	return structures.AnalyzeItem{
		CleanerID: request.CleanerID,
		OptionID:  request.OptionID,
		Size:      size,
		FileCount: fileCount,
		Paths:     foundPaths,
	}
}

// ProcessAction acts as a router to determine the correct file discovery strategy.
//
// It expands environment variables in paths (e.g., %APPDATA%) and selects between:
// - Globbing (if "*" is present or explicitly set)
// - Recursive walking ("walk.files")
// - Single file verification
func ProcessAction(action structures.Action, currentPathCount int) (uint64, uint64, []string) {
	searchPath := detector.ExpandPath(action.Path)

	if action.Search == "glob" || strings.Contains(searchPath, "*") {
		return ProcessGlobAction(searchPath, currentPathCount)
	} else if action.Search == "walk.files" {
		return ProcessWalkAction(searchPath)
	}

	return ProcessFileAction(searchPath)
}

// ProcessGlobAction handles file discovery using standard filesystem glob patterns.
//
// It performs a concurrent stat() on all matches found by filepath.Glob.
// Thread-safe: Uses a mutex to safely update the size, count, and path slice.
func ProcessGlobAction(searchPath string, currentPathCount int) (uint64, uint64, []string) {
	var size uint64 = 0
	var fileCount uint64 = 0
	var paths []string

	matches, err := filepath.Glob(searchPath)
	if err != nil {
		log.Printf("Error in glob %s: %v\n", searchPath, err)
		return 0, 0, nil
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	semaphore := make(chan struct{}, workers)

	for _, match := range matches {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(match string) {
			defer wg.Done()
			defer func() { <-semaphore }()

			info, err := os.Stat(match)
			if err != nil || info.IsDir() {
				return
			}

			mutex.Lock()
			size += uint64(info.Size())
			fileCount++
			if currentPathCount+len(paths) < maxPathsToCollect {
				paths = append(paths, match)
			}
			mutex.Unlock()

		}(match)
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
func ProcessWalkAction(searchPath string) (uint64, uint64, []string) {
	var size, fileCount uint64
	var paths []string

	var wg sync.WaitGroup
	var mutex sync.Mutex
	fileChan := make(chan string, 100)


	// start worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go ProcessFileWorker(fileChan, &size, &fileCount, &paths, &mutex, &wg)
	}

	CollectFilePaths(searchPath, fileChan)

	close(fileChan)
	wg.Wait()

	return size, fileCount, paths
}

// ProcessFileWorker is the consumer for ProcessWalkAction.
// It reads file paths from a channel, calculates their size, and aggregates results.
func ProcessFileWorker(fileChan chan string, size *uint64, fileCount *uint64,
	paths *[]string, mutex *sync.Mutex, wg *sync.WaitGroup) {

	defer wg.Done()

	for path := range fileChan {
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

// CollectFilePaths is the producer for ProcessWalkAction.
// It walks the directory tree and sends valid file paths to the fileChan.
func CollectFilePaths(searchPath string, fileChan chan string) {
	err := filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		fileChan <- path
		return nil
	})

	if err != nil {
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

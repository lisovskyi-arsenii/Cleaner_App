package handlers

import (
	"backend/src/cleaners_util"
	"backend/src/detector"
	"backend/src/structures"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const maxPathsToCollect = 10
var workers = runtime.NumCPU() * 2

func GetCleaners(w http.ResponseWriter, _ *http.Request) {
	allCleaners, err := cleaners_util.LoadAllCleaners()
	if err != nil {
		fmt.Printf("Error loading all cleaners: %v\n", err)
		http.Error(w, fmt.Sprintf("Error loading cleaners: %v", err),
			http.StatusInternalServerError)
		return
	}

	installedCleaners, err := cleaners_util.FilterOnlyInstalledCleaners(allCleaners)
	if err != nil {
		fmt.Printf("Error filtering installed cleaners: %v\n", err)
		http.Error(w, fmt.Sprintf("Error filtering installed cleaners: %v", err),
			http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&installedCleaners); err != nil {
		fmt.Printf("Error encoding cleaners: %v\n", err)
		http.Error(w, fmt.Sprintf("Error encoding cleaners: %v", err),
			http.StatusInternalServerError)
		return
	}

	for _, cleaner := range installedCleaners {
		fmt.Printf("Found cleaner %s\n", cleaner.Name)
		fmt.Printf("Option ID: %s\n", cleaner.ID)
		fmt.Printf("Option Name: %s\n", cleaner.Name)
		fmt.Printf("Option Description: %s\n", cleaner.Description)
	}
}

func HandleAnalyze(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	var requests []structures.CleanRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	fmt.Println("DEBUG: Cleaners - ", requests)

	cleanerMap, err := loadCleanerMap()
	if err != nil {
		http.Error(w, "Error loading cleaners", http.StatusInternalServerError)
		return
	}

	response := analyzeRequests(requests, cleanerMap)

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		fmt.Printf("Error encoding cleaners: %v\n", err)
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("DEBUG: AnalyzeResponse - ", *response)
}


func loadCleanerMap() (map[string]map[string][]structures.Action, error) {
	allCleaners, err := cleaners_util.LoadAllCleaners()
	if err != nil {
		fmt.Printf("Error loading all cleaners: %v\n", err)
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

func analyzeRequests(requests []structures.CleanRequest,
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

			item := analyzeActions(request, actions)
			resultsChan	<- item
		} (request, actions)
	}

	wg.Wait()
	close(resultsChan)
	close(semaphore)

	for item := range resultsChan {
		response.Items = append(response.Items, item)
		response.TotalSize += item.Size
		response.TotalFiles += item.FileCount
	}

	return response
}

func analyzeActions(request structures.CleanRequest, actions []structures.Action) structures.AnalyzeItem {
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

		go func(action structures.Action, currentPathCount int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			actionSize, actionCount, actionPaths := processAction(action, currentPathCount)

			resultChan <- structures.ActionResult{
				Size:      actionSize,
				FileCount: actionCount,
				Paths:     actionPaths,
			}
		}(action, 0)
	}

	wg.Wait()
	close(resultChan)
	close(semaphore)

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

func processAction(action structures.Action, currentPathCount int) (uint64, uint64, []string) {
	searchPath := detector.ExpandPath(action.Path)

	if action.Search == "glob" || strings.Contains(searchPath, "*") {
		return processGlobAction(searchPath, currentPathCount)
	} else if action.Search == "walk.files" {
		return processWalkAction(searchPath, currentPathCount)
	}

	return processFileAction(searchPath)
}

func processGlobAction(searchPath string, currentPathCount int) (uint64, uint64, []string) {
	var size uint64 = 0
	var fileCount uint64 = 0
	var paths []string

	matches, err := filepath.Glob(searchPath)
	if err != nil {
		fmt.Printf("Error in glob %s: %v\n", searchPath, err)
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

func processWalkAction(searchPath string, currentPathCount int) (uint64, uint64, []string) {
	var size, fileCount uint64
	var paths []string

	var wg sync.WaitGroup
	var mutex sync.Mutex
	fileChan := make(chan string, 100)


	// start worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go processFileWorker(fileChan, &size, &fileCount, &paths, currentPathCount, &mutex, &wg)
	}

	collectFilePaths(searchPath, fileChan)

	wg.Wait()
	close(fileChan)

	return size, fileCount, paths
}

func processFileWorker(fileChan chan string, size *uint64, fileCount *uint64,
	paths *[]string, currentPathCount int, mutex *sync.Mutex, wg *sync.WaitGroup) {

	defer wg.Done()

	for path := range fileChan {
		info, err := os.Stat(path)
		if err != nil || info.IsDir() {
			continue
		}

		mutex.Lock()
		*size += uint64(info.Size())
		*fileCount++
		if currentPathCount+len(*paths) < maxPathsToCollect {
			*paths = append(*paths, path)
		}
		mutex.Unlock()
	}
}

func collectFilePaths(searchPath string, fileChan chan string) {
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
		fmt.Printf("Error walking directory %s: %v\n", searchPath, err)
	}
}

func processFileAction(searchPath string) (uint64, uint64, []string) {
	info, err := os.Stat(searchPath)
	if err != nil || info.IsDir() {
		return 0, 0, nil
	}

	return uint64(info.Size()), 1, []string{searchPath}
}

package handlers

import (
	"backend/src/cleaners_util"
	"backend/src/detector"
	"backend/src/structures"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
	err = json.NewEncoder(w).Encode(&installedCleaners)
	if err != nil {
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
		http.Error(w, "Error encoding response: %v", http.StatusInternalServerError)
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

	for _, request := range requests {
		actions, ok := cleanerMap[request.CleanerID][request.OptionID]
		if !ok {
			continue
		}

		item := analyzeActions(request, actions)
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

	for _, action := range actions {
		if !detector.IsOSSupported(action.OS) {
			continue
		}

		actionSize, actionCount, actionPaths := processAction(action, len(foundPaths))
		size += actionSize
		fileCount += actionCount
		foundPaths = append(foundPaths, actionPaths...)
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
	} else {
		return processFileAction(searchPath)
	}
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

	for _, match := range matches {
		if info, err := os.Stat(match); err == nil && !info.IsDir() {
			size += uint64(info.Size())
			fileCount++
			if currentPathCount+len(paths) < 10 {
				paths = append(paths, match)
			}
		}
	}

	return size, fileCount, paths
}

func processWalkAction(searchPath string, currentPathCount int) (uint64, uint64, []string) {
	var size uint64 = 0
	var fileCount uint64 = 0
	var paths []string

	err := filepath.WalkDir(searchPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		size += uint64(info.Size())
		fileCount++
		if currentPathCount+len(paths) < 10 {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking directory %s: %v\n", searchPath, err)
	}

	return size, fileCount, paths
}

func processFileAction(searchPath string) (uint64, uint64, []string) {
	info, err := os.Stat(searchPath)
	if err != nil || info.IsDir() {
		return 0, 0, nil
	}

	return uint64(info.Size()), 1, []string{searchPath}
}

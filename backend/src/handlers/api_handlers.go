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
		return
	}

	installedCleaners, err := cleaners_util.FilterOnlyInstalledCleaners(allCleaners)
	if err != nil {
		fmt.Printf("Error filtering installed cleaners: %v\n", err)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&installedCleaners)
	if err != nil {
		fmt.Printf("Error encoding cleaners: %v\n", err)
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


	allCleaners, err := cleaners_util.LoadAllCleaners()
	if err != nil {
		fmt.Printf("Error loading all cleaners: %v\n", err)
		return
	}

	cleanerMap := make(map[string]map[string][]structures.Action)

	for _, cleaner := range allCleaners {
		cleanerMap[cleaner.ID] = make(map[string][]structures.Action)
		for _, option := range cleaner.Options {
			cleanerMap[cleaner.ID][option.ID] = option.Actions
		}
	}


	response := &structures.AnalyzeResponse{
		Items: make([]structures.AnalyzeItem, 0),
	}

	for _, request := range requests {
		actions, ok := cleanerMap[request.CleanerID][request.OptionID]
		if !ok {
			continue
		}

		var size uint64 = 0
		var fileCount uint64 = 0
		var foundPaths []string

		for _, action := range actions {
			if !detector.IsOSSupported(action.OS) {
				continue
			}

			searchPath := detector.ExpandPath(action.Path)


			if action.Search == "glob" || strings.Contains(searchPath, "*") {
				matches, _ := filepath.Glob(searchPath)
				for _, match := range matches {
					if info, err := os.Stat(match); err == nil && !info.IsDir() {
						size += uint64(info.Size())
						fileCount++
						if len(foundPaths) < 10 {
							foundPaths = append(foundPaths, match)
						}
					}
				}
			} else if action.Search == "walk.files" {
				err := filepath.WalkDir(searchPath, func(path string, d os.DirEntry, err error) error {
					if err == nil && !d.IsDir() {
						info, _ := d.Info()
						size += uint64(info.Size())
						fileCount++
						if len(foundPaths) < 10 {
							foundPaths = append(foundPaths, path)
						}
					}
					return nil
				})
				if err != nil {
					return
				}
			} else {
				if info, err := os.Stat(searchPath); err == nil && !info.IsDir() {
					size += uint64(info.Size())
					fileCount++
					foundPaths = append(foundPaths, searchPath)
				}
			}
		}

		item := &structures.AnalyzeItem{
			CleanerID: request.CleanerID,
			OptionID:  request.OptionID,
			Size:      size,
			FileCount: fileCount,
			Paths:     foundPaths,
		}

		response.Items = append(response.Items, *item)
		response.TotalSize += size
		response.TotalFiles += fileCount
	}


	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&response)
	if err != nil {
		fmt.Printf("Error encoding cleaners: %v\n", err)
		return
	}

	fmt.Println("DEBUG: AnalyzeResponse - ", *response)
}

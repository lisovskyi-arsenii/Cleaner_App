package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)


func loadAllCleaners() ([]Cleaner, error) {
	cleanersDir := "./resources"

	files, err := os.ReadDir(cleanersDir)
	if err != nil {
		fmt.Printf("Error reading dir %s: %v\n", cleanersDir, err)
		return nil, err
	}

	var cleaners []Cleaner

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(cleanersDir, file.Name())

		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			continue
		}

		var cleaner Cleaner
		if err := json.Unmarshal(data, &cleaner); err != nil {
			fmt.Printf("Error parsing file %s: %v\n", filePath, err)
			continue
		}

		cleaners = append(cleaners, cleaner)
		fmt.Printf("Loaded cleaner %s from %s\n", cleaner.Name, file.Name())
	}

	fmt.Printf("Total cleaners: %d\n", len(cleaners))
	return cleaners, nil
}


func filterOnlyInstalledCleaners(cleaners []Cleaner) ([]Cleaner, error) {
	var installedCleaners []Cleaner

	for _, cleaner := range cleaners {
		if DetectInstalled(cleaner.Detect) {
			installedCleaners = append(installedCleaners, cleaner)
		}
	}

	return installedCleaners, nil
}


func getCleaners(w http.ResponseWriter, r *http.Request) {
	allCleaners, err := loadAllCleaners()
	if err != nil {
		fmt.Printf("Error loading all cleaners: %v\n", err)
		return
	}

	installedCleaners, err := filterOnlyInstalledCleaners(allCleaners)
	if err != nil {
		fmt.Printf("Error filtering installed cleaners: %v\n", err)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&installedCleaners)
	for _, cleaner := range installedCleaners {
		fmt.Printf("Found cleaner %s\n", cleaner.Name)
		fmt.Printf("Option ID: %s\n", cleaner.ID)
		fmt.Printf("Option Name: %s\n", cleaner.Name)
		fmt.Printf("Option Description: %s\n", cleaner.Description)
	}
}

func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	var requests []CleanRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	fmt.Println("DEBUG: Cleaners - ", requests)


	allCleaners, err := loadAllCleaners()
	if err != nil {
		fmt.Printf("Error loading all cleaners: %v\n", err)
		return
	}

	cleanerMap := make(map[string]map[string][]Action)

	for _, cleaner := range allCleaners {
		cleanerMap[cleaner.ID] = make(map[string][]Action)
		for _, option := range cleaner.Options {
			cleanerMap[cleaner.ID][option.ID] = option.Actions
		}
	}



	response := &AnalyzeResponse{
		Items: make([]AnalyzeItem, 0),
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
			if !isOSSupported(action.OS) {
				return
			}

			searchPath := expandPath(action.Path)


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
				filepath.WalkDir(searchPath, func(path string, d os.DirEntry, err error) error {
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
			} else {
				if info, err := os.Stat(searchPath); err == nil && !info.IsDir() {
					size += uint64(info.Size())
					fileCount++
					foundPaths = append(foundPaths, searchPath)
				}
			}
		}

		item := &AnalyzeItem{
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
	json.NewEncoder(w).Encode(&response)

	fmt.Println("DEBUG: AnalyzeResponse - ", *response)
}



func main() {
	http.Handle("/api/cleaners", http.HandlerFunc(getCleaners))
	http.Handle("/api/analyze", http.HandlerFunc(handleAnalyze))
	log.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

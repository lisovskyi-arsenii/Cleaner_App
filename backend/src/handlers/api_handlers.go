// Package handlers provides the HTTP request handlers for the application's backend API.
// It includes functionality for discovering installed system cleaners, reviewing
// cleanup targets, and executing cleaning operations
package handlers

import (
	"backend/src/cleaners_util"
	"backend/src/service"
	"backend/src/structures"
	"encoding/json"
	"fmt"
	"net/http"
)


// GetCleaners handles the discovery of system cleaners.
//
// It loads all available cleaner definitions, checks which ones are actually
// installed on the host system, and returns the filtered list as a JSON response.
//
// GET /api/cleaners
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

// HandleAnalyze processes requests to analyze specific cleanup targets.
//
// It expects a JSON body containing a list of structures.CleanRequest.
// It loads the cleaner configuration map, performs the analysis to determine
// space to be freed or files to be removed, and returns a detailed JSON response.
//
// POST /api/analyze
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

	cleanerMap, err := service.LoadCleanerMap()
	if err != nil {
		http.Error(w, "Error loading cleaners", http.StatusInternalServerError)
		return
	}

	response := service.AnalyzeRequests(requests, cleanerMap)

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		fmt.Printf("Error encoding cleaners: %v\n", err)
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("DEBUG: AnalyzeResponse - ", *response)
}

// HandleClean executes the cleanup process.
//
// It is intended to receive confirmation of items to delete and perform
// the actual file removal operations.
//
// POST /api/clean
func HandleClean(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the cleaning logic
}



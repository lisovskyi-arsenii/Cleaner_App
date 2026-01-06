package handlers

import (
	"backend/src/cleaners_util"
	"backend/src/structures"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

const maxPathsToCollect = 500
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

	cleanerMap, err := LoadCleanerMap()
	if err != nil {
		http.Error(w, "Error loading cleaners", http.StatusInternalServerError)
		return
	}

	response := AnalyzeRequests(requests, cleanerMap)

	if err := json.NewEncoder(w).Encode(&response); err != nil {
		fmt.Printf("Error encoding cleaners: %v\n", err)
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("DEBUG: AnalyzeResponse - ", *response)
}

func HandleClean(w http.ResponseWriter, r *http.Request) {
	
}



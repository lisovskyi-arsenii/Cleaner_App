package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	//installed := make([]Cleaner, 0)
	//
	// for _, cleaner := range cleaners {
    //    // System завжди показувати (немає detection)
    //    if len(cleaner.Detection.Paths) == 0 &&
    //       len(cleaner.Detection.Registry) == 0 {
    //        installed = append(installed, cleaner)
    //        log.Printf("ℹ️  Always show: %s", cleaner.Name)
    //        continue
    //    }
	//
    //    // Перевірити чи встановлено
    //    if detector.DetectInstalled(cleaner.Detection) {
    //        installed = append(installed, cleaner)
    //        log.Printf("✅ Detected: %s", cleaner.Name)
    //    } else {
    //        log.Printf("❌ Not found: %s", cleaner.Name)
    //    }
    //}
	return cleaners, nil
}


func getCleaners(w http.ResponseWriter, r *http.Request) {
	//// test
	//cleaners := &[]Cleaner{
	//	{
	//		ID:          "firefox",
	//		Name:        "Mozilla Firefox",
	//		Description: "Web browser",
	//		Running:     false,
	//		Options: []Option{
	//			{
	//				ID:          "cache",
	//				Label:       "Cache",
	//				Description: "Delete the web cache",
	//				Warning:     "This may slow down first browser start",
	//				Actions: []Action{
	//					{
	//						Command: "delete",
	//						Search:  "walk.files",
	//						Path:    "%AppData%/Mozilla/Firefox/Profiles/*/cache2",
	//						OS:      []string{"windows"},
	//						Type:    "f",
	//					},
	//					{
	//						Command: "delete",
	//						Search:  "walk.files",
	//						Path:    "~/.mozilla/firefox/*/cache2",
	//						OS:      []string{"linux"},
	//						Type:    "f",
	//					},
	//				},
	//			},
	//			{
	//				ID:          "cookies",
	//				Label:       "Cookies",
	//				Description: "Delete cookies which track you",
	//				Warning:     "You will be logged out of websites",
	//				Actions: []Action{
	//					{
	//						Command: "delete",
	//						Search:  "glob",
	//						Path:    "%AppData%/Mozilla/Firefox/Profiles/*/cookies.sqlite*",
	//						OS:      []string{"windows"},
	//					},
	//					{
	//						Command: "delete",
	//						Search:  "glob",
	//						Path:    "~/.mozilla/firefox/*/cookies.sqlite*",
	//						OS:      []string{"linux"},
	//					},
	//				},
	//			},
	//			{
	//				ID:          "history",
	//				Label:       "History",
	//				Description: "Delete browsing history",
	//				Actions: []Action{
	//					{
	//						Command: "vacuum",
	//						Search:  "file",
	//						Path:    "%AppData%/Mozilla/Firefox/Profiles/*/places.sqlite",
	//						OS:      []string{"windows"},
	//					},
	//				},
	//			},
	//		},
	//	},
	//	{
	//		ID:          "system",
	//		Name:        "System",
	//		Description: "Windows system files",
	//		Running:     false,
	//		Options: []Option{
	//			{
	//				ID:          "temp",
	//				Label:       "Temporary files",
	//				Description: "Delete temporary files",
	//				Actions: []Action{
	//					{
	//						Command: "delete",
	//						Search:  "walk.all",
	//						Path:    "%TEMP%",
	//						OS:      []string{"windows"},
	//					},
	//				},
	//			},
	//			{
	//				ID:          "recycle",
	//				Label:       "Recycle Bin",
	//				Description: "Empty recycle bin",
	//				Warning:     "Cannot be undone!",
	//				Actions: []Action{
	//					{
	//						Command: "delete",
	//						Search:  "walk.all",
	//						Path:    "C:\\$Recycle.Bin",
	//						OS:      []string{"windows"},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}

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
	var requests []CleanRequest
	json.NewDecoder(r.Body).Decode(&requests)

	fmt.Println("DEBUG: Cleaners - ", requests)

	response := &AnalyzeResponse{
		TotalSize: uint64(len(requests)),
		TotalFiles: uint64(len(requests)),
		Items:      make([]AnalyzeItem, len(requests)),
	}

	for _, request := range requests {
		item := &AnalyzeItem{
			CleanerID: request.CleanerID,
			OptionID:  request.OptionID,
			Size: 50000000,
			FileCount: uint64(len(requests)),
			Paths: []string{"C:/test/file1.txt", "C:/test/file2.txt"},
		}

		response.Items = append(response.Items, *item)
		response.TotalSize += item.Size
		response.TotalFiles += item.FileCount
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

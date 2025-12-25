package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func getCleaners(w http.ResponseWriter, r *http.Request) {
	// test
	cleaners := []Cleaner{
		{
			ID:          "firefox",
			Name:        "Mozilla Firefox",
			Description: "Web browser",
			Running:     false,
			Options: []Option{
				{
					ID:          "cache",
					Label:       "Cache",
					Description: "Delete the web cache",
					Warning:     "This may slow down first browser start",
					Actions: []Action{
						{
							Command: "delete",
							Search:  "walk.files",
							Path:    "%AppData%/Mozilla/Firefox/Profiles/*/cache2",
							OS:      []string{"windows"},
							Type:    "f",
						},
						{
							Command: "delete",
							Search:  "walk.files",
							Path:    "~/.mozilla/firefox/*/cache2",
							OS:      []string{"linux"},
							Type:    "f",
						},
					},
				},
				{
					ID:          "cookies",
					Label:       "Cookies",
					Description: "Delete cookies which track you",
					Warning:     "You will be logged out of websites",
					Actions: []Action{
						{
							Command: "delete",
							Search:  "glob",
							Path:    "%AppData%/Mozilla/Firefox/Profiles/*/cookies.sqlite*",
							OS:      []string{"windows"},
						},
						{
							Command: "delete",
							Search:  "glob",
							Path:    "~/.mozilla/firefox/*/cookies.sqlite*",
							OS:      []string{"linux"},
						},
					},
				},
				{
					ID:          "history",
					Label:       "History",
					Description: "Delete browsing history",
					Actions: []Action{
						{
							Command: "vacuum",
							Search:  "file",
							Path:    "%AppData%/Mozilla/Firefox/Profiles/*/places.sqlite",
							OS:      []string{"windows"},
						},
					},
				},
			},
		},
		{
			ID:          "system",
			Name:        "System",
			Description: "Windows system files",
			Running:     false,
			Options: []Option{
				{
					ID:          "temp",
					Label:       "Temporary files",
					Description: "Delete temporary files",
					Actions: []Action{
						{
							Command: "delete",
							Search:  "walk.all",
							Path:    "%TEMP%",
							OS:      []string{"windows"},
						},
					},
				},
				{
					ID:          "recycle",
					Label:       "Recycle Bin",
					Description: "Empty recycle bin",
					Warning:     "Cannot be undone!",
					Actions: []Action{
						{
							Command: "delete",
							Search:  "walk.all",
							Path:    "C:\\$Recycle.Bin",
							OS:      []string{"windows"},
						},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cleaners)
}

func main() {
	http.Handle("/api/cleaners", http.HandlerFunc(getCleaners))
	log.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

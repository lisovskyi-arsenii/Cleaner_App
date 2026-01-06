package main

import (
	"backend/src/handlers"
	"log"
	"net/http"
)

func main() {
	http.Handle("/api/cleaners", http.HandlerFunc(handlers.GetCleaners))
	http.Handle("/api/analyze", http.HandlerFunc(handlers.HandleAnalyze))
	http.Handle("/api/clean", http.HandlerFunc(handlers.HandleClean))
	log.Println("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

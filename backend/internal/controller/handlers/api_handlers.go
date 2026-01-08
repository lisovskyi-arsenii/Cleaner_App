// Package handlers provides the HTTP request handlers for the application's backend API.
// It includes functionality for discovering installed system cleaners, reviewing
// cleanup targets, and executing cleaning operations
package handlers

import (
	"backend/src/cleaners_util"
	"backend/src/service"
	"backend/src/structures"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCleaners handles the discovery of system cleaners.
//
// It loads all available cleaner definitions, checks which ones are actually
// installed on the host system, and returns the filtered list as a JSON response.
//
// GET /api/cleaners
func GetCleaners(c *gin.Context) {
	allCleaners, err := cleaners_util.LoadAllCleaners()
	if err != nil {
		slog.Error("Error loading all cleaners: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error loading cleaners: %v", err)})
		return
	}

	installedCleaners, err := cleaners_util.FilterOnlyInstalledCleaners(allCleaners)
	if err != nil {
		slog.Error("Error filtering installed cleaners: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error filtering installed cleaners: %v", err)})
		return
	}

	c.JSON(http.StatusOK, &installedCleaners)

	for _, cleaner := range installedCleaners {
		slog.Info("Found cleaner %s\n", cleaner.Name)
		slog.Info("Option ID: %s\n", cleaner.ID)
		slog.Info("Option Name: %s\n", cleaner.Name)
		slog.Info("Option Description: %s\n", cleaner.Description)
	}
}

// HandlePreview processes requests to analyze specific cleanup targets.
//
// It expects a JSON body containing a list of structures.CleanRequest.
// It loads the cleaner configuration map, performs the analysis to determine
// space to be freed or files to be removed, and returns a detailed JSON response.
//
// POST /api/preview
func HandlePreview(c *gin.Context) {
	var requests []structures.CleanRequest
	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	log.Println("DEBUG: Cleaners - ", requests)

	cleanerMap, err := service.LoadCleanerMap()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error loading cleaners: %v", err)})
		return
	}

	response := service.AnalyzeRequests(requests, cleanerMap)
	slog.Debug("DEBUG: AnalyzeResponse - ", *response)

	c.JSON(http.StatusOK, &response)
}

// HandleClean executes the cleanup process.
//
// It is intended to receive confirmation of items to delete and perform
// the actual file removal operations.
//
// POST /api/clean
func HandleClean(c *gin.Context) {
	// TODO: Implement the cleaning logic
}

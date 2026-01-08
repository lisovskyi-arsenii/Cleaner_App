package cleaners

import (
	"backend/internal/detector"
	"backend/internal/models"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
)

func LoadAllCleaners() ([]structures.Cleaner, error) {
	cleanersDir := "./resources"

	files, err := os.ReadDir(cleanersDir)
	if err != nil {
		slog.Error("Error reading dir %s: %v\n", cleanersDir, err)
		return nil, err
	}

	var cleaners []structures.Cleaner

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(cleanersDir, file.Name())

		data, err := os.ReadFile(filePath)
		if err != nil {
			slog.Error("Error reading file %s: %v\n", filePath, err)
			continue
		}

		var cleaner structures.Cleaner
		if err := json.Unmarshal(data, &cleaner); err != nil {
			slog.Error("Error parsing file %s: %v\n", filePath, err)
			continue
		}

		cleaners = append(cleaners, cleaner)
		slog.Info("Loaded cleaner %s from %s\n", cleaner.Name, file.Name())
	}

	slog.Debug("Total cleaners: %d\n", len(cleaners))
	return cleaners, nil
}


func FilterOnlyInstalledCleaners(cleaners []structures.Cleaner) ([]structures.Cleaner, error) {
	var installedCleaners []structures.Cleaner

	for _, cleaner := range cleaners {
		if detector.DetectInstalled(cleaner.Detect) {
			installedCleaners = append(installedCleaners, cleaner)
		}
	}

	return installedCleaners, nil
}

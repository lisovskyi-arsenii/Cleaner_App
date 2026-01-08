package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// SetupLogger initializes the global logger with a specific level.
// It configures file rotation for all logs and optionally writes to stdout for debug levels.
func SetupLogger(logDir string, level slog.Level) {
	var writer io.Writer
	var handler slog.Handler

	// Ensure log directory exists
	// We ignore the error here because if we can't create the dir,
	// the lumberjack logger will likely fail or just log to stderr anyway.
	_ = os.MkdirAll(logDir, 0755)

	// Configure Log Rotation
	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "app.log"),
		MaxSize:    10,   // Megabytes before rotating
		MaxBackups: 3,    // Keep 3 old files
		MaxAge:     28,   // Days to keep files
		Compress:   true, // Compress old files (gzip)
	}

	// Logic:
	// - If Level is Debug: Write to BOTH Console (Stdout) and File. Use Text format for readability.
	// - If Level is Info/Warn/Error: Write ONLY to File. Use JSON format for parsing tools.
	if level == slog.LevelDebug {
		writer = io.MultiWriter(os.Stdout, fileWriter)
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		writer = fileWriter
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: level,
		})
	}

	// Create the logger and set it as the global default
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
package main

import (
	"backend/internal/controller/handlers"
	"backend/internal/logger"
	"backend/internal/middleware"
	"backend/internal/routes"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnvFile(filenames... string) {
	if err := godotenv.Load(filenames...); err != nil {
		slog.Warn("No .env file found (using system env vars)", "error", err)
	}
}

func loadPortFromEnvFile() string {
	port := os.Getenv("PORT")
	if port == "" {
		slog.Error("$PORT must be set")
		os.Exit(1)
	}
	return port
}

func getLogLevel() slog.Level {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// main function entry point to the program
func main() {
	loadEnvFile(".env")

	// setup logger for whole project
	logLevel := getLogLevel()
	logger.SetupLogger("./logs", logLevel)

	slog.Info("Environment loaded", "level", logLevel.String())

	// Set Gin to Release mode if we aren't in debug to keep console clean
	if logLevel != slog.LevelDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	// setup router
	router := gin.New()
	router.Use(gin.Recovery())

	// middleware
	router.Use(middleware.SlogLogger())
	router.Use(middleware.CORSMiddleware())

	// create API group '/api/*'
	api := router.Group(routes.APIGroup)
	{
		api.GET(routes.GetCleaners, handlers.GetCleaners)
		api.POST(routes.Preview, handlers.HandlePreview)
		api.POST(routes.Clean, handlers.HandleClean)
	}

	// load port from .env file
	port := loadPortFromEnvFile()
	slog.Info("Server started on port " + port)

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		slog.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}

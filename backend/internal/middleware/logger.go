package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// SlogLogger replaces Gin's default logger to use our custom slog setup
func SlogLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()
		errors := c.Errors.String()

		// Log using our global slog instance
		// We use LogAttrs for high performance, but simple slog.Info works too
		slog.Info("HTTP Request",
			"status", status,
			"method", method,
			"path", path,
			"query", query,
			"ip", clientIP,
			"latency", latency,
			"errors", errors,
		)
	}
}
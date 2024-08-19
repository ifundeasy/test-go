package middleware

import (
	"encoding/json"
	"time"

	"test-go/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
)

// logEntry represents the structure of a log entry
type logEntry struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	URL       string `json:"url"`
	Status    int    `json:"status"`
	Duration  string `json:"duration"`
	Error     string `json:"error,omitempty"`
}

// LoggingMiddleware logs the details of each request and response in a single line JSON format
func LoggingMiddleware(logger *logging.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Stop timer
		duration := time.Since(start)

		// Create log entry
		entry := logEntry{
			Level:     "info",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Method:    c.Method(),
			URL:       c.OriginalURL(),
			Status:    c.Response().StatusCode(),
			Duration:  duration.String(),
		}

		// Add error to log entry if exists
		if err != nil {
			entry.Error = err.Error()
		}

		// Serialize log entry to JSON
		logLine, _ := json.Marshal(entry)

		// Log the request details
		logger.Info(string(logLine))

		// Return the error if there was one
		return err
	}
}

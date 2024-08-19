package logging

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// Logger is a simple wrapper around the standard log package
type Logger struct {
	logger *log.Logger
}

// logEntry represents the structure of a log entry
type logEntry struct {
	Level     string      `json:"level"`
	Timestamp string      `json:"timestamp"`
	Message   interface{} `json:"message"`
}

// NewLogger creates a new instance of Logger
func NewLogger(prefix string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, prefix, 0), // Disable standard flags since we will use JSON format
	}
}

// logInternal handles the internal logic for logging a message in JSON format
func (l *Logger) logInternal(level string, msg interface{}) {
	entry := logEntry{
		Level:     level,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Message:   msg,
	}
	logLine, _ := json.Marshal(entry)
	l.logger.Println(string(logLine))
}

// Info logs informational messages
func (l *Logger) Info(msg string) {
	l.logInternal("info", msg)
}

// InfoJSON logs informational messages in JSON format
func (l *Logger) InfoJSON(msg interface{}) {
	l.logInternal("info", msg)
}

// Warn logs warning messages
func (l *Logger) Warn(msg string) {
	l.logInternal("warn", msg)
}

// Error logs error messages
func (l *Logger) Error(msg string) {
	l.logInternal("error", msg)
}

// Fatal logs a fatal message and exits the application
func (l *Logger) Fatal(msg string) {
	l.logInternal("fatal", msg)
	os.Exit(1)
}

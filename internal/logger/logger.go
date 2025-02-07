package logger

import (
	"log/slog"
	"os"
)

// Initializes the logger to log in JSON format
func Init() *slog.Logger {
	// Create a JSON encoder for the logs
	encoder := slog.NewJSONHandler(os.Stdout, nil)

	// Create a logger with the JSON handler
	logger := slog.New(encoder)

	// Optional: Set the log level (you can adjust as needed)
	logger.Info("Logger initialized")

	return logger
}

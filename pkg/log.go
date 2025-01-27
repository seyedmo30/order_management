package pkg

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

// GetLogger returns an instance of the configured logger.
func GetLogger() *slog.Logger {
	if logger == nil {
		logger = initLogger()
	}
	return logger
}

// initLogger initializes a new JSON logger that writes to stdout.
func initLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

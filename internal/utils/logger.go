package utils

import (
	"os"

	"github.com/charmbracelet/log"
)

func NewLogger(noColor bool) *log.Logger {
	options := log.Options{
		ReportTimestamp: true,
		TimeFormat:      "15:04:05",
		Level:           log.DebugLevel,
	}
	if noColor {
		// Simple approach: use a custom formatter
	}
	logger := log.NewWithOptions(os.Stderr, options)
	return logger
}
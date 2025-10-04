package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// level: "debug", "info", "warn", "error"
// pretty: if true, enables human-readable console output (useful for local dev)
func NewLogger(level string, pretty bool) zerolog.Logger {
	// Set timestamp format
	zerolog.TimeFieldFormat = time.RFC3339

	// Parse log level
	logLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// Ensure logs directory exists
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		panic("Failed to create logs directory: " + err.Error())
	}

	// Create (or append) to log file
	logFilePath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	// Prepare output writers
	var writers []io.Writer

	if pretty {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
		writers = append(writers, consoleWriter)
	} else {
		writers = append(writers, os.Stdout)
	}

	// Always include file output
	writers = append(writers, logFile)

	multi := io.MultiWriter(writers...)

	logger := zerolog.New(multi).With().
		Timestamp().
		Str("service", "app").
		Logger()

	return logger
}

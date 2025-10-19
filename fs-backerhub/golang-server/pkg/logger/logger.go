package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logging struct {
	Level  string `json:"level" mapstructure:"level"`
	Path   string `json:"path" mapstructure:"path"`
	Format string `json:"format" mapstructure:"format"` // "json" | "console"
	Output string `json:"output" mapstructure:"output"` // "stdout" | "stderr" | "file" | "both"
	Caller bool   `json:"caller" mapstructure:"caller"` // true | false
}

func NewZeroLogger(cfg Logging) zerolog.Logger {
	if cfg.Path == "" {
		cfg.Path = "logs/app.log"
	}

	var writers []io.Writer

	fileWriter := mustOpenFile(cfg.Path)

	switch strings.ToLower(cfg.Output) {
	case "stdout":
		writers = append(writers, os.Stdout)
	case "stderr":
		writers = append(writers, os.Stderr)
	case "file":
		writers = append(writers, mustOpenFile(cfg.Path))
	case "both", "all", "":
		if cfg.Format == "console" {
			consoleWriter := zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: "2006-01-02 15:04:05",
			}
			consoleWriter.FormatLevel = func(i interface{}) string {
				return strings.ToUpper(fmt.Sprintf("%-5s", i))
			}
			writers = append(writers, consoleWriter)
		} else {
			writers = append(writers, os.Stdout)
		}
		writers = append(writers, fileWriter)
	}

	writer := io.MultiWriter(writers...)

	level := parseLevel(cfg.Level)

	builder := zerolog.New(writer).
		Level(level).
		With().
		Timestamp()

	if cfg.Caller {
		builder = builder.Caller()
	}

	logger := builder.Logger()
	log.Logger = logger

	return logger
}

func mustOpenFile(path string) *os.File {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatal().Err(err).Msg("failed to create log directory")
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open log file")
	}
	return file
}

func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

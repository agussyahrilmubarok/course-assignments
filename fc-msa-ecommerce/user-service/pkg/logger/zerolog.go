package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"example.com/user-service/pkg/config"
	"github.com/rs/zerolog"
)

func NewZerolog(cfg *config.Config) zerolog.Logger {
	var writers []io.Writer

	if cfg.Log.Output == "stdout" || cfg.Log.Output == "both" {
		if cfg.Log.PrettyConsole {
			consoleWriter := zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}
			writers = append(writers, consoleWriter)
		} else {
			writers = append(writers, os.Stdout)
		}
	}

	if cfg.Log.Output == "file" || cfg.Log.Output == "both" {
		if cfg.Log.FilePath == "" {
			cfg.Log.FilePath = "app.log"
		}
		dir := filepath.Dir(cfg.Log.FilePath)
		if dir != "." { 
			if err := os.MkdirAll(dir, 0755); err != nil {
				panic("cannot create log directory: " + err.Error())
			}
		}
		file, err := os.OpenFile(cfg.Log.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("cannot open log file: " + err.Error())
		}
		writers = append(writers, file)
	}

	multi := io.MultiWriter(writers...)

	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Log.Level))
	if err != nil {
		level = zerolog.InfoLevel
	}

	logger := zerolog.New(multi).
		Level(level).
		With().
		Timestamp().
		Logger()

	return logger
}

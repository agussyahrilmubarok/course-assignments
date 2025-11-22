package config

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

func NewZerolog(cfg *Config) (zerolog.Logger, error) {
	logDir := filepath.Dir(cfg.Logger.Filepath)
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return zerolog.Logger{}, err
	}

	logFile, err := os.OpenFile(cfg.Logger.Filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return zerolog.Logger{}, err
	}

	multi := zerolog.MultiLevelWriter(os.Stdout, logFile)

	level, err := zerolog.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	logger := zerolog.New(multi).
		With().
		Timestamp().
		Logger()

	return logger, nil
}

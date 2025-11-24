package config

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logging struct {
	Logger *logrus.Logger
}

func NewLogging(cfg *Config) (*Logging, error) {
	log := logrus.New()

	level, err := logrus.ParseLevel(strings.ToLower(cfg.Logging.Level))
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	log.SetLevel(level)

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})

	log.SetOutput(os.Stdout)

	if cfg.Logging.Filepath != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.Logging.Filepath,
			MaxSize:    20,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		}

		log.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
	}

	return &Logging{Logger: log}, nil
}

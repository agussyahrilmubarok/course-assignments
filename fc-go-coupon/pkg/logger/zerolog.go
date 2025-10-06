package logger

import (
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	"example.com/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"gopkg.in/natefinch/lumberjack.v2"
)

var syncOnce sync.Once
var log zerolog.Logger

func NewZerolog(cfg *config.Config) zerolog.Logger {
	syncOnce.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		logLevel, err := strconv.Atoi(cfg.Server.LogLevel)
		if err != nil {
			logLevel = int(zerolog.InfoLevel) // default to INFO
		}

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if cfg.Server.Environment != "development" {
			fileLogger := &lumberjack.Logger{
				Filename:   cfg.Server.LogFilepath,
				MaxSize:    5, //
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("service-name", cfg.Server.Name).
			Logger()
	})

	return log
}

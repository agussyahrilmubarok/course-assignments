package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logging struct {
	Logger *zap.Logger
}

func NewLogging(cfg *Config) (*Logging, error) {
	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText([]byte(cfg.Logging.Level)); err != nil {
		logLevel = zapcore.InfoLevel
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Logging.Filepath,
		MaxSize:    100, // MB
		MaxBackups: 7,
		MaxAge:     30,   // days
		Compress:   true, // gzip
	})

	consoleWriter := zapcore.Lock(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), consoleWriter, logLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileWriter, logLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logging{Logger: logger}, nil
}

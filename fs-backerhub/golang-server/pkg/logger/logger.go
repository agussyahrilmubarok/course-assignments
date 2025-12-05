package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func NewLogger(level, filepath, gelfAddr string) error {
	var err error
	once.Do(func() {
		var logLevel zapcore.Level
		if e := logLevel.UnmarshalText([]byte(level)); e != nil {
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
			Filename:   filepath,
			MaxSize:    100, // MB
			MaxBackups: 7,
			MaxAge:     30,   // days
			Compress:   true, // gzip
		})

		consoleWriter := zapcore.Lock(os.Stdout)

		var gelfWriter zapcore.WriteSyncer
		if gelfAddr != "" {
			w, e := gelf.NewUDPWriter(gelfAddr) // "e.g:192.168.1.10:12201"
			if e != nil {
				return
			}
			gelfWriter = zapcore.AddSync(w)
		}

		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), consoleWriter, logLevel),
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), fileWriter, logLevel),
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), gelfWriter, logLevel),
		)

		instance = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	})

	return err
}

func GetLogger() *zap.Logger {
	if instance == nil {
		panic("logger not initialized, call init logging first")
	}
	return instance
}

type ctxKey struct{}

func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return instance
	}
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	}
	return instance
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	l := instance.With(zap.String("request_id", requestID))
	return context.WithValue(ctx, ctxKey{}, l)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	l := instance.With(zap.String("trace_id", traceID))
	return context.WithValue(ctx, ctxKey{}, l)
}

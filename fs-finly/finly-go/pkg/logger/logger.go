package logger

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

func init() {
	// Create the directory for storing log files
	logDir := "./logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic("Failed to create logs directory: " + err.Error())
	}

	// Configure log rotation
	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "finly.log"),
		MaxSize:    10,   // Maximum size of a log file in MB
		MaxBackups: 5,    // Maximum number of backup log files
		MaxAge:     30,   // Maximum age of log files in days
		Compress:   true, // Compress old log files
	}

	// Configure JSON encoder for structured logging
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		MessageKey:    "message",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeTime:    zapcore.ISO8601TimeEncoder, // Timestamp format
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	})

	// Set up file and console log outputs
	fileWriter := zapcore.AddSync(logFile)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Combine file and console outputs
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, fileWriter, zapcore.InfoLevel),
		zapcore.NewCore(encoder, consoleWriter, zapcore.InfoLevel),
	)

	// Initialize the logger
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Validate logger initialization
	if Log == nil {
		panic("Failed to initialize logger")
	}

	// Log a message to confirm logger setup
	Log.Info("Logger initialized", zap.String("time", time.Now().Format(time.RFC3339)))
}

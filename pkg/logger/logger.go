package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	Level       slog.Level
	Development bool
	Filename    string
}

var globalLogger *slog.Logger

func InitLogger(cfg Config) (*slog.Logger, error) {
	var logOutput io.Writer = os.Stdout

	if cfg.Filename != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.Filename), 0755); err != nil {
			return nil, err
		}

		logFile, err := os.OpenFile(cfg.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		logOutput = io.MultiWriter(os.Stdout, logFile)
	}

	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.Development,
	}

	var handler slog.Handler

	if cfg.Development {
		handler = slog.NewTextHandler(logOutput, opts)
	} else {
		handler = slog.NewJSONHandler(logOutput, opts)
	}

	globalLogger = slog.New(handler)
	slog.SetDefault(globalLogger)

	return globalLogger, nil
}

func GetLogger() *slog.Logger {
	if globalLogger == nil {
		InitLogger(Config{
			Level:       slog.LevelInfo,
			Development: true,
			Filename:    "logs/app.log",
		})
	}
	return globalLogger
}

func Close() error {
	return nil
}

func Info(msg string, args ...any) {
	globalLogger.Info(msg, args...)
}

func Error(msg string, err error, args ...any) {
	args = append(args, "error", err)
	globalLogger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	globalLogger.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	globalLogger.Warn(msg, args...)
}

package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/rampartfw/rampart/internal/config"
)

var GlobalLogger *slog.Logger

func Init(cfg *config.Config) error {
	var level slog.Level
	switch strings.ToLower(cfg.Logging.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var writer io.Writer = os.Stderr
	if cfg.Logging.File != "" {
		f, err := os.OpenFile(cfg.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writer = io.MultiWriter(os.Stderr, f)
	}

	var handler slog.Handler
	if cfg.Logging.Format == "json" {
		handler = slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level})
	}

	GlobalLogger = slog.New(handler)
	slog.SetDefault(GlobalLogger)
	return nil
}

func With(args ...interface{}) *slog.Logger {
	return GlobalLogger.With(args...)
}

func Component(name string) *slog.Logger {
	return GlobalLogger.With("component", name)
}

func Debug(msg string, args ...interface{}) {
	GlobalLogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	GlobalLogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	GlobalLogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	GlobalLogger.Error(msg, args...)
}

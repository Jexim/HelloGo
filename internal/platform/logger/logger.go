package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a zap logger with the provided level string (e.g., "debug", "info").
func New(levelString string) *zap.Logger {
	level := zapcore.InfoLevel
	if err := level.Set(levelString); err != nil {
		level = zapcore.InfoLevel
	}
	cfgZap := zap.NewProductionConfig()
	cfgZap.Level = zap.NewAtomicLevelAt(level)
	logger, _ := cfgZap.Build()
	return logger
}

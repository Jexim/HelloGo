package sentry

import (
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

// Init initializes Sentry with the given DSN
func Init(dsn string, logger *zap.Logger) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: 1.0,
		Environment:      "development", // TODO: make this configurable
		EnableTracing:    true,
	})
	if err != nil {
		logger.Error("sentry initialization failed", zap.Error(err))
		return err
	}

	logger.Info("sentry initialized successfully")
	return nil
}

// Flush flushes Sentry events
func Flush(timeout time.Duration) {
	sentry.Flush(timeout)
}

// Recover recovers from panics and reports them to Sentry
func Recover(logger *zap.Logger) {
	if err := recover(); err != nil {
		sentry.CurrentHub().Recover(err)
		logger.Error("recovered from panic", zap.Any("error", err))
	}
}

package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg *viper.Viper) *zap.Logger {
	level := zapcore.InfoLevel
	if err := level.Set(viper.GetString("logger.level")); err != nil {
		level = zapcore.InfoLevel
	}
	cfgZap := zap.NewProductionConfig()
	cfgZap.Level = zap.NewAtomicLevelAt(level)
	logger, _ := cfgZap.Build()
	return logger
}

func Middleware(log *zap.Logger, cfg *viper.Viper) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("request",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.Duration("duration", time.Since(start)),
		)
	}
}

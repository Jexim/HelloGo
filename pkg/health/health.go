package health

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Status struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]Status `json:"services"`
}

type Checker struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewChecker(db *gorm.DB, logger *zap.Logger) *Checker {
	return &Checker{
		db:     db,
		logger: logger,
	}
}

func (c *Checker) Check(ctx context.Context) HealthStatus {
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Services:  make(map[string]Status),
	}

	// Check database
	if err := c.checkDatabase(ctx); err != nil {
		status.Status = "degraded"
		status.Services["database"] = Status{
			Status:  "error",
			Message: err.Error(),
		}
	} else {
		status.Services["database"] = Status{
			Status: "ok",
		}
	}

	return status
}

func (c *Checker) checkDatabase(ctx context.Context) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

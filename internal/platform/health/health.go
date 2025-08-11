package health

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"
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
	db     *sql.DB
	logger *zap.Logger
}

func NewChecker(db *sql.DB, logger *zap.Logger) *Checker {
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
	return c.db.PingContext(ctx)
}

package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/Jexim/HelloGo/pkg/health"
)

type REST struct {
	checker *health.Checker
	logger  *zap.Logger
}

func New(r chi.Router, path string, checker *health.Checker, logger *zap.Logger) *REST {
	rest := &REST{
		checker: checker,
		logger:  logger,
	}

	r.Get(path, rest.Handler())

	return rest
}

func (r *REST) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
		defer cancel()

		status := r.checker.Check(ctx)

		// Set response status code
		if status.Status == "ok" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")

		// Encode response
		if err := json.NewEncoder(w).Encode(status); err != nil {
			r.logger.Error("failed to encode health check response", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

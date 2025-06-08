package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	TraceIDHeader = "X-Trace-ID"
)

// TraceID is a middleware that adds a trace ID to each request
func TraceID(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get trace ID from header or generate new one
			traceID := r.Header.Get(TraceIDHeader)
			if traceID == "" {
				traceID = uuid.New().String()
			}

			// Add trace ID to response headers
			w.Header().Set(TraceIDHeader, traceID)

			// Add trace ID to request context and logger
			ctx := r.Context()
			loggerWithTrace := logger.With(zap.String("trace_id", traceID))
			loggerWithTrace.Debug("processing request with trace ID")

			// Create new request with context containing logger
			r = r.WithContext(ctx)

			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}

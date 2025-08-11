package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	TraceIDHeader = "X-Trace-ID"
)

type contextKey string

const traceIDContextKey contextKey = "trace_id"

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
			ctx = contextWithTraceID(ctx, traceID)
			loggerWithTrace := logger.With(zap.String("trace_id", traceID))
			loggerWithTrace.Debug("processing request with trace ID")

			// Create new request with context containing trace id
			r = r.WithContext(ctx)

			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}

// contextWithTraceID sets trace id into context
func contextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDContextKey, traceID)
}

// GetTraceID extracts trace id from context if present
func GetTraceID(r *http.Request) string {
	if v := r.Context().Value(traceIDContextKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	// fallback to request header if present
	if h := r.Header.Get(TraceIDHeader); h != "" {
		return h
	}
	return ""
}

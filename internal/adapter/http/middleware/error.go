package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Jexim/HelloGo/internal/platform/apperr"
	"github.com/Jexim/HelloGo/internal/platform/sentry"
	"go.uber.org/zap"
)

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		TraceID string `json:"trace_id,omitempty"`
	} `json:"error"`
}

// ErrorHandler middleware standardizes error responses and logs them
func ErrorHandler(logger *zap.Logger, capture func(err error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap ResponseWriter to observe status codes
			rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			defer func() {
				if rec := recover(); rec != nil {
					// Convert panic to error and capture
					err, ok := rec.(error)
					if !ok {
						err = fmt.Errorf("panic: %v", rec)
					}
					if capture != nil {
						capture(err)
					} else {
						sentry.CaptureError(err)
					}
					respondError(rw, r, logger, http.StatusInternalServerError, "internal_error", "internal server error")
				}
			}()
			// Call next
			next.ServeHTTP(rw, r)

			// After next: capture unexpected 5xx responses
			if rw.status >= 500 {
				err := fmt.Errorf("http %d %s %s", rw.status, r.Method, r.URL.Path)
				if capture != nil {
					capture(err)
				} else {
					sentry.CaptureError(err)
				}
			}
		})
	}
}

func respondError(w http.ResponseWriter, r *http.Request, logger *zap.Logger, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	traceID := GetTraceID(r)
	resp := errorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	resp.Error.TraceID = traceID
	_ = json.NewEncoder(w).Encode(resp)
	logger.Error("http_error", zap.Int("status", status), zap.String("code", code), zap.String("message", message), zap.String("trace_id", traceID))
}

// MapError maps known errors to HTTP status and codes
func MapError(err error) (status int, code, msg string) {
	switch {
	case err == nil:
		return http.StatusOK, "", ""
	case errors.Is(err, apperr.ErrBadRequest):
		return http.StatusBadRequest, "bad_request", err.Error()
	case errors.Is(err, apperr.ErrNotFound):
		return http.StatusNotFound, "not_found", err.Error()
	case errors.Is(err, apperr.ErrAlreadyExists):
		return http.StatusConflict, "already_exists", err.Error()
	default:
		return http.StatusInternalServerError, "internal_error", "internal server error"
	}
}

// statusRecorder captures the status code written by handlers
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

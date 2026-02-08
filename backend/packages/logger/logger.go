package logger

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// New creates a new zap logger instance.
func New(isProduction bool) (*zap.Logger, error) {
	if isProduction {
		return zap.NewProduction()
	}
	// NewDevelopment builds a development logger that writes to standard error,
	// writes in a human-friendly format, and enables stacktraces.
	return zap.NewDevelopment()
}

// FromContext returns a logger instance with the request_id from the context baked in.
// If no request_id is found, it returns the original logger.
// This is the key function for creating context-aware loggers.
func FromContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	reqID := middleware.GetReqID(ctx)
	if reqID != "" {
		// With returns a new logger with the given field attached.
		return logger.With(zap.String("request_id", reqID))
	}
	return logger
}

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
	return zap.NewDevelopment()
}

// FromContext returns a logger instance with the request_id from the context baked in.
func FromContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	reqID := middleware.GetReqID(ctx)
	if reqID != "" {
		return logger.With(zap.String("request_id", reqID))
	}
	return logger
}

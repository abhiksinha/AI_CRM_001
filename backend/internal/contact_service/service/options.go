package service

import (
	"CRM/internal/contact_service/repo"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Option is a function that configures a ContactService.
type Option func(*ContactService)

// WithRepo creates a new repository from a DB connection and applies it to the service.
// This encapsulates the repository creation logic.
func WithRepo(db *gorm.DB) Option {
	return func(s *ContactService) {
		s.repo = repo.NewRepository(db)
	}
}

// WithLogger applies a logger to the service.
func WithLogger(logger *zap.Logger) Option {
	return func(s *ContactService) {
		s.logger = logger
	}
}

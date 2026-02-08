package service

import (
	"CRM/internal/user_service/repo"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Option func(*UserService)

func WithUserRepo(db *gorm.DB) Option {
	return func(s *UserService) {
		s.userRepo = repo.NewUserRepository(db)
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *UserService) {
		s.logger = logger
	}
}

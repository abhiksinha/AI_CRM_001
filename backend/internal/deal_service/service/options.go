package service

import (
	"CRM/internal/deal_service/repo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Option func(*DealService)

func WithDealRepo(db *gorm.DB) Option {
	return func(s *DealService) {
		s.dealRepo = repo.NewDealRepository(db)
	}
}

func WithTaskRepo(db *gorm.DB) Option {
	return func(s *DealService) {
		s.taskRepo = repo.NewTaskRepository(db)
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *DealService) {
		s.logger = logger
	}
}

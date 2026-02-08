package service

import (
	"CRM/internal/deal_service/contracts"
	"CRM/internal/deal_service/model"
	"CRM/internal/deal_service/repo"
	"CRM/packages/logger"
	"CRM/packages/public_response"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DealService struct {
	dealRepo *repo.DealRepository
	taskRepo *repo.TaskRepository
	logger   *zap.Logger
}

func NewDealService(opts ...Option) *DealService {
	s := &DealService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *DealService) CreateDeal(ctx context.Context, req contracts.CreateDealRequest) (*contracts.DealResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	closeDate, _ := time.Parse("2006-01-02", req.ExpectedCloseDate)
	newDeal := &model.Deal{
		Name:              req.Name,
		Stage:             req.Stage,
		Value:             req.Value,
		ExpectedCloseDate: closeDate,
		ContactID:         req.ContactID,
		OwnerID:           req.OwnerID,
	}
	err := s.dealRepo.ExecTxn(ctx, func(txnRepo *repo.DealRepository) error {
		return txnRepo.CreateDeal(newDeal)
	})
	if err != nil {
		log.Error("Error creating deal", zap.Error(err))
		return nil, err
	}
	log.Info("Successfully created deal", zap.String("deal_id", newDeal.ID))
	return s.dealRepo.ToDealApiResponse(newDeal), nil
}

func (s *DealService) GetDealByID(ctx context.Context, id, ownerID string) (*contracts.DealResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	deal, err := s.dealRepo.GetByIDAndOwner(id, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, public_response.ErrNotFound
		}
		log.Error("Error fetching deal", zap.Error(err))
		return nil, err
	}
	return s.dealRepo.ToDealApiResponse(deal), nil
}

func (s *DealService) ListDeals(ctx context.Context, ownerID string, req contracts.ListDealsRequest) (*contracts.ListDealsResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	deals, total, err := s.dealRepo.ListByOwner(ownerID, req)
	if err != nil {
		log.Error("Error listing deals", zap.Error(err))
		return nil, err
	}
	response := make([]contracts.DealResponse, len(deals))
	for i, d := range deals {
		response[i] = *s.dealRepo.ToDealApiResponse(&d)
	}
	return &contracts.ListDealsResponse{Data: response, TotalCount: total}, nil
}

func (s *DealService) UpdateDeal(ctx context.Context, id, ownerID string, req contracts.UpdateDealRequest) (*contracts.DealResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	deal, err := s.dealRepo.GetByIDAndOwner(id, ownerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, public_response.ErrNotFound
		}
		log.Error("Error fetching deal for update", zap.Error(err))
		return nil, err
	}

	if req.Name != nil {
		deal.Name = *req.Name
	}
	if req.Stage != nil {
		deal.Stage = *req.Stage
	}
	if req.Value != nil {
		deal.Value = *req.Value
	}
	if req.ExpectedCloseDate != nil {
		deal.ExpectedCloseDate, _ = time.Parse("2006-01-02", *req.ExpectedCloseDate)
	}

	err = s.dealRepo.ExecTxn(ctx, func(txnRepo *repo.DealRepository) error {
		return txnRepo.UpdateDeal(deal)
	})
	if err != nil {
		log.Error("Error during deal update transaction", zap.Error(err))
		return nil, err
	}
	return s.dealRepo.ToDealApiResponse(deal), nil
}

func (s *DealService) DeleteDeal(ctx context.Context, id, ownerID string) error {
	log := logger.FromContext(ctx, s.logger)
	err := s.dealRepo.ExecTxn(ctx, func(txnRepo *repo.DealRepository) error {
		rowsAffected, err := txnRepo.DeleteByIDAndOwner(id, ownerID)
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return public_response.ErrNotFound
		}
		return nil
	})
	if err != nil {
		log.Error("Error during deal deletion transaction", zap.Error(err))
		return err
	}
	log.Info("Successfully deleted deal", zap.String("id", id))
	return nil
}

func (s *DealService) CreateTask(ctx context.Context, dealID, ownerID string, req contracts.CreateTaskRequest) (*contracts.TaskResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	if _, err := s.dealRepo.GetByIDAndOwner(dealID, ownerID); err != nil {
		log.Warn("Attempted to add task to non-existent or unauthorized deal", zap.String("deal_id", dealID), zap.String("owner_id", ownerID))
		return nil, public_response.ErrNotFound
	}

	dueDate, _ := time.Parse("2006-01-02", req.DueDate)
	task := &model.Task{
		DealID:       dealID,
		AssignedToID: req.AssignedToID,
		Title:        req.Title,
		DueDate:      dueDate.Unix(),
		IsCompleted:  req.IsCompleted,
	}

	if err := s.taskRepo.CreateTask(task); err != nil {
		log.Error("Error creating task", zap.Error(err))
		return nil, err
	}

	return &contracts.TaskResponse{
		ID:           task.ID,
		Title:        task.Title,
		DueDate:      task.DueDate,
		IsCompleted:  task.IsCompleted,
		AssignedToID: task.AssignedToID,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
	}, nil
}

func (s *DealService) ListTasks(ctx context.Context, dealID, ownerID string) (*contracts.ListTasksResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	if _, err := s.dealRepo.GetByIDAndOwner(dealID, ownerID); err != nil {
		log.Warn("Attempted to list tasks for non-existent or unauthorized deal", zap.String("deal_id", dealID), zap.String("owner_id", ownerID))
		return nil, public_response.ErrNotFound
	}

	tasks, err := s.taskRepo.ListByDealID(dealID)
	if err != nil {
		log.Error("Error listing tasks", zap.Error(err))
		return nil, err
	}

	response := make([]contracts.TaskResponse, len(tasks))
	for i, t := range tasks {
		response[i] = contracts.TaskResponse{
			ID:           t.ID,
			Title:        t.Title,
			DueDate:      t.DueDate,
			IsCompleted:  t.IsCompleted,
			AssignedToID: t.AssignedToID,
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
		}
	}
	return &contracts.ListTasksResponse{Data: response, TotalCount: int64(len(response))}, nil
}

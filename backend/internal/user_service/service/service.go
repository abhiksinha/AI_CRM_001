package service

import (
	"CRM/internal/user_service/contracts"
	"CRM/internal/user_service/model"
	"CRM/internal/user_service/repo"
	"CRM/packages/logger"
	"CRM/packages/public_response"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repo.UserRepository
	logger   *zap.Logger
}

func NewUserService(opts ...Option) *UserService {
	s := &UserService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *UserService) VerifyPassword(ctx context.Context, req contracts.VerifyPasswordRequest) (bool, error) {
	log := logger.FromContext(ctx, s.logger)
	user, err := s.userRepo.GetByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found for password verification", zap.String("id", req.ID))
			return false, public_response.ErrNotFound
		}
		log.Error("Error fetching user for password verification", zap.Error(err))
		return false, err
	}
	// Convert Password type to string for comparison
	return user.CheckPassword(string(req.Password)), nil
}

func (s *UserService) CreateUser(ctx context.Context, req contracts.CreateUserRequest) (*contracts.UserResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	_, err := s.userRepo.GetByEmail(req.Email)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, public_response.ErrDuplicateEntry
	}
	newUser := &model.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: string(req.Password), // Convert Password type to string
		Role:         req.Role,
	}
	err = s.userRepo.ExecTxn(ctx, func(txnRepo *repo.UserRepository) error {
		return txnRepo.CreateUser(newUser)
	})
	if err != nil {
		log.Error("Error creating user", zap.Error(err))
		return nil, err
	}
	log.Info("Successfully created user", zap.String("user_id", newUser.ID))
	return toUserResponse(newUser), nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*contracts.UserResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found", zap.String("id", id))
			return nil, public_response.ErrNotFound
		}
		log.Error("Error fetching user", zap.Error(err))
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *UserService) ListUsers(ctx context.Context, req contracts.ListUsersRequest) (*contracts.ListUsersResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	users, total, err := s.userRepo.List(req)
	if err != nil {
		log.Error("Error listing users", zap.Error(err))
		return nil, err
	}
	response := make([]contracts.UserResponse, len(users))
	for i, u := range users {
		response[i] = *toUserResponse(&u)
	}
	return &contracts.ListUsersResponse{Data: response, TotalCount: total}, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req contracts.UpdateUserRequest) (*contracts.UserResponse, error) {
	log := logger.FromContext(ctx, s.logger)
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found for update", zap.String("id", id))
			return nil, public_response.ErrNotFound
		}
		log.Error("Error fetching user for update", zap.Error(err))
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	err = s.userRepo.ExecTxn(ctx, func(txnRepo *repo.UserRepository) error {
		return txnRepo.UpdateUser(user)
	})
	if err != nil {
		log.Error("Error during user update transaction", zap.Error(err))
		return nil, err
	}
	return toUserResponse(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	log := logger.FromContext(ctx, s.logger)
	err := s.userRepo.ExecTxn(ctx, func(txnRepo *repo.UserRepository) error {
		rowsAffected, err := txnRepo.DeleteByID(id)
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return public_response.ErrNotFound
		}
		return nil
	})
	if err != nil {
		log.Error("Error during user deletion transaction", zap.Error(err))
		return err
	}
	log.Info("Successfully deleted user", zap.String("id", id))
	return nil
}

func toUserResponse(user *model.User) *contracts.UserResponse {
	return &contracts.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

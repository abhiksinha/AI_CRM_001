package repo

import (
	"CRM/internal/user_service/contracts"
	"CRM/internal/user_service/model"
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ExecTxn(ctx context.Context, fn func(repo *UserRepository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(NewUserRepository(tx))
	})
}

func (r *UserRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "email = ?", email).Error
	return &user, err
}

func (r *UserRepository) List(req contracts.ListUsersRequest) ([]model.User, int64, error) {
	var users []model.User
	var totalCount int64

	query := r.db.Model(&model.User{})
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	query = query.Offset(offset).Limit(req.PageSize)

	err := query.Find(&users).Error
	return users, totalCount, err
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) DeleteByID(id string) (int64, error) {
	result := r.db.Delete(&model.User{}, "id = ?", id)
	return result.RowsAffected, result.Error
}

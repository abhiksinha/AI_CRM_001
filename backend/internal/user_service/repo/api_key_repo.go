package repo

import (
	"CRM/internal/user_service/model"

	"gorm.io/gorm"
)

type ApiKeyRepository struct {
	db *gorm.DB
}

func NewApiKeyRepository(db *gorm.DB) *ApiKeyRepository {
	return &ApiKeyRepository{db: db}
}

// Create inserts a new API key.
func (r *ApiKeyRepository) Create(apiKey *model.ApiKey) error {
	return r.db.Create(apiKey).Error
}

// GetByID retrieves an active API key by its ID.
func (r *ApiKeyRepository) GetByID(id string) (*model.ApiKey, error) {
	var apiKey model.ApiKey
	// GORM's soft delete automatically handles the "deleted_at IS NULL" condition.
	err := r.db.First(&apiKey, "id = ?", id).Error
	return &apiKey, err
}

// Expire soft deletes an API key.
func (r *ApiKeyRepository) Expire(apiKey *model.ApiKey) error {
	// GORM's Delete performs a soft delete if the model has a DeletedAt field.
	return r.db.Delete(apiKey).Error
}

package model

import (
	"CRM/packages/uniqueid"
	"gorm.io/gorm"
	"time"
)

// ApiKey represents the api_keys table in the database.
type ApiKey struct {
	ID        string         `gorm:"type:varchar(18);primary_key"`
	UserID    string         `gorm:"type:varchar(36);not null"`
	CreatedAt int64          `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // Use GORM's soft delete
}

func (a *ApiKey) TableName() string { return "api_keys" }

// BeforeCreate is a GORM hook that is called before a new record is created.
func (a *ApiKey) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		// Generate a new 18-character ID for the API key.
		a.ID = uniqueid.New()
	}
	a.CreatedAt = time.Now().Unix()
	return
}

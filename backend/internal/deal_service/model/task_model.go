package model

import (
	"CRM/packages/uniqueid"
	"time"

	"gorm.io/gorm"
)

// Task represents the tasks table in the database.
type Task struct {
	ID           string `gorm:"type:varchar(36);primary_key"`
	DealID       string `gorm:"type:varchar(36)"`
	AssignedToID string `gorm:"type:varchar(36)"`
	Title        string `gorm:"type:varchar(255);not null"`
	DueDate      int64
	IsCompleted  bool  `gorm:"default:false"`
	CreatedAt    int64 `gorm:"not null"`
	UpdatedAt    int64 `gorm:"not null"`
}

func (t *Task) TableName() string { return "tasks" }

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uniqueid.New()
	}
	now := time.Now().Unix()
	t.CreatedAt = now
	t.UpdatedAt = now
	return
}

func (t *Task) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now().Unix()
	return
}

package model

import (
	"CRM/packages/uniqueid"
	"time"

	"gorm.io/gorm"
)

// Deal represents the deals table in the database.
type Deal struct {
	ID                string    `gorm:"type:varchar(36);primary_key"`
	Name              string    `gorm:"type:varchar(255);not null"`
	Stage             string    `gorm:"type:varchar(100);not null"`
	Value             float64   `gorm:"type:decimal"`
	ExpectedCloseDate time.Time `gorm:"type:date"`
	ContactID         string    `gorm:"type:varchar(36)"`
	OwnerID           string    `gorm:"type:varchar(36)"`
	CreatedAt         int64     `gorm:"not null"`
	UpdatedAt         int64     `gorm:"not null"`
}

func (d *Deal) TableName() string { return "deals" }

func (d *Deal) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = uniqueid.New()
	}
	now := time.Now().Unix()
	d.CreatedAt = now
	d.UpdatedAt = now
	return
}

func (d *Deal) BeforeUpdate(tx *gorm.DB) (err error) {
	d.UpdatedAt = time.Now().Unix()
	return
}

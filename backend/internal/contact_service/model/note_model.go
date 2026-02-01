package model

import (
	"CRM/packages/uniqueid"
	"gorm.io/gorm"
	"time"
)

// Note represents the notes table in the database.
type Note struct {
	ID        string `gorm:"type:varchar(36);primary_key"`
	ContactID string `gorm:"type:varchar(36);not null"`
	AuthorID  string `gorm:"type:varchar(36);not null"`
	Content   string `gorm:"type:text"`
	CreatedAt int64  `gorm:"not null"`
}

// TableName explicitly sets the table name for the Note model.
func (n *Note) TableName() string {
	return "notes"
}

// BeforeCreate is a GORM hook that is called before a new record is created.
func (n *Note) BeforeCreate(tx *gorm.DB) (err error) {
	if n.ID == "" {
		n.ID = uniqueid.New()
	}
	n.CreatedAt = time.Now().Unix()
	return
}

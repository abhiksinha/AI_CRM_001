package model

import (
	"CRM/packages/uniqueid" // Import the new package
	"time"

	"gorm.io/gorm"
)

// Contact represents the contacts table in the database.
type Contact struct {
	ID        string `gorm:"type:varchar(36);primary_key"`
	FirstName string `gorm:"type:varchar(255)"`
	LastName  string `gorm:"type:varchar(255)"`
	Email     string `gorm:"type:varchar(255);unique"`
	Phone     string `gorm:"type:varchar(50)"`
	OwnerID   string `gorm:"type:varchar(36)"`
	CreatedAt int64  `gorm:"not null"`
	UpdatedAt int64  `gorm:"not null"`
}

// TableName explicitly sets the table name for the Contact model.
func (c *Contact) TableName() string {
	return "contacts"
}

// BeforeCreate is a GORM hook that is called before a new record is created.
func (c *Contact) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new 14-digit unique ID.
	if c.ID == "" {
		c.ID = uniqueid.New()
	}

	now := time.Now().Unix()
	c.CreatedAt = now
	c.UpdatedAt = now
	return
}

// BeforeUpdate is a GORM hook that is called before a record is updated.
func (c *Contact) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now().Unix()
	return
}

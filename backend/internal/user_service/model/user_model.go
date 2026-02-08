package model

import (
	"CRM/packages/uniqueid"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the users table in the database.
type User struct {
	ID           string `gorm:"type:varchar(36);primary_key"`
	FirstName    string `gorm:"type:varchar(255)"`
	LastName     string `gorm:"type:varchar(255)"`
	Email        string `gorm:"type:varchar(255);unique;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Role         string `gorm:"type:varchar(50);not null"`
	IsActive     bool   `gorm:"default:true"`
	CreatedAt    int64  `gorm:"not null"`
	UpdatedAt    int64  `gorm:"not null"`
}

func (u *User) TableName() string { return "users" }

// BeforeCreate is a GORM hook that hashes the password before creating a user.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uniqueid.New()
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	now := time.Now().Unix()
	u.CreatedAt = now
	u.UpdatedAt = now
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now().Unix()
	return
}

// CheckPassword compares a provided password with the stored hash.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// NewGormDB creates a new GORM database instance.
func NewGormDB(cfg Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Database connection established successfully.")
	return db, nil
}

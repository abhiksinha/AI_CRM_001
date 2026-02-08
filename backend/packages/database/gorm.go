package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// NewGormDB creates a new GORM database instance from a DBConfig.
func NewGormDB(cfg DBConfig) (*gorm.DB, error) {
	// Use the DSN method from the config struct.
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Database connection established successfully.")
	return db, nil
}

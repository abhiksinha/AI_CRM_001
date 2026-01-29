package main

import (
	"CRM/internal/contact_service"
	"CRM/internal/contact_service/repo"
	"CRM/internal/contact_service/service"
	"CRM/packages/database"
	"CRM/packages/logger" // Import the new logger package
	"CRM/packages/server"
	"log"

	"go.uber.org/zap"
)

func main() {
	// --- Logger Initialization ---
	// TODO: Use an env var to set production mode.
	appLogger, err := logger.New(false) // false = development mode
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer appLogger.Sync() // Flushes any buffered log entries

	// --- Database Connection ---
	dbConfig := database.Config{
		DSN: "host=localhost user=postgres password=postgres dbname=crm port=5432 sslmode=disable",
	}
	db, err := database.NewGormDB(dbConfig)
	if err != nil {
		appLogger.Fatal("failed to connect to database", zap.Error(err))
	}

	// --- Dependency Injection ---
	// 1. Create the repository.
	contactRepo := repo.NewRepository(db)

	// 2. Create the service, injecting the logger.
	contactSvc := service.NewContactService(contactRepo, appLogger)

	// --- Server Setup ---
	srv := server.New()

	// 3. Create the handler, injecting the service.
	contact_service.NewContactHandlerServer(srv.Router(), contactSvc)

	// Start the server.
	srv.Start(":8080")
}

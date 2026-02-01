package main

import (
	"CRM/internal/contact_service"
	"CRM/internal/contact_service/repo"
	"CRM/internal/contact_service/service"
	"CRM/packages/configloader"
	"CRM/packages/database"
	"CRM/packages/logger"
	"CRM/packages/server"
	"fmt"
	"go.uber.org/zap"
	"log"
)

func main() {
	// --- Load Configuration ---
	cfg, err := configloader.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// --- Logger Initialization ---
	isProduction := cfg.App.Env == "prod"
	appLogger, err := logger.New(isProduction)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	// --- Database Connection ---
	db, err := database.NewGormDB(cfg.Database)
	if err != nil {
		appLogger.Fatal("failed to connect to database", zap.Error(err))
	}

	// --- Dependency Injection ---
	contactRepo := repo.NewRepository(db)
	notesRepo := repo.NewNotesRepository(db)
	contactSvc := service.NewContactService(contactRepo, notesRepo, appLogger)

	// --- Server Setup ---
	srv := server.New()
	contact_service.NewContactHandlerServer(srv.Router(), contactSvc)

	// Start the server.
	serverPort := fmt.Sprintf(":%s", cfg.App.Port)
	srv.Start(serverPort)
}

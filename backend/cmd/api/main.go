package main

import (
	"CRM/internal/contact_service"
	"CRM/internal/contact_service/service"
	"CRM/packages/configloader"
	"CRM/packages/database"
	"CRM/packages/logger"
	"CRM/packages/server"
	"fmt"
	"log"

	"go.uber.org/zap"
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

	// --- Dependency Injection using Functional Options ---
	// 1. Create the service using the With... option functions.
	contactSvc := service.NewContactService(
		service.WithRepo(db),
		service.WithLogger(appLogger),
	)

	// --- Server Setup ---
	srv := server.New()

	// 2. Create the handler, injecting the fully configured service.
	contact_service.NewContactHandlerServer(srv.Router(), contactSvc)

	// Start the server.
	serverPort := fmt.Sprintf(":%s", cfg.App.Port)
	srv.Start(serverPort)
}

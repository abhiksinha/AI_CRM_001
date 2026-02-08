package main

import (
	contactServer "CRM/internal/contact_service"
	contactrepo "CRM/internal/contact_service/repo"
	contact_service "CRM/internal/contact_service/service"
	dealServer "CRM/internal/deal_service"
	deal_service "CRM/internal/deal_service/service"
	userServer "CRM/internal/user_service"
	user_service "CRM/internal/user_service/service"
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

	// --- Dependency Injection ---
	contactRepo := contactrepo.NewRepository(db)
	notesRepo := contactrepo.NewNotesRepository(db)
	contactSvc := contact_service.NewContactService(contactRepo, notesRepo, appLogger)

	dealSvc := deal_service.NewDealService(
		deal_service.WithDealRepo(db),
		deal_service.WithTaskRepo(db),
		deal_service.WithLogger(appLogger),
	)

	userSvc := user_service.NewUserService(
		user_service.WithUserRepo(db),
		user_service.WithLogger(appLogger),
	)

	// --- Server Setup ---
	srv := server.New()
	contactServer.NewContactHandlerServer(srv.Router(), contactSvc)
	dealServer.NewDealHandlerServer(srv.Router(), dealSvc)
	userServer.NewUserHandlerServer(srv.Router(), userSvc) // Register the new handler

	// Start the server.
	serverPort := fmt.Sprintf(":%s", cfg.App.Port)
	srv.Start(serverPort)
}

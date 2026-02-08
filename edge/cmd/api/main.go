package main

import (
	"context"
	edgeserver "edge/internal/edge_service"
	edgeservice "edge/internal/edge_service/service"
	"edge/packages/configloader"
	"edge/packages/logger"
	"edge/packages/redis"
	"edge/packages/server"
	"fmt"
	"log"

	"go.uber.org/zap"
)

func main() {
	cfg, err := configloader.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	isProduction := cfg.App.Env == "prod"
	appLogger, err := logger.New(isProduction)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer appLogger.Sync()

	redisClient := redis.NewClient(cfg.Redis)
	if err := redis.Ping(context.Background(), redisClient); err != nil {
		appLogger.Warn("failed to ping redis", zap.Error(err))
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			appLogger.Warn("failed to close redis", zap.Error(err))
		}
	}()

	authSvc := edgeservice.NewAuthService(redisClient, cfg.BackendServices["backend_service"], appLogger)
	proxySvc := edgeservice.NewProxyService(redisClient, appLogger)
	routeConfig := edgeserver.BuildRouteConfig()

	srv := server.New()
	edgeserver.NewEdgeHandlerServer(srv.Router(), authSvc, proxySvc, routeConfig, cfg.BackendServices)

	serverPort := fmt.Sprintf(":%s", cfg.App.Port)
	srv.Start(serverPort)
}

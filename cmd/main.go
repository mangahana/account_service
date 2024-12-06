package main

import (
	"account/internal/application"
	"account/internal/infrastructure/configuration"
	"account/internal/infrastructure/logger"
	"account/internal/infrastructure/postgresql"
	ban_repository "account/internal/infrastructure/repository/ban"
	code_repository "account/internal/infrastructure/repository/code"
	session_repository "account/internal/infrastructure/repository/session"
	user_repository "account/internal/infrastructure/repository/user"
	ban_service "account/internal/service/ban"
	code_service "account/internal/service/code"
	session_service "account/internal/service/session"
	user_service "account/internal/service/user"
	"account/internal/transport/grpc"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := configuration.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	db, err := postgresql.New(&cfg.DB)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// repositories
	userRepository := user_repository.New(db)
	banRepository := ban_repository.New(db)
	codeRepository := code_repository.New(db)
	sessionRepository := session_repository.New(db)

	// services
	userService := user_service.New(userRepository)
	banService := ban_service.New(banRepository)
	codeService := code_service.New(&cfg.SMS, codeRepository)
	sessionService := session_service.New(sessionRepository)

	useCase := application.New(logger, userService, banService, codeService, sessionService)

	grpcServer := grpc.New(useCase)

	go func() {
		if err := grpcServer.Run(cfg.Server.GrpcSocker); err != nil {
			logger.Fatal("grpc server down", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	_, shutdown := context.WithTimeout(context.Background(), time.Second*3)
	defer shutdown()

	grpcServer.Stop()
}

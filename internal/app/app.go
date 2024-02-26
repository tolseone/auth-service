package app

import (
	"log/slog"
	"time"

	grpcapp "auth-service/internal/app/grpc"
	"auth-service/internal/services/auth"
	db "auth-service/internal/storage/postgresql"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, gprcPort int, tokenTTL time.Duration) *App {
	storage := db.New(log)
	if storage == nil {
		panic("Failed to create storage")
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, gprcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}

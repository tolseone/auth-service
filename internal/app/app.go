package app

import (
	"log/slog"
	"time"

	"auth-service/internal/app/grpc"

)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, gprcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: init storage

	// TODO: init auth service

	grpcApp := grpcapp.New(log, gprcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}

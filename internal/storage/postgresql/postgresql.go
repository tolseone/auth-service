package db

import (
	"context"
	"log/slog"
	"strings"

	_ "github.com/jackc/pgx/v5/pgconn"

	"auth-service/internal/config"
	"auth-service/internal/domain/models"
	"auth-service/pkg/client/postgresql"

)

type Storage struct {
	client postgresql.Client
	log    *slog.Logger
}

func New(log *slog.Logger) *Storage {
	cfg := config.MustLoad()
	client, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		log.Info("Failed to connect to PostgreSQL: %v", err)
		return &Storage{}
	}
	log.Info("connected to PostgreSQL")

	return &Storage{
		client: client,
		log:    log,
	}
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "Storage.SaveUser"
	return 0, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	return models.User{}, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	return false, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	return models.App{}, nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

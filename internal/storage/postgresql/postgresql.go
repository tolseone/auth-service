package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"auth-service/internal/config"
	"auth-service/internal/domain/models"
	"auth-service/internal/lib/logger/sl"
	"auth-service/internal/storage"
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

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uuid.UUID, error) {
	const op = "Storage.SaveUser"

	q := `
		INSERT INTO users (
			id, 
			email,
			password
		) 
		VALUES (
			gen_random_uuid(), 
			$1, 
			$2
		)
		RETURNING id
	`
	s.log.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var id string

	if err := s.client.QueryRow(ctx, q, email, passHash).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s, op: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState(), op))
			return uuid.Nil, newErr
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)

	}
	UID, err := uuid.Parse(id)
	if err != nil {
		s.log.Error("failed to parse UUID", sl.Err(err))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	s.log.Info("Completed to create user")

	return UID, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	q := `
        SELECT 
			id,  
			email,
			password
		FROM users 
		WHERE 
			email = $1
	`
	s.log.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var user models.User

	if err := s.client.QueryRow(ctx, q, email).Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	q := `
        SELECT 
			is_admin
		FROM users
		WHERE 
			id = $1
	`
	s.log.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var isAdmin bool

	if err := s.client.QueryRow(ctx, q, userID).Scan(&isAdmin); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.sqlite.App"

	q := `
        SELECT 
			id,  
			name,
			secret
		FROM apps
		WHERE 
			id = $1
	`
	s.log.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var app models.App

	if err := s.client.QueryRow(ctx, q, appID).Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

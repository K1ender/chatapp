package database

import (
	"chatapp/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn, err := getDSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("get dsn: %w", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	return pool, nil
}

func getDSN(cfg config.DatabaseConfig) (string, error) {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name,
	), nil
}

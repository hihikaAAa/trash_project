// Package database
package database

import (
	"context"
	"fmt"

	"github.com/hihikaAAa/trash_project/pkg/config"
	"github.com/hihikaAAa/trash_project/pkg/database/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(ctx context.Context, cfg *config.Configuration) (*pgxpool.Pool, error) {
	pool, err := postgres.ConnectDB(ctx, cfg.PostgreSQL)
	if err != nil {
		return nil, fmt.Errorf("init db: %w", err)
	}
	return pool, nil
}

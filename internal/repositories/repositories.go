// Package repositories contains repository interfaces and constructors.
package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/hihikaAAa/trash_project/internal/domain/task"
	"github.com/hihikaAAa/trash_project/internal/repositories/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Orders interface {
	Create(ctx context.Context, order *task.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error)
	ListByClientID(ctx context.Context, clientID uuid.UUID) ([]*task.Task, error)
	ListByWorkerID(ctx context.Context, workerID uuid.UUID) ([]*task.Task, error)
	ListAvailable(ctx context.Context) ([]*task.Task, error)
	ListAll(ctx context.Context) ([]*task.Task, error)
	Update(ctx context.Context, order *task.Task) error
}

type Repository struct {
	Orders Orders
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Orders: postgres.NewTaskRepository(db),
	}
}

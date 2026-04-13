// Package service wires use-cases.
package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hihikaAAa/trash_project/internal/domain/task"
	"github.com/hihikaAAa/trash_project/internal/repositories"
	"github.com/hihikaAAa/trash_project/internal/service/services"
)

type Orders interface {
	Create(ctx context.Context, actor services.Actor, input services.CreateOrderInput) (*task.Task, error)
	GetByID(ctx context.Context, actor services.Actor, orderID uuid.UUID) (*task.Task, error)
	ListOwn(ctx context.Context, actor services.Actor) ([]*task.Task, error)
	ListAvailable(ctx context.Context, actor services.Actor) ([]*task.Task, error)
	ListAssigned(ctx context.Context, actor services.Actor) ([]*task.Task, error)
	Accept(ctx context.Context, actor services.Actor, orderID uuid.UUID) (*task.Task, error)
	Complete(ctx context.Context, actor services.Actor, orderID uuid.UUID) (*task.Task, error)
	ListAll(ctx context.Context, actor services.Actor) ([]*task.Task, error)
	Assign(ctx context.Context, actor services.Actor, orderID, workerID uuid.UUID) (*task.Task, error)
	Cancel(ctx context.Context, actor services.Actor, orderID uuid.UUID) (*task.Task, error)
}

type Service struct {
	Orders Orders
}

func NewService(repository *repositories.Repository) *Service {
	return &Service{
		Orders: services.NewOrdersService(repository.Orders),
	}
}

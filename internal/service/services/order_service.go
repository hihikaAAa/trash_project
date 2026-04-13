package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hihikaAAa/trash_project/internal/domain/task"
	domainerrors "github.com/hihikaAAa/trash_project/internal/domainerrors"
)

type OrdersRepository interface {
	Create(ctx context.Context, order *task.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error)
	ListByClientID(ctx context.Context, clientID uuid.UUID) ([]*task.Task, error)
	ListByWorkerID(ctx context.Context, workerID uuid.UUID) ([]*task.Task, error)
	ListAvailable(ctx context.Context) ([]*task.Task, error)
	ListAll(ctx context.Context) ([]*task.Task, error)
	Update(ctx context.Context, order *task.Task) error
}

type OrdersService struct {
	repo OrdersRepository
}

type Actor struct {
	ID   uuid.UUID
	Role task.Role
}

type CreateOrderInput struct {
	Address       string
	Description   *string
	PreferredTime *time.Time
}

func NewOrdersService(repo OrdersRepository) *OrdersService {
	return &OrdersService{repo: repo}
}

func (s *OrdersService) Create(ctx context.Context, actor Actor, input CreateOrderInput) (*task.Task, error) {
	if actor.Role != task.RoleUser {
		return nil, domainerrors.ErrForbidden
	}

	order, err := task.NewTask(actor.ID, input.Address, input.Description, input.PreferredTime, time.Now().UTC(), actor.Role)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrdersService) GetByID(ctx context.Context, actor Actor, orderID uuid.UUID) (*task.Task, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if !order.CanBeViewedBy(actor.ID, actor.Role) {
		return nil, domainerrors.ErrForbidden
	}
	return order, nil
}

func (s *OrdersService) ListOwn(ctx context.Context, actor Actor) ([]*task.Task, error) {
	if actor.Role != task.RoleUser {
		return nil, domainerrors.ErrForbidden
	}
	return s.repo.ListByClientID(ctx, actor.ID)
}

func (s *OrdersService) ListAvailable(ctx context.Context, actor Actor) ([]*task.Task, error) {
	if actor.Role != task.RoleWorker {
		return nil, domainerrors.ErrForbidden
	}
	return s.repo.ListAvailable(ctx)
}

func (s *OrdersService) ListAssigned(ctx context.Context, actor Actor) ([]*task.Task, error) {
	if actor.Role != task.RoleWorker {
		return nil, domainerrors.ErrForbidden
	}
	return s.repo.ListByWorkerID(ctx, actor.ID)
}

func (s *OrdersService) Accept(ctx context.Context, actor Actor, orderID uuid.UUID) (*task.Task, error) {
	if actor.Role != task.RoleWorker {
		return nil, domainerrors.ErrForbidden
	}

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	switch order.Status {
	case task.StatusCreated:
		if err = order.AssignWorker(actor.ID, now); err != nil {
			return nil, err
		}
		if err = order.StartByWorker(actor.ID, now, actor.Role); err != nil {
			return nil, err
		}
	case task.StatusAssigned:
		if err = order.StartByWorker(actor.ID, now, actor.Role); err != nil {
			return nil, err
		}
	default:
		return nil, domainerrors.ErrInvalidStatusTransition
	}

	if err = s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrdersService) Complete(ctx context.Context, actor Actor, orderID uuid.UUID) (*task.Task, error) {
	if actor.Role != task.RoleWorker {
		return nil, domainerrors.ErrForbidden
	}

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if err = order.CompleteByWorker(actor.ID, time.Now().UTC(), actor.Role); err != nil {
		return nil, err
	}
	if err = s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrdersService) ListAll(ctx context.Context, actor Actor) ([]*task.Task, error) {
	if actor.Role != task.RoleAdmin {
		return nil, domainerrors.ErrForbidden
	}
	return s.repo.ListAll(ctx)
}

func (s *OrdersService) Assign(ctx context.Context, actor Actor, orderID, workerID uuid.UUID) (*task.Task, error) {
	if actor.Role != task.RoleAdmin {
		return nil, domainerrors.ErrForbidden
	}

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if err = order.AssignWorker(workerID, time.Now().UTC()); err != nil {
		return nil, err
	}
	if err = s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrdersService) Cancel(ctx context.Context, actor Actor, orderID uuid.UUID) (*task.Task, error) {
	if actor.Role != task.RoleAdmin {
		return nil, domainerrors.ErrForbidden
	}

	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if err = order.CancelByAdmin(time.Now().UTC(), actor.Role); err != nil {
		return nil, err
	}
	if err = s.repo.Update(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func ParseRole(role string) (task.Role, error) {
	switch task.Role(role) {
	case task.RoleUser, task.RoleWorker, task.RoleAdmin:
		return task.Role(role), nil
	default:
		return "", domainerrors.ErrWrongRole
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, domainerrors.ErrTaskNotFound)
}

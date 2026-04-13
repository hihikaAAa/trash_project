// Package postgres contains postgres implementations for repositories.
package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hihikaAAa/trash_project/internal/domain/task"
	domainerrors "github.com/hihikaAAa/trash_project/internal/domainerrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository interface {
	Create(ctx context.Context, order *task.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error)
	ListByClientID(ctx context.Context, clientID uuid.UUID) ([]*task.Task, error)
	ListByWorkerID(ctx context.Context, workerID uuid.UUID) ([]*task.Task, error)
	ListAvailable(ctx context.Context) ([]*task.Task, error)
	ListAll(ctx context.Context) ([]*task.Task, error)
	Update(ctx context.Context, order *task.Task) error
}

type taskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, order *task.Task) error {
	const op = "internal.repositories.postgres.task_repo.Create"

	const q = `
	INSERT INTO orders(
		id, user_id, worker_id, address, description, preferred_time,
		status, created_at, updated_at, assigned_at, started_at, completed_at, canceled_at
	)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`

	_, err := r.db.Exec(
		ctx,
		q,
		order.ID,
		order.ClientID,
		order.WorkerID,
		order.Address,
		order.Description,
		order.PreferredTime,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
		order.AssignedAt,
		order.StartedAt,
		order.CompletedAt,
		order.CanceledAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	const op = "internal.repositories.postgres.task_repo.GetByID"

	const q = `
	SELECT id, user_id, worker_id, address, description, preferred_time,
	       status, created_at, updated_at, assigned_at, started_at, completed_at, canceled_at
	FROM orders
	WHERE id = $1
	`

	order, err := scanOrderRow(r.db.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainerrors.ErrTaskNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return order, nil
}

func (r *taskRepository) ListByClientID(ctx context.Context, clientID uuid.UUID) ([]*task.Task, error) {
	const op = "internal.repositories.postgres.task_repo.ListByClientID"

	const q = `
	SELECT id, user_id, worker_id, address, description, preferred_time,
	       status, created_at, updated_at, assigned_at, started_at, completed_at, canceled_at
	FROM orders
	WHERE user_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, q, clientID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return scanOrders(rows)
}

func (r *taskRepository) ListByWorkerID(ctx context.Context, workerID uuid.UUID) ([]*task.Task, error) {
	const op = "internal.repositories.postgres.task_repo.ListByWorkerID"

	const q = `
	SELECT id, user_id, worker_id, address, description, preferred_time,
	       status, created_at, updated_at, assigned_at, started_at, completed_at, canceled_at
	FROM orders
	WHERE worker_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, q, workerID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return scanOrders(rows)
}

func (r *taskRepository) ListAvailable(ctx context.Context) ([]*task.Task, error) {
	const op = "internal.repositories.postgres.task_repo.ListAvailable"

	const q = `
	SELECT id, user_id, worker_id, address, description, preferred_time,
	       status, created_at, updated_at, assigned_at, started_at, completed_at, canceled_at
	FROM orders
	WHERE status = $1
	ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, q, task.StatusCreated)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return scanOrders(rows)
}

func (r *taskRepository) ListAll(ctx context.Context) ([]*task.Task, error) {
	const op = "internal.repositories.postgres.task_repo.ListAll"

	const q = `
	SELECT id, user_id, worker_id, address, description, preferred_time,
	       status, created_at, updated_at, assigned_at, started_at, completed_at, canceled_at
	FROM orders
	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	return scanOrders(rows)
}

func (r *taskRepository) Update(ctx context.Context, order *task.Task) error {
	const op = "internal.repositories.postgres.task_repo.Update"

	const q = `
	UPDATE orders
	SET user_id = $2,
	    worker_id = $3,
	    address = $4,
	    description = $5,
	    preferred_time = $6,
	    status = $7,
	    updated_at = $8,
	    assigned_at = $9,
	    started_at = $10,
	    completed_at = $11,
	    canceled_at = $12
	WHERE id = $1
	`

	res, err := r.db.Exec(
		ctx,
		q,
		order.ID,
		order.ClientID,
		order.WorkerID,
		order.Address,
		order.Description,
		order.PreferredTime,
		order.Status,
		order.UpdatedAt,
		order.AssignedAt,
		order.StartedAt,
		order.CompletedAt,
		order.CanceledAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if res.RowsAffected() == 0 {
		return domainerrors.ErrTaskNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanOrderRow(row scanner) (*task.Task, error) {
	order := &task.Task{}

	var workerID pgtype.UUID
	var description pgtype.Text
	var preferredTime pgtype.Timestamptz
	var assignedAt pgtype.Timestamptz
	var startedAt pgtype.Timestamptz
	var completedAt pgtype.Timestamptz
	var canceledAt pgtype.Timestamptz

	err := row.Scan(
		&order.ID,
		&order.ClientID,
		&workerID,
		&order.Address,
		&description,
		&preferredTime,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
		&assignedAt,
		&startedAt,
		&completedAt,
		&canceledAt,
	)
	if err != nil {
		return nil, err
	}

	if workerID.Valid {
		id := uuid.UUID(workerID.Bytes)
		order.WorkerID = &id
	}
	if description.Valid {
		v := description.String
		order.Description = &v
	}
	if preferredTime.Valid {
		v := preferredTime.Time
		order.PreferredTime = &v
	}
	if assignedAt.Valid {
		v := assignedAt.Time
		order.AssignedAt = &v
	}
	if startedAt.Valid {
		v := startedAt.Time
		order.StartedAt = &v
	}
	if completedAt.Valid {
		v := completedAt.Time
		order.CompletedAt = &v
	}
	if canceledAt.Valid {
		v := canceledAt.Time
		order.CanceledAt = &v
	}

	return order, nil
}

func scanOrders(rows pgx.Rows) ([]*task.Task, error) {
	orders := make([]*task.Task, 0)
	for rows.Next() {
		order, err := scanOrderRow(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

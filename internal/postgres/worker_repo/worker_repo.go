// Package workerrepo
package workerrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
	"github.com/hihikaAAa/TrashProject/internal/domain/worker"
	postgreserrors "github.com/hihikaAAa/TrashProject/internal/postgres/postgres_errors"
)

type WorkerRepository interface {
	AddWorker(ctx context.Context, worker *worker.Worker) error
	SetIsActive(ctx context.Context, id uuid.UUID, active bool) (*worker.Worker, error)
	FindActive(ctx context.Context) ([]*worker.Worker, error)
	GetByID(ctx context.Context, id uuid.UUID) (*worker.Worker, error)
	List(ctx context.Context) ([]*worker.Worker, error)
	DeleteWorker(ctx context.Context, id uuid.UUID) error
	UpdateWorker(ctx context.Context, w *worker.Worker) (*worker.Worker, error)
}
type workerRepository struct {
	db *sql.DB
}

func NewWorkerRepository(db *sql.DB) WorkerRepository {
	return &workerRepository{db: db}
}

func (r *workerRepository) AddWorker(ctx context.Context, worker *worker.Worker) error {
	const op = "internal.postgres.worker_repo.AddWorker"

	const q = `
	INSERT INTO workers(worker_id, account_id, first_name, surname, last_name, is_active)
	VALUES ($1, $2, $3, $4, $5, $6) 
	`

	_, err := r.db.ExecContext(ctx, q, worker.ID, worker.AccountID ,worker.Person.FirstName, worker.Person.Surname, worker.Person.LastName, false)
	if err != nil {
		return fmt.Errorf("%s, ExecContext: %w", op, err)
	}
	return nil
}

func (r *workerRepository) SetIsActive(ctx context.Context, id uuid.UUID, active bool) (*worker.Worker, error) {
	const op = "internal.postgres.worker_repo.SetIsActive"

	const q = `
	UPDATE workers 
	SET is_active = $2, updated_at = now()
	WHERE worker_id = $1
	RETURNING worker_id, account_id, first_name, surname, last_name, is_active
	`

	w := &worker.Worker{Person: &person.Person{}}
	err := r.db.QueryRowContext(ctx, q, id, active).Scan(&w.ID,&w.AccountID, &w.Person.FirstName, &w.Person.Surname, &w.Person.LastName, &w.IsActive)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s, %w", op, postgreserrors.ErrWorkerNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return w, nil
}

func (r *workerRepository) FindActive(ctx context.Context) ([]*worker.Worker, error) {
	const op = "internal.postgres.worker_repo.FindActive"

	const q = `
	SELECT worker_id, account_id, first_name, surname, is_active
	FROM workers 
	WHERE is_active = $1
	`

	rows, err := r.db.QueryContext(ctx, q, true)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	activeWorkers := make([]*worker.Worker, 0)
	for rows.Next() {
		w := &worker.Worker{Person: &person.Person{}}
		err := rows.Scan(&w.ID, &w.AccountID, &w.Person.FirstName, &w.Person.Surname, &w.IsActive)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		activeWorkers = append(activeWorkers, w)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}
	return activeWorkers, nil
}

func (r *workerRepository) GetByID(ctx context.Context, id uuid.UUID) (*worker.Worker, error) {
	const op = "internal.postgres.worker_repo.GetByID"

	const q = `
	SELECT worker_id, account_id, first_name, surname, last_name, is_active
	FROM workers
	WHERE worker_id = $1
	`

	w := &worker.Worker{Person: &person.Person{}}
	err := r.db.QueryRowContext(ctx, q, id).Scan(&w.ID, &w.AccountID, &w.Person.FirstName, &w.Person.Surname, &w.Person.LastName, &w.IsActive)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: %w", op, postgreserrors.ErrWorkerNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return w, nil
}

func (r *workerRepository) List(ctx context.Context) ([]*worker.Worker, error) {
	const op = "internal.postgres.worker_repo.List"

	const q = `
	SELECT worker_id, account_id, first_name, surname, last_name, is_active
	FROM workers
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	workers := make([]*worker.Worker, 0)

	for rows.Next() {
		w := &worker.Worker{Person: &person.Person{}}
		err := rows.Scan(&w.ID, &w.AccountID, &w.Person.FirstName, &w.Person.Surname, &w.Person.LastName, &w.IsActive)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		workers = append(workers, w)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return workers, nil
}

func (r *workerRepository) DeleteWorker(ctx context.Context, id uuid.UUID) error {
	const op = "internal.postgres.worker_repo.DeleteWorker"

	const q = `
		DELETE FROM workers
		WHERE worker_id = $1
	`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("%s, ExecContext: %w", op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s, RowsAffected: %w", op, err)
	}
	if affected == 0 {
		return fmt.Errorf("%s: %w", op, postgreserrors.ErrWorkerNotFound)
	}

	return nil
}

func (r *workerRepository) UpdateWorker(ctx context.Context, w *worker.Worker) (*worker.Worker, error) {
	const op = "internal.postgres.worker_repo.UpdateWorker"

	const q = `
	UPDATE workers
	SET account_id = $2, first_name = $3, surname = $4, last_name = $5, is_active = $6, updated_at = now()
	WHERE worker_id = $1
	RETURNING worker_id, account_id, first_name, surname, last_name, is_active
	`

	worker := &worker.Worker{Person: &person.Person{}}
	err := r.db.QueryRowContext(ctx, q, w.ID, w.AccountID, w.Person.FirstName, w.Person.Surname, w.Person.LastName, w.IsActive).Scan(
		&worker.ID, &worker.AccountID, &worker.Person.FirstName, &worker.Person.Surname, &worker.Person.LastName, &worker.IsActive,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s, %w", op, postgreserrors.ErrWorkerNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return worker, nil
}

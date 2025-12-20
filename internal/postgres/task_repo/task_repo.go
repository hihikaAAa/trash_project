package taskrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hihikaAAa/TrashProject/internal/domain/task"
	repoerrors "github.com/hihikaAAa/TrashProject/internal/repository/postgres/repo_errors"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) AddTask(ctx context.Context, task *task.Task) error {
	const op = "internal.repository.postgres.task_repo.AddTask"

	const q = `
	INSERT INTO tasks(task_id, client_id, address_id, worker_id, status)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, q, task.ID, task.ClientID, task.AddressID, task.WorkerID, task.Status)
	if err != nil {
		return fmt.Errorf("%s, ExecContext: %w", op, err)
	}
	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	const op = "internal.repository.postgres.task_repo.GetByID"

	const q = `
	SELECT task_id, client_id, address_id, worker_id, status
	FROM tasks
	WHERE task_id = $1
	`

	t := &task.Task{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.ClientID, &t.AddressID, &t.WorkerID, &t.Status)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: %w", op, repoerrors.ErrTaskNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}

	return t, nil
}

func (r *TaskRepository) ListByClientID(ctx context.Context, clientID uuid.UUID) ([]*task.Task, error) {
	const op = "internal.repository.postgres.task_repo.ListByClientID"

	const q = `
	SELECT task_id, client_id, address_id, worker_id, status
	FROM tasks
	WHERE client_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, q, clientID)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	tasks := make([]*task.Task, 0)
	for rows.Next() {
		t := &task.Task{}
		err := rows.Scan(&t.ID, &t.ClientID, &t.AddressID, &t.WorkerID, &t.Status)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return tasks, nil
}

func (r *TaskRepository) ListActiveByWorkerID(ctx context.Context, workerID uuid.UUID) ([]*task.Task, error) {
	const op = "internal.repository.postgres.task_repo.ListActiveByWorkerID"

	const q = `
	SELECT task_id, client_id, address_id, worker_id, status
	FROM tasks
	WHERE worker_id = $1 AND status = $2
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, q, workerID, task.StatusInProgress)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	tasks := make([]*task.Task, 0)
	for rows.Next() {
		t := &task.Task{}
		err := rows.Scan(&t.ID, &t.ClientID, &t.AddressID, &t.WorkerID, &t.Status)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return tasks, nil
}

func (r *TaskRepository) ListDoneByWorkerID(ctx context.Context, workerID uuid.UUID) ([]*task.Task, error) {
	const op = "internal.repository.postgres.task_repo.ListDoneByWorkerID"

	const q = `
	SELECT task_id, client_id, address_id, worker_id, status
	FROM tasks
	WHERE worker_id = $1 AND status = $2
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, q, workerID, task.StatusDone)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	tasks := make([]*task.Task, 0)
	for rows.Next() {
		t := &task.Task{}
		err := rows.Scan(&t.ID, &t.ClientID, &t.AddressID, &t.WorkerID, &t.Status)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return tasks, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	const op = "internal.repository.postgres.task_repo.DeleteTask"

	const q = `
	DELETE FROM tasks
	WHERE task_id = $1
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
		return fmt.Errorf("%s: %w", op, repoerrors.ErrTaskNotFound)
	}

	return nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status task.Status) (*task.Task, error) {
	const op = "internal.repository.postgres.task_repo.UpdateStatus"

	const q = `
	UPDATE tasks
	SET status = $2, updated_at = now()
	WHERE task_id = $1
	RETURNING task_id, client_id, address_id, worker_id, status
	`

	task := &task.Task{}
	err := r.db.QueryRowContext(ctx, q, id, status).Scan(
		&task.ID, &task.ClientID, &task.AddressID, &task.WorkerID, &task.Status,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s, %w", op, repoerrors.ErrTaskNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return task, nil
}

func (r *TaskRepository) HasOpenTaskForClient(ctx context.Context, clientID uuid.UUID) (bool, error) {
	const op = "internal.repository.postgres.task_repo.HasOpenTaskForClient"

	const q = `
	SELECT 1 
	FROM tasks
	WHERE client_id = $1 AND status IN ($2, $3)
	LIMIT 1
	`

	var dummy int
	err := r.db.QueryRowContext(ctx, q, clientID, task.StatusOpen, task.StatusInProgress).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return true, nil
}

func (r *TaskRepository) AssignWorker(ctx context.Context, taskID uuid.UUID, workerID uuid.UUID) (*task.Task, error) {
	const op = "internal.repository.postgres.task_repo.AssignWorker"

	const q = `
	UPDATE tasks
	SET worker_id = $2, status = $3, updated_at = now()
	WHERE task_id = $1
	RETURNING task_id, client_id, address_id, worker_id, status
	`

	t := &task.Task{}
	err := r.db.QueryRowContext(ctx, q, taskID, workerID, task.StatusInProgress).Scan(
		&t.ID, &t.ClientID, &t.AddressID, &t.WorkerID, &t.Status,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s, %w", op, repoerrors.ErrTaskNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return t, nil
}

func (r *TaskRepository) ListOpenTasks(ctx context.Context) ([]*task.Task, error) {
	const op = "internal.repository.postgres.task_repo.ListOpenTasks"

	const q = `
	SELECT task_id, client_id, address_id, worker_id, status
	FROM tasks
	WHERE status = $1
	`

	rows, err := r.db.QueryContext(ctx, q, task.StatusOpen)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	tasks := make([]*task.Task, 0)

	for rows.Next() {
		t := &task.Task{}
		err := rows.Scan(&t.ID, &t.ClientID, &t.AddressID, &t.WorkerID, &t.Status)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		tasks = append(tasks, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return tasks, nil
}

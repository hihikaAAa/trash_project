package taskrepo

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/task"
	postgreserrors "github.com/hihikaAAa/TrashProject/internal/postgres/postgres_errors"
)

func newTestTaskRepo(t *testing.T) (TaskRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	repo := NewTaskRepository(db)
	cleanup := func() { _ = db.Close() }

	return repo, mock, cleanup
}

func TestTaskRepository_AddTask_Success(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	tsk := &task.Task{
		ID:        uuid.New(),
		ClientID:  uuid.New(),
		AddressID: uuid.New(),
		Status:    task.StatusOpen,
	}

	mock.ExpectExec("INSERT INTO tasks").
		WithArgs(tsk.ID, tsk.ClientID, tsk.AddressID, tsk.Status).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.AddTask(ctx, tsk); err != nil {
		t.Fatalf("AddTask error: %v", err)
	}
}

func TestTaskRepository_GetByID_Success(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	rows := sqlmock.NewRows([]string{
		"task_id", "client_id", "address_id", "worker_id", "status",
	}).AddRow(id, uuid.New(), uuid.New(), uuid.New(), task.StatusOpen)

	mock.ExpectQuery("FROM tasks").
		WithArgs(id).
		WillReturnRows(rows)

	tsk, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID error: %v", err)
	}

	if tsk.ID != id {
		t.Fatalf("unexpected task id")
	}
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery("FROM tasks").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(ctx, id)
	if !errors.Is(err, postgreserrors.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskRepository_ListByClientID(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	clientID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"task_id", "client_id", "address_id", "worker_id", "status",
	}).
		AddRow(uuid.New(), clientID, uuid.New(), uuid.New(), task.StatusOpen).
		AddRow(uuid.New(), clientID, uuid.New(), uuid.New(), task.StatusDone)

	mock.ExpectQuery("FROM tasks").
		WithArgs(clientID).
		WillReturnRows(rows)

	tasks, err := repo.ListByClientID(ctx, clientID)
	if err != nil {
		t.Fatalf("ListByClientID error: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestTaskRepository_ListActiveByWorkerID(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	workerID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"task_id", "client_id", "address_id", "worker_id", "status",
	}).AddRow(uuid.New(), uuid.New(), uuid.New(), workerID, task.StatusInProgress)

	mock.ExpectQuery("FROM tasks").
		WithArgs(workerID, task.StatusInProgress).
		WillReturnRows(rows)

	_, err := repo.ListActiveByWorkerID(ctx, workerID)
	if err != nil {
		t.Fatalf("ListActiveByWorkerID error: %v", err)
	}
}

func TestTaskRepository_ListDoneByWorkerID(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	workerID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"task_id", "client_id", "address_id", "worker_id", "status",
	}).AddRow(uuid.New(), uuid.New(), uuid.New(), workerID, task.StatusDone)

	mock.ExpectQuery("FROM tasks").
		WithArgs(workerID, task.StatusDone).
		WillReturnRows(rows)

	_, err := repo.ListDoneByWorkerID(ctx, workerID)
	if err != nil {
		t.Fatalf("ListDoneByWorkerID error: %v", err)
	}
}

func TestTaskRepository_ListOpenTasks(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{
		"task_id", "client_id", "address_id", "worker_id", "status",
	}).AddRow(uuid.New(), uuid.New(), uuid.New(), nil, task.StatusOpen)

	mock.ExpectQuery("FROM tasks").
		WithArgs(task.StatusOpen).
		WillReturnRows(rows)

	_, err := repo.ListOpenTasks(ctx)
	if err != nil {
		t.Fatalf("ListOpenTasks error: %v", err)
	}
}

func TestTaskRepository_HasOpenTaskForClient_True(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	clientID := uuid.New()

	rows := sqlmock.NewRows([]string{"dummy"}).AddRow(1)

	mock.ExpectQuery("FROM tasks").
		WithArgs(clientID, task.StatusOpen, task.StatusInProgress).
		WillReturnRows(rows)

	ok, err := repo.HasOpenTaskForClient(ctx, clientID)
	if err != nil || !ok {
		t.Fatalf("expected true, got %v err=%v", ok, err)
	}
}

func TestTaskRepository_HasOpenTaskForClient_False(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	clientID := uuid.New()

	mock.ExpectQuery("FROM tasks").
		WithArgs(clientID, task.StatusOpen, task.StatusInProgress).
		WillReturnError(sql.ErrNoRows)

	ok, err := repo.HasOpenTaskForClient(ctx, clientID)
	if err != nil || ok {
		t.Fatalf("expected false, got %v err=%v", ok, err)
	}
}

func TestTaskRepository_UpdateStatus(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	complete := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

	clientID := uuid.New()
	addressID := uuid.New()
	workerID := uuid.New()

	tsk, err := task.NewTask(uuid.New(), clientID, time.Now(), "user")
	if err != nil {
		t.Fatalf("NewTask error: %v", err)
	}


	err = tsk.StartTask("worker") 
	if err != nil {
		t.Fatalf("StartTask error: %v", err)
	}

	err = tsk.CompleteTask(complete, "worker")
	if err != nil {
		t.Fatalf("CompleteTask error: %v", err)
	}

	rows := sqlmock.NewRows([]string{
		"task_id", "client_id", "address_id", "worker_id", "status",
	}).AddRow(tsk.ID, clientID, addressID, workerID, task.StatusDone)
	mock.ExpectQuery("UPDATE tasks").
		WithArgs(tsk.ID, task.StatusDone, complete).
		WillReturnRows(rows)
	_, err = repo.UpdateStatus(ctx, tsk)
	if err != nil {
		t.Fatalf("UpdateStatus error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("mock expectations were not met: %v", err)
	}
}


func TestTaskRepository_AssignWorker(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	taskID := uuid.New()
	workerID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"task_id", "client_id", "address_id", "worker_id", "status",
	}).AddRow(taskID, uuid.New(), uuid.New(), workerID, task.StatusInProgress)

	mock.ExpectQuery("UPDATE tasks").
		WithArgs(taskID, workerID, task.StatusInProgress).
		WillReturnRows(rows)

	_, err := repo.AssignWorker(ctx, taskID, workerID)
	if err != nil {
		t.Fatalf("AssignWorker error: %v", err)
	}
}

func TestTaskRepository_DeleteTask_Success(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectExec("DELETE FROM tasks").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.DeleteTask(ctx, id); err != nil {
		t.Fatalf("DeleteTask error: %v", err)
	}
}

func TestTaskRepository_DeleteTask_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestTaskRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectExec("DELETE FROM tasks").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteTask(ctx, id)
	if !errors.Is(err, postgreserrors.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

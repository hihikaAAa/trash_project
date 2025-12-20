package workerrepo

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/worker"
	postgreserrors "github.com/hihikaAAa/TrashProject/internal/postgres/postgres_errors"
)

func newTestWorkerRepo(t *testing.T) (*WorkerRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	repo := NewWorkerRepository(db)

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func TestWorkerRepository_AddWorker_Success(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	w := &worker.Worker{
		ID:        uuid.New(),
		FirstName: "Ivan",
		Surname:   "Ivanov",
		LastName:  "Ivanovich",
		IsActive:  false,
	}

	mock.ExpectQuery("SELECT 1 FROM workers").
		WithArgs(w.FirstName, w.Surname, w.LastName).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO workers").
		WithArgs(w.ID, w.FirstName, w.Surname, w.LastName, w.IsActive).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.AddWorker(ctx, w)
	if err != nil {
		t.Fatalf("AddWorker returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWorkerRepository_AddWorker_WorkerAlreadyExists(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	w := &worker.Worker{
		ID:        uuid.New(),
		FirstName: "Ivan",
		Surname:   "Ivanov",
		LastName:  "Ivanovich",
		IsActive:  false,
	}

	rows := sqlmock.NewRows([]string{"dummy"}).AddRow(1)
	mock.ExpectQuery("SELECT 1 FROM workers").
		WithArgs(w.FirstName, w.Surname, w.LastName).
		WillReturnRows(rows)

	err := repo.AddWorker(ctx, w)
	if !errors.Is(err, postgreserrors.ErrWorkerExists) {
		t.Fatalf("expected ErrWorkerExists, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWorkerRepository_CheckNotExists_NoRows(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	mock.ExpectQuery("SELECT 1 FROM workers").
		WithArgs("Ivan", "Ivanov", "Ivanovich").
		WillReturnError(sql.ErrNoRows)

	err := repo.CheckNotExists(ctx, "Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWorkerRepository_CheckNotExists_Exists(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"dummy"}).AddRow(1)
	mock.ExpectQuery("SELECT 1 FROM workers").
		WithArgs("Ivan", "Ivanov", "Ivanovich").
		WillReturnRows(rows)

	err := repo.CheckNotExists(ctx, "Ivan", "Ivanov", "Ivanovich")
	if !errors.Is(err, postgreserrors.ErrWorkerExists) {
		t.Fatalf("expected ErrWorkerExists, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWorkerRepository_CheckNotExists_DBError(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	dbErr := errors.New("db error")

	mock.ExpectQuery("SELECT 1 FROM workers").
		WithArgs("Ivan", "Ivanov", "Ivanovich").
		WillReturnError(dbErr)

	err := repo.CheckNotExists(ctx, "Ivan", "Ivanov", "Ivanovich")
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestWorkerRepository_SetIsActive_Success(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	rows := sqlmock.NewRows(
		[]string{"worker_id", "first_name", "surname", "last_name", "is_active"},
	).AddRow(id, "Ivan", "Ivanov", "Ivanovich", true)

	mock.ExpectQuery("UPDATE workers").
		WithArgs(id, true).
		WillReturnRows(rows)

	w, err := repo.SetIsActive(ctx, id, true)
	if err != nil {
		t.Fatalf("SetIsActive returned error: %v", err)
	}

	if w.ID != id || !w.IsActive {
		t.Fatalf("unexpected worker: %+v", w)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWorkerRepository_SetIsActive_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery("UPDATE workers").
		WithArgs(id, true).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.SetIsActive(ctx, id, true)
	if !errors.Is(err, postgreserrors.ErrWorkerNotFound) {
		t.Fatalf("expected ErrWorkerNotFound, got: %v", err)
	}
}

func TestWorkerRepository_SetIsActive_DBError(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()
	dbErr := errors.New("db error")

	mock.ExpectQuery("UPDATE workers").
		WithArgs(id, true).
		WillReturnError(dbErr)

	_, err := repo.SetIsActive(ctx, id, true)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestWorkerRepository_FindActive_Success(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	id1 := uuid.New()
	id2 := uuid.New()

	rows := sqlmock.NewRows([]string{"worker_id", "first_name", "surname", "is_active"}).
		AddRow(id1, "Ivan", "Ivanov", true).
		AddRow(id2, "Petr", "Petrov", true)

	mock.ExpectQuery("FROM workers").
		WithArgs(true).
		WillReturnRows(rows)

	workers, err := repo.FindActive(ctx)
	if err != nil {
		t.Fatalf("FindActive returned error: %v", err)
	}

	if len(workers) != 2 {
		t.Fatalf("expected 2 workers, got: %d", len(workers))
	}
}

func TestWorkerRepository_FindActive_QueryError(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	dbErr := errors.New("db error")

	mock.ExpectQuery("FROM workers").
		WithArgs(true).
		WillReturnError(dbErr)

	_, err := repo.FindActive(ctx)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestWorkerRepository_GetByID_Success(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	rows := sqlmock.NewRows(
		[]string{"worker_id", "first_name", "surname", "last_name", "is_active"},
	).AddRow(id, "Ivan", "Ivanov", "Ivanovich", true)

	mock.ExpectQuery("FROM workers").
		WithArgs(id).
		WillReturnRows(rows)

	w, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	if w.ID != id || w.FirstName != "Ivan" {
		t.Fatalf("unexpected worker: %+v", w)
	}
}

func TestWorkerRepository_GetByID_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery("FROM workers").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(ctx, id)
	if !errors.Is(err, postgreserrors.ErrWorkerNotFound) {
		t.Fatalf("expected ErrWorkerNotFound, got: %v", err)
	}
}

func TestWorkerRepository_GetByID_DBError(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()
	dbErr := errors.New("db error")

	mock.ExpectQuery("FROM workers").
		WithArgs(id).
		WillReturnError(dbErr)

	_, err := repo.GetByID(ctx, id)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestWorkerRepository_List_Success(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	id1 := uuid.New()
	id2 := uuid.New()

	rows := sqlmock.NewRows([]string{"worker_id", "first_name", "surname", "last_name", "is_active"}).
		AddRow(id1, "Ivan", "Ivanov", "Ivanovich", true).
		AddRow(id2, "Petr", "Petrov", "Petrovich", false)

	mock.ExpectQuery("SELECT worker_id").
		WillReturnRows(rows)

	workers, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(workers) != 2 {
		t.Fatalf("expected 2 workers, got: %d", len(workers))
	}
}

func TestWorkerRepository_List_QueryError(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	dbErr := errors.New("db error")

	mock.ExpectQuery("SELECT worker_id").
		WillReturnError(dbErr)

	_, err := repo.List(ctx)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestWorkerRepository_DeleteWorker_Success(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectExec("DELETE FROM workers").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteWorker(ctx, id)
	if err != nil {
		t.Fatalf("DeleteWorker returned error: %v", err)
	}
}

func TestWorkerRepository_DeleteWorker_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectExec("DELETE FROM workers").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteWorker(ctx, id)
	if !errors.Is(err, postgreserrors.ErrWorkerNotFound) {
		t.Fatalf("expected ErrWorkerNotFound, got: %v", err)
	}
}

func TestWorkerRepository_DeleteWorker_ExecError(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()
	dbErr := errors.New("db error")

	mock.ExpectExec("DELETE FROM workers").
		WithArgs(id).
		WillReturnError(dbErr)

	err := repo.DeleteWorker(ctx, id)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestWorkerRepository_UpdateWorker_Success(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	w := &worker.Worker{
		ID:        uuid.New(),
		FirstName: "Ivan",
		Surname:   "Ivanov",
		LastName:  "Ivanovich",
		IsActive:  true,
	}

	rows := sqlmock.NewRows([]string{"worker_id", "first_name", "surname", "last_name", "is_active"}).
		AddRow(w.ID, w.FirstName, w.Surname, w.LastName, w.IsActive)

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE workers
		SET first_name = $2, surname = $3, last_name = $4, is_active = $5, updated_at = now()
		WHERE worker_id = $1
		RETURNING worker_id, first_name, surname, last_name, is_active
	`)).
		WithArgs(w.ID, w.FirstName, w.Surname, w.LastName, w.IsActive).
		WillReturnRows(rows)

	updated, err := repo.UpdateWorker(ctx, w)
	if err != nil {
		t.Fatalf("UpdateWorker returned error: %v", err)
	}

	if updated.ID != w.ID || updated.FirstName != w.FirstName {
		t.Fatalf("unexpected updated worker: %+v", updated)
	}
}

func TestWorkerRepository_UpdateWorker_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	w := &worker.Worker{
		ID:        uuid.New(),
		FirstName: "Ivan",
		Surname:   "Ivanov",
		LastName:  "Ivanovich",
		IsActive:  true,
	}

	mock.ExpectQuery("UPDATE workers").
		WithArgs(w.ID, w.FirstName, w.Surname, w.LastName, w.IsActive).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.UpdateWorker(ctx, w)
	if !errors.Is(err, postgreserrors.ErrWorkerNotFound) {
		t.Fatalf("expected ErrWorkerNotFound, got: %v", err)
	}
}

func TestWorkerRepository_UpdateWorker_DBError(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	w := &worker.Worker{
		ID:        uuid.New(),
		FirstName: "Ivan",
		Surname:   "Ivanov",
		LastName:  "Ivanovich",
		IsActive:  true,
	}

	dbErr := errors.New("db error")

	mock.ExpectQuery("UPDATE workers").
		WithArgs(w.ID, w.FirstName, w.Surname, w.LastName, w.IsActive).
		WillReturnError(dbErr)

	_, err := repo.UpdateWorker(ctx, w)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

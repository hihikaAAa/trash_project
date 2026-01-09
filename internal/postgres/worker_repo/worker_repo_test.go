package workerrepo

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/person"
	"github.com/hihikaAAa/TrashProject/internal/domain/worker"
	postgreserrors "github.com/hihikaAAa/TrashProject/internal/postgres/postgres_errors"
)

func newTestWorkerRepo(t *testing.T) (WorkerRepository, sqlmock.Sqlmock, func()) {
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
		AccountID: uuid.New(),
		Person:    &person.Person{FirstName: "Иван", Surname: "Иванов", LastName: "Иванович"},
		IsActive:  false,
	}

	mock.ExpectExec(
		regexp.QuoteMeta(`
	INSERT INTO workers(worker_id, account_id, first_name, surname, last_name, is_active)
	VALUES ($1, $2, $3, $4, $5, $6) 
	`),
	).
		WithArgs(w.ID, w.AccountID, w.Person.FirstName, w.Person.Surname, w.Person.LastName, false).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.AddWorker(ctx, w)
	if err != nil {
		t.Fatalf("AddWorker returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWorkerRepository_SetIsActive_Success(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()
	accID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"worker_id", "account_id", "first_name", "surname", "last_name", "is_active",
	}).AddRow(id, accID, "Иван", "Иванов", "Иванович", true)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	UPDATE workers 
	SET is_active = $2, updated_at = now()
	WHERE worker_id = $1
	RETURNING worker_id, account_id, first_name, surname, last_name, is_active
	`),
	).
		WithArgs(id, true).
		WillReturnRows(rows)

	w, err := repo.SetIsActive(ctx, id, true)
	if err != nil {
		t.Fatalf("SetIsActive returned error: %v", err)
	}

	if w.ID != id || w.AccountID != accID || !w.IsActive {
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

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	UPDATE workers 
	SET is_active = $2, updated_at = now()
	WHERE worker_id = $1
	RETURNING worker_id, account_id, first_name, surname, last_name, is_active
	`),
	).
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

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	UPDATE workers 
	SET is_active = $2, updated_at = now()
	WHERE worker_id = $1
	RETURNING worker_id, account_id, first_name, surname, last_name, is_active
	`),
	).
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
	acc1 := uuid.New()
	acc2 := uuid.New()

	rows := sqlmock.NewRows([]string{"worker_id", "account_id", "first_name", "surname", "is_active"}).
		AddRow(id1, acc1, "Иван", "Иванов", true).
		AddRow(id2, acc2, "Пётр", "Петров", true)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT worker_id, account_id, first_name, surname, is_active
	FROM workers 
	WHERE is_active = $1
	`),
	).
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

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT worker_id, account_id, first_name, surname, is_active
	FROM workers 
	WHERE is_active = $1
	`),
	).
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
	accID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"worker_id", "account_id", "first_name", "surname", "last_name", "is_active",
	}).AddRow(id, accID, "Иван", "Иванов", "Иванович", true)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT worker_id, account_id, first_name, surname, last_name, is_active
	FROM workers
	WHERE worker_id = $1
	`),
	).
		WithArgs(id).
		WillReturnRows(rows)

	w, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	if w.ID != id || w.AccountID != accID || w.Person.FirstName != "Иван" {
		t.Fatalf("unexpected worker: %+v", w)
	}
}

func TestWorkerRepository_GetByID_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT worker_id, account_id, first_name, surname, last_name, is_active
	FROM workers
	WHERE worker_id = $1
	`),
	).
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

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT worker_id, account_id, first_name, surname, last_name, is_active
	FROM workers
	WHERE worker_id = $1
	`),
	).
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
	acc1 := uuid.New()
	acc2 := uuid.New()

	rows := sqlmock.NewRows([]string{
		"worker_id", "account_id", "first_name", "surname", "last_name", "is_active",
	}).
		AddRow(id1, acc1, "Иван", "Иванов", "Иванович", true).
		AddRow(id2, acc2, "Пётр", "Петров", "Петрович", false)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT worker_id, account_id, first_name, surname, last_name, is_active
	FROM workers
	`),
	).
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

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT worker_id, account_id, first_name, surname, last_name, is_active
	FROM workers
	`),
	).
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

	mock.ExpectExec(
		regexp.QuoteMeta(`
		DELETE FROM workers
		WHERE worker_id = $1
	`),
	).
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

	mock.ExpectExec(
		regexp.QuoteMeta(`
		DELETE FROM workers
		WHERE worker_id = $1
	`),
	).
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

	mock.ExpectExec(
		regexp.QuoteMeta(`
		DELETE FROM workers
		WHERE worker_id = $1
	`),
	).
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
		AccountID: uuid.New(),
		Person:    &person.Person{FirstName: "Иван", Surname: "Иванов", LastName: "Иванович"},
		IsActive:  true,
	}

	rows := sqlmock.NewRows([]string{
		"worker_id", "account_id", "first_name", "surname", "last_name", "is_active",
	}).AddRow(w.ID, w.AccountID, w.Person.FirstName, w.Person.Surname, w.Person.LastName, w.IsActive)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	UPDATE workers
	SET account_id = $2, first_name = $3, surname = $4, last_name = $5, is_active = $6, updated_at = now()
	WHERE worker_id = $1
	RETURNING worker_id, account_id, first_name, surname, last_name, is_active
	`),
	).
		WithArgs(w.ID, w.AccountID, w.Person.FirstName, w.Person.Surname, w.Person.LastName, w.IsActive).
		WillReturnRows(rows)

	updated, err := repo.UpdateWorker(ctx, w)
	if err != nil {
		t.Fatalf("UpdateWorker returned error: %v", err)
	}

	if updated.ID != w.ID || updated.AccountID != w.AccountID || updated.Person.FirstName != w.Person.FirstName {
		t.Fatalf("unexpected updated worker: %+v", updated)
	}
}

func TestWorkerRepository_UpdateWorker_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestWorkerRepo(t)
	defer cleanup()

	ctx := context.Background()

	w := &worker.Worker{
		ID:        uuid.New(),
		AccountID: uuid.New(),
		Person:    &person.Person{FirstName: "Иван", Surname: "Иванов", LastName: "Иванович"},
		IsActive:  false,
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	UPDATE workers
	SET account_id = $2, first_name = $3, surname = $4, last_name = $5, is_active = $6, updated_at = now()
	WHERE worker_id = $1
	RETURNING worker_id, account_id, first_name, surname, last_name, is_active
	`),
	).
		WithArgs(w.ID, w.AccountID, w.Person.FirstName, w.Person.Surname, w.Person.LastName, w.IsActive).
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
		AccountID: uuid.New(),
		Person:    &person.Person{FirstName: "Иван", Surname: "Иванов", LastName: "Иванович"},
		IsActive:  true,
	}

	dbErr := errors.New("db error")

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	UPDATE workers
	SET account_id = $2, first_name = $3, surname = $4, last_name = $5, is_active = $6, updated_at = now()
	WHERE worker_id = $1
	RETURNING worker_id, account_id, first_name, surname, last_name, is_active
	`),
	).
		WithArgs(w.ID, w.AccountID, w.Person.FirstName, w.Person.Surname, w.Person.LastName, w.IsActive).
		WillReturnError(dbErr)

	_, err := repo.UpdateWorker(ctx, w)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

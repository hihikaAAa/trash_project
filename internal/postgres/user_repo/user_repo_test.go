package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/user"
	postgreserrors "github.com/hihikaAAa/TrashProject/internal/postgres/postgres_errors"
)

func newTestUserRepo(t *testing.T) (*UserRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	repo := NewUserRepository(db)

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func TestUserRepository_AddUser_Success(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	u := &user.User{
		ID:        uuid.New(),
		FirstName: "Ivan",
		Surname:   "Ivanov",
		LastName:  "Ivanovich",
		AddressID: uuid.New(),
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT 1 
			FROM users 
			WHERE first_name = $1 AND surname = $2 AND last_name = $3
		`),
	).
		WithArgs(u.FirstName, u.Surname, u.LastName).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec(
		regexp.QuoteMeta(`
			INSERT INTO users(user_id, first_name, surname, last_name, address_id)
			VALUES ($1, $2, $3, $4, $5)
		`),
	).
		WithArgs(u.ID, u.FirstName, u.Surname, u.LastName, u.AddressID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.AddUser(ctx, u)
	if err != nil {
		t.Fatalf("AddUser returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUserRepository_AddUser_UserAlreadyExists(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	u := &user.User{
		ID:        uuid.New(),
		FirstName: "Ivan",
		Surname:   "Ivanov",
		LastName:  "Ivanovich",
		AddressID: uuid.New(),
	}

	rows := sqlmock.NewRows([]string{"dummy"}).AddRow(1)
	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT 1 
			FROM users 
			WHERE first_name = $1 AND surname = $2 AND last_name = $3
		`),
	).
		WithArgs(u.FirstName, u.Surname, u.LastName).
		WillReturnRows(rows)

	err := repo.AddUser(ctx, u)
	if !errors.Is(err, postgreserrors.ErrUserExists) {
		t.Fatalf("expected ErrUserExists, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUserRepository_CheckNotExists_NoRows(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT 1 
			FROM users 
			WHERE first_name = $1 AND surname = $2 AND last_name = $3
		`),
	).
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

func TestUserRepository_CheckNotExists_Exists(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"dummy"}).AddRow(1)
	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT 1 
			FROM users 
			WHERE first_name = $1 AND surname = $2 AND last_name = $3
		`),
	).
		WithArgs("Ivan", "Ivanov", "Ivanovich").
		WillReturnRows(rows)

	err := repo.CheckNotExists(ctx, "Ivan", "Ivanov", "Ivanovich")
	if !errors.Is(err, postgreserrors.ErrUserExists) {
		t.Fatalf("expected ErrUserExists, got: %v", err)
	}
}

func TestUserRepository_CheckNotExists_DBError(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	dbErr := errors.New("db error")
	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT 1 
			FROM users 
			WHERE first_name = $1 AND surname = $2 AND last_name = $3
		`),
	).
		WithArgs("Ivan", "Ivanov", "Ivanovich").
		WillReturnError(dbErr)

	err := repo.CheckNotExists(ctx, "Ivan", "Ivanov", "Ivanovich")
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestUserRepository_GetByID_Success(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()
	addrID := uuid.New()

	rows := sqlmock.NewRows([]string{"user_id", "first_name", "surname", "last_name", "address_id"}).
		AddRow(id, "Ivan", "Ivanov", "Ivanovich", addrID)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT user_id, first_name, surname, last_name, address_id
			FROM users
			WHERE user_id = $1
		`),
	).
		WithArgs(id).
		WillReturnRows(rows)

	u, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	if u.ID != id || u.FirstName != "Ivan" || u.Surname != "Ivanov" || u.LastName != "Ivanovich" || u.AddressID != addrID {
		t.Fatalf("unexpected user: %+v", u)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT user_id, first_name, surname, last_name, address_id
			FROM users
			WHERE user_id = $1
		`),
	).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(ctx, id)
	if !errors.Is(err, postgreserrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserRepository_FindByFullName_Success(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()
	addrID := uuid.New()

	rows := sqlmock.NewRows([]string{"user_id", "first_name", "surname", "last_name", "address_id"}).
		AddRow(id, "Ivan", "Ivanov", "Ivanovich", addrID)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT user_id, first_name, surname, last_name, address_id
			FROM users
			WHERE first_name = $1 AND surname = $2 AND last_name = $3
		`),
	).
		WithArgs("Ivan", "Ivanov", "Ivanovich").
		WillReturnRows(rows)

	u, err := repo.FindByFullName(ctx, "Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("FindByFullName returned error: %v", err)
	}

	if u.ID != id || u.AddressID != addrID {
		t.Fatalf("unexpected user: %+v", u)
	}
}

func TestUserRepository_FindByFullName_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT user_id, first_name, surname, last_name, address_id
			FROM users
			WHERE first_name = $1 AND surname = $2 AND last_name = $3
		`),
	).
		WithArgs("Ivan", "Ivanov", "Ivanovich").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.FindByFullName(ctx, "Ivan", "Ivanov", "Ivanovich")
	if !errors.Is(err, postgreserrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserRepository_List_Success(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	id1 := uuid.New()
	id2 := uuid.New()
	addr1 := uuid.New()
	addr2 := uuid.New()

	rows := sqlmock.NewRows([]string{"user_id", "first_name", "surname", "last_name", "address_id"}).
		AddRow(id1, "Ivan", "Ivanov", "Ivanovich", addr1).
		AddRow(id2, "Petr", "Petrov", "Petrovich", addr2)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT user_id, first_name, surname, last_name, address_id
			FROM users
		`),
	).
		WillReturnRows(rows)

	users, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got: %d", len(users))
	}
}

func TestUserRepository_List_QueryError(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	dbErr := errors.New("db error")
	mock.ExpectQuery(
		regexp.QuoteMeta(`
			SELECT user_id, first_name, surname, last_name, address_id
			FROM users
		`),
	).
		WillReturnError(dbErr)

	_, err := repo.List(ctx)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestUserRepository_DeleteUser_Success(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectExec(
		regexp.QuoteMeta(`
			DELETE FROM users
			WHERE user_id = $1
		`),
	).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(ctx, id)
	if err != nil {
		t.Fatalf("DeleteUser returned error: %v", err)
	}
}

func TestUserRepository_DeleteUser_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectExec(
		regexp.QuoteMeta(`
			DELETE FROM users
			WHERE user_id = $1
		`),
	).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteUser(ctx, id)
	if !errors.Is(err, postgreserrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestUserRepository_UpdateUser_Success(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	u := &user.User{
		ID:        uuid.New(),
		FirstName: "Ivan",
		Surname:   "Ivanov",
		LastName:  "Ivanovich",
		AddressID: uuid.New(),
	}

	rows := sqlmock.NewRows([]string{"user_id", "first_name", "surname", "last_name", "address_id"}).
		AddRow(u.ID, u.FirstName, u.Surname, u.LastName, u.AddressID)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
			UPDATE users
			SET first_name = $2,
				surname    = $3,
				last_name  = $4,
				address_id = $5,
				updated_at = now()
			WHERE user_id = $1
			RETURNING user_id, first_name, surname, last_name, address_id
		`),
	).
		WithArgs(u.ID, u.FirstName, u.Surname, u.LastName, u.AddressID).
		WillReturnRows(rows)

	updated, err := repo.UpdateUser(ctx, u)
	if err != nil {
		t.Fatalf("UpdateUser returned error: %v", err)
	}

	if updated.ID != u.ID || updated.FirstName != u.FirstName {
		t.Fatalf("unexpected updated user: %+v", updated)
	}
}

func TestUserRepository_UpdateUser_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	u := &user.User{
		ID:        uuid.New(),
		FirstName: "Ivan",
		Surname:   "Ivanov",
		LastName:  "Ivanovich",
		AddressID: uuid.New(),
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`
			UPDATE users
			SET first_name = $2,
				surname    = $3,
				last_name  = $4,
				address_id = $5,
				updated_at = now()
			WHERE user_id = $1
			RETURNING user_id, first_name, surname, last_name, address_id
		`),
	).
		WithArgs(u.ID, u.FirstName, u.Surname, u.LastName, u.AddressID).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.UpdateUser(ctx, u)
	if !errors.Is(err, postgreserrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

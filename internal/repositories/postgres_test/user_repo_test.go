package postgres

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
	"github.com/hihikaAAa/trash_project/internal/repositories/postgres"
)

func newTestUserRepo(t *testing.T) (postgres.userrepo, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	repo := postgres.userrepo.NewUserRepository(db)

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func TestUserRepository_AddUser_Success(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	u, _ := user.NewUser("Иван", "Иванов", "Иванович", uuid.New(), uuid.New())

	mock.ExpectExec(
		regexp.QuoteMeta(`
	INSERT INTO users(user_id, account_id, first_name, surname, last_name, address_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	`),
	).
		WithArgs(u.ID, u.AccountID, u.Person.FirstName, u.Person.Surname, u.Person.LastName, u.AddressID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.AddUser(ctx, u)
	if err != nil {
		t.Fatalf("AddUser returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

/*
func TestUserRepository_AddUser_UserAlreadyExists(t *testing.T) { //TODO: CHECK TESTS EMAIL

		repo, mock, cleanup := newTestUserRepo(t)
		defer cleanup()

		ctx := context.Background()
		p, _ := person.NewPerson("Ivan", "Ivanov", "Ivanovich")
		u, _ := user.NewUser(p, uuid.New())

		rows := sqlmock.NewRows([]string{"dummy"}).AddRow(1)
		mock.ExpectQuery(
			regexp.QuoteMeta(`
				SELECT 1
				FROM users
				WHERE first_name = $1 AND surname = $2 AND last_name = $3
			`),
		).
			WithArgs(u.Person.FirstName, u.Person.Surname, u.Person.LastName).
			WillReturnRows(rows)

		err := repo.AddUser(ctx, u)
		if !errors.Is(err, postgreserrors.ErrUserExists) {
			t.Fatalf("expected ErrUserExists, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	}

func TestUserRepository_CheckNotExists_NoRows(t *testing.T) { //TODO: CHECK TESTS EMAIL

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

		err := repo.CheckNotExists(ctx, "email")
		if err != nil {
			t.Fatalf("expected nil, got: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	}

func TestUserRepository_CheckNotExists_Exists(t *testing.T) { //TODO: CHECK TESTS EMAIL

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

		err := repo.CheckNotExists(ctx, "email")
		if !errors.Is(err, postgreserrors.ErrUserExists) {
			t.Fatalf("expected ErrUserExists, got: %v", err)
		}
	}

func TestUserRepository_CheckNotExists_DBError(t *testing.T) { //TODO: CHECK TESTS EMAIL

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

		err := repo.CheckNotExists(ctx, "email")
		if err == nil || !errors.Is(err, dbErr) {
			t.Fatalf("expected wrapped db error, got: %v", err)
		}
	}
*/
func TestUserRepository_GetByID_Success(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()
	accID := uuid.New()
	addrID := uuid.New()

	rows := sqlmock.NewRows([]string{"user_id", "account_id", "first_name", "surname", "last_name", "address_id"}).
		AddRow(id, accID, "Иван", "Иванов", "Иванович", addrID)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT user_id, account_id, first_name, surname, last_name, address_id
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

	if u.ID != id ||
		u.AccountID != accID ||
		u.Person.FirstName != "Иван" ||
		u.Person.Surname != "Иванов" ||
		u.Person.LastName != "Иванович" ||
		u.AddressID != addrID {
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
	SELECT user_id, account_id, first_name, surname, last_name, address_id
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
	accID := uuid.New()
	addrID := uuid.New()

	rows := sqlmock.NewRows([]string{"user_id", "account_id", "first_name", "surname", "last_name", "address_id"}).
		AddRow(id, accID, "Иван", "Иванов", "Иванович", addrID)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT user_id, account_id, first_name, surname, last_name, address_id
	FROM users
	WHERE first_name = $1 AND surname = $2 AND last_name = $3
	`),
	).
		WithArgs("Иван", "Иванов", "Иванович").
		WillReturnRows(rows)

	u, err := repo.FindByFullName(ctx, "Иван", "Иванов", "Иванович")
	if err != nil {
		t.Fatalf("FindByFullName returned error: %v", err)
	}

	if u.ID != id || u.AccountID != accID || u.AddressID != addrID {
		t.Fatalf("unexpected user: %+v", u)
	}
}

func TestUserRepository_FindByFullName_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT user_id, account_id, first_name, surname, last_name, address_id
	FROM users
	WHERE first_name = $1 AND surname = $2 AND last_name = $3
	`),
	).
		WithArgs("Иван", "Иванов", "Иванович").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.FindByFullName(ctx, "Иван", "Иванов", "Иванович")
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
	acc1 := uuid.New()
	acc2 := uuid.New()
	addr1 := uuid.New()
	addr2 := uuid.New()

	rows := sqlmock.NewRows([]string{"user_id", "account_id", "first_name", "surname", "last_name", "address_id"}).
		AddRow(id1, acc1, "Иван", "Иванов", "Иванович", addr1).
		AddRow(id2, acc2, "Петр", "Петров", "Петрович", addr2)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	SELECT user_id, account_id, first_name, surname, last_name, address_id
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
	SELECT user_id, account_id, first_name, surname, last_name, address_id
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
	u, _ := user.NewUser("Иван", "Иванов", "Иванович", uuid.New(), uuid.New())

	rows := sqlmock.NewRows([]string{"user_id", "account_id", "first_name", "surname", "last_name", "address_id"}).
		AddRow(u.ID, u.AccountID, u.Person.FirstName, u.Person.Surname, u.Person.LastName, u.AddressID)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	UPDATE users
	SET account_id = $2, first_name = $3, surname = $4, last_name = $5, address_id = $6, updated_at = now()
	WHERE user_id = $1
	RETURNING user_id, account_id, first_name, surname, last_name, address_id
	`),
	).
		WithArgs(u.ID, u.AccountID, u.Person.FirstName, u.Person.Surname, u.Person.LastName, u.AddressID).
		WillReturnRows(rows)

	updated, err := repo.UpdateUser(ctx, u)
	if err != nil {
		t.Fatalf("UpdateUser returned error: %v", err)
	}

	if updated.ID != u.ID || updated.AccountID != u.AccountID || updated.Person.FirstName != u.Person.FirstName {
		t.Fatalf("unexpected updated user: %+v", updated)
	}
}

func TestUserRepository_UpdateUser_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	u, _ := user.NewUser("Иван", "Иванов", "Иванович", uuid.New(), uuid.New())

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	UPDATE users
	SET account_id = $2, first_name = $3, surname = $4, last_name = $5, address_id = $6, updated_at = now()
	WHERE user_id = $1
	RETURNING user_id, account_id, first_name, surname, last_name, address_id
	`),
	).
		WithArgs(u.ID, u.AccountID, u.Person.FirstName, u.Person.Surname, u.Person.LastName, u.AddressID).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.UpdateUser(ctx, u)
	if !errors.Is(err, postgreserrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

package addressrepo

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/address"
	repoerrors "github.com/hihikaAAa/TrashProject/internal/repository/postgres/repo_errors"
)

func newTestAddressRepo(t *testing.T) (*AddressRepository, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}

	repo := NewAddressRepository(db)

	cleanup := func() {
		_ = db.Close()
	}

	return repo, mock, cleanup
}

func TestAddressRepository_AddAddress_Success(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()

	a := &address.Address{
		ID:              uuid.New(),
		Street:          "Main",
		HouseNumber:     "10",
		Entrance:        "1",
		FloorNumber:     5,
		ApartmentNumber: 12,
	}

	mock.ExpectQuery("SELECT 1 FROM addresses").
		WithArgs(a.Street, a.HouseNumber, a.Entrance, a.FloorNumber, a.ApartmentNumber).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO addresses").
		WithArgs(a.ID, a.Street, a.HouseNumber, a.Entrance, a.FloorNumber, a.ApartmentNumber).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.AddAddress(ctx, a)
	if err != nil {
		t.Fatalf("AddAddress returned error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAddressRepository_AddAddress_AlreadyExists(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()

	a := &address.Address{
		ID:              uuid.New(),
		Street:          "Main",
		HouseNumber:     "10",
		Entrance:        "1",
		FloorNumber:     5,
		ApartmentNumber: 12,
	}

	rows := sqlmock.NewRows([]string{"dummy"}).AddRow(1)
	mock.ExpectQuery("SELECT 1 FROM addresses").
		WithArgs(a.Street, a.HouseNumber, a.Entrance, a.FloorNumber, a.ApartmentNumber).
		WillReturnRows(rows)

	err := repo.AddAddress(ctx, a)
	if !errors.Is(err, nil) {
		t.Fatalf("expected ErrAddressExists, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAddressRepository_CheckNotExists_NoRows(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()

	mock.ExpectQuery("SELECT 1 FROM addresses").
		WithArgs("Main", "10", "1", 5, 12).
		WillReturnError(sql.ErrNoRows)

	err := repo.CheckNotExists(ctx, "Main", "10", "1", 5, 12)
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAddressRepository_CheckNotExists_Exists(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"dummy"}).AddRow(1)
	mock.ExpectQuery("SELECT 1 FROM addresses").
		WithArgs("Main", "10", "1", 5, 12).
		WillReturnRows(rows)

	err := repo.CheckNotExists(ctx, "Main", "10", "1", 5, 12)
	if !errors.Is(err, repoerrors.ErrAddressExists) {
		t.Fatalf("expected ErrAddressExists, got: %v", err)
	}
}

func TestAddressRepository_CheckNotExists_DBError(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()
	dbErr := errors.New("db error")

	mock.ExpectQuery("SELECT 1 FROM addresses").
		WithArgs("Main", "10", "1", 5, 12).
		WillReturnError(dbErr)

	err := repo.CheckNotExists(ctx, "Main", "10", "1", 5, 12)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestAddressRepository_GetByID_Success(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	rows := sqlmock.NewRows(
		[]string{"address_id", "street", "house_number", "entrance", "floor_number", "apartment_number"},
	).AddRow(id, "Main", "10", "1", 5, 12)

	mock.ExpectQuery("FROM addresses").
		WithArgs(id).
		WillReturnRows(rows)

	a, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID returned error: %v", err)
	}

	if a.ID != id || a.Street != "Main" || a.HouseNumber != "10" {
		t.Fatalf("unexpected address: %+v", a)
	}
}

func TestAddressRepository_GetByID_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectQuery("FROM addresses").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(ctx, id)
	if !errors.Is(err, repoerrors.ErrAddressNotFound) {
		t.Fatalf("expected ErrAddressNotFound, got: %v", err)
	}
}

func TestAddressRepository_GetByID_DBError(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()
	dbErr := errors.New("db error")

	mock.ExpectQuery("FROM addresses").
		WithArgs(id).
		WillReturnError(dbErr)

	_, err := repo.GetByID(ctx, id)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestAddressRepository_List_Success(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()

	id1 := uuid.New()
	id2 := uuid.New()

	rows := sqlmock.NewRows([]string{"address_id", "street", "house_number", "entrance", "floor_number", "apartment_number"}).
		AddRow(id1, "Main", "10", "1", 5, 12).
		AddRow(id2, "Main", "12", "2", 6, 24)

	mock.ExpectQuery("SELECT address_id").
		WillReturnRows(rows)

	addresses, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(addresses) != 2 {
		t.Fatalf("expected 2 addresses, got: %d", len(addresses))
	}
}

func TestAddressRepository_List_QueryError(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()
	dbErr := errors.New("db error")

	mock.ExpectQuery("SELECT address_id").
		WillReturnError(dbErr)

	_, err := repo.List(ctx)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestAddressRepository_DeleteAddress_Success(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectExec("DELETE FROM addresses").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteAddress(ctx, id)
	if err != nil {
		t.Fatalf("DeleteAddress returned error: %v", err)
	}
}

func TestAddressRepository_DeleteAddress_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()

	mock.ExpectExec("DELETE FROM addresses").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteAddress(ctx, id)
	if !errors.Is(err, repoerrors.ErrAddressNotFound) {
		t.Fatalf("expected ErrAddressNotFound, got: %v", err)
	}
}

func TestAddressRepository_DeleteAddress_ExecError(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()
	id := uuid.New()
	dbErr := errors.New("db error")

	mock.ExpectExec("DELETE FROM addresses").
		WithArgs(id).
		WillReturnError(dbErr)

	err := repo.DeleteAddress(ctx, id)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

func TestAddressRepository_UpdateAddress_Success(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()

	a := &address.Address{
		ID:              uuid.New(),
		Street:          "Main",
		HouseNumber:     "10",
		Entrance:        "1",
		FloorNumber:     5,
		ApartmentNumber: 12,
	}

	rows := sqlmock.NewRows([]string{"address_id", "street", "house_number", "entrance", "floor_number", "apartment_number"}).
		AddRow(a.ID, a.Street, a.HouseNumber, a.Entrance, a.FloorNumber, a.ApartmentNumber)

	mock.ExpectQuery("UPDATE addresses").
		WithArgs(a.ID, a.Street, a.HouseNumber, a.Entrance, a.FloorNumber, a.ApartmentNumber).
		WillReturnRows(rows)

	updated, err := repo.UpdateAddress(ctx, a)
	if err != nil {
		t.Fatalf("UpdateAddress returned error: %v", err)
	}

	if updated.ID != a.ID || updated.Street != a.Street {
		t.Fatalf("unexpected updated address: %+v", updated)
	}
}

func TestAddressRepository_UpdateAddress_NotFound(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()

	a := &address.Address{
		ID:              uuid.New(),
		Street:          "Main",
		HouseNumber:     "10",
		Entrance:        "1",
		FloorNumber:     5,
		ApartmentNumber: 12,
	}

	mock.ExpectQuery("UPDATE addresses").
		WithArgs(a.ID, a.Street, a.HouseNumber, a.Entrance, a.FloorNumber, a.ApartmentNumber).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.UpdateAddress(ctx, a)
	if !errors.Is(err, repoerrors.ErrAddressNotFound) {
		t.Fatalf("expected ErrAddressNotFound, got: %v", err)
	}
}

func TestAddressRepository_UpdateAddress_DBError(t *testing.T) {
	repo, mock, cleanup := newTestAddressRepo(t)
	defer cleanup()

	ctx := context.Background()

	a := &address.Address{
		ID:              uuid.New(),
		Street:          "Main",
		HouseNumber:     "10",
		Entrance:        "1",
		FloorNumber:     5,
		ApartmentNumber: 12,
	}

	dbErr := errors.New("db error")

	mock.ExpectQuery("UPDATE addresses").
		WithArgs(a.ID, a.Street, a.HouseNumber, a.Entrance, a.FloorNumber, a.ApartmentNumber).
		WillReturnError(dbErr)

	_, err := repo.UpdateAddress(ctx, a)
	if err == nil || !errors.Is(err, dbErr) {
		t.Fatalf("expected wrapped db error, got: %v", err)
	}
}

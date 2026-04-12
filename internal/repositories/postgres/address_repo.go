// Package addressrepo
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hihikaAAa/TrashProject/internal/domain/address"
	postgreserrors "github.com/hihikaAAa/TrashProject/internal/postgres/postgres_errors"
)

type AddressRepository interface {
	AddAddress(ctx context.Context, address *address.Address) error
	CheckNotExists(ctx context.Context, street, houseNumber, entrance string, floorNumber, apartmentNumber int) error
	GetByID(ctx context.Context, id uuid.UUID) (*address.Address, error)
	List(ctx context.Context) ([]*address.Address, error)
	DeleteAddress(ctx context.Context, id uuid.UUID) error
	UpdateAddress(ctx context.Context, addr *address.Address) (*address.Address, error)
}
type addressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) AddAddress(ctx context.Context, address *address.Address) error {
	const op = "internal.postgres.address_repo.AddAddress"

	const q = `
	INSERT INTO addresses(address_id, street, house_number, entrance, floor_number, apartment_number)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	err := r.CheckNotExists(ctx, address.Street, address.HouseNumber, address.Entrance, address.FloorNumber, address.ApartmentNumber)
	if errors.Is(err, postgreserrors.ErrAddressExists) {
		return nil
	}
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, q, address.ID, address.Street, address.HouseNumber, address.Entrance, address.FloorNumber, address.ApartmentNumber)
	if err != nil {
		return fmt.Errorf("%s, ExecContext: %w", op, err)
	}
	return nil
}

func (r *addressRepository) CheckNotExists(ctx context.Context, street, houseNumber, entrance string, floorNumber, apartmentNumber int) error {
	const op = "internal.postgres.address_repo.CheckNotExists"

	const q = `
	SELECT 1 FROM addresses 
	WHERE street = $1 AND house_number = $2 AND entrance = $3 AND floor_number = $4 AND apartment_number = $5
	`

	var dummy int
	err := r.db.QueryRowContext(ctx, q, street, houseNumber, entrance, floorNumber, apartmentNumber).Scan(&dummy)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return postgreserrors.ErrAddressExists
}

func (r *addressRepository) GetByID(ctx context.Context, id uuid.UUID) (*address.Address, error) {
	const op = "internal.postgres.address_repo.GetByID"

	const q = `
	SELECT address_id, street, house_number, entrance, floor_number, apartment_number
	FROM addresses
	WHERE address_id = $1
	`
	a := &address.Address{}

	err := r.db.QueryRowContext(ctx, q, id).Scan(&a.ID, &a.Street, &a.HouseNumber, &a.Entrance, &a.FloorNumber, &a.ApartmentNumber)
	if err == sql.ErrNoRows {
		return nil, postgreserrors.ErrAddressNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}

	return a, nil
}

func (r *addressRepository) List(ctx context.Context) ([]*address.Address, error) {
	const op = "internal.postgres.address_repo.List"

	const q = `
	SELECT address_id, street, house_number, entrance, floor_number, apartment_number
	FROM addresses
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	addresses := make([]*address.Address, 0)
	for rows.Next() {
		a := &address.Address{}
		err := rows.Scan(&a.ID, &a.Street, &a.HouseNumber, &a.Entrance, &a.FloorNumber, &a.ApartmentNumber)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		addresses = append(addresses, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return addresses, nil
}

func (r *addressRepository) DeleteAddress(ctx context.Context, id uuid.UUID) error {
	const op = "internal.postgres.address_repo.DeleteAddress"

	const q = `
	DELETE FROM addresses
	WHERE address_id = $1
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
		return fmt.Errorf("%s: %w", op, postgreserrors.ErrAddressNotFound)
	}

	return nil
}

func (r *addressRepository) UpdateAddress(ctx context.Context, addr *address.Address) (*address.Address, error) {
	const op = "internal.postgres.address_repo.UpdateAddress"

	const q = `
	UPDATE addresses
	SET street = $2,  house_number = $3, entrance = $4,  floor_number = $5,  apartment_number = $6
	WHERE address_id = $1
	RETURNING address_id, street, house_number, entrance, floor_number, apartment_number
	`

	a := &address.Address{}

	err := r.db.QueryRowContext(ctx, q, addr.ID, addr.Street, addr.HouseNumber, addr.Entrance, addr.FloorNumber, addr.ApartmentNumber).Scan(
		&a.ID, &a.Street, &a.HouseNumber, &a.Entrance, &a.FloorNumber, &a.ApartmentNumber,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s, %w", op, postgreserrors.ErrAddressNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return a, nil
}

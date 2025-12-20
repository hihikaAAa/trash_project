package userrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hihikaAAa/TrashProject/internal/domain/user"
	repoerrors "github.com/hihikaAAa/TrashProject/internal/repository/postgres/repo_errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) AddUser(ctx context.Context, user *user.User) error {
	const op = "internal.repository.postgres.user_repo.AddUser"

	const q = `
	INSERT INTO users(user_id, first_name, surname, last_name, address_id)
	VALUES ($1, $2, $3, $4, $5)
	`
	err := r.CheckNotExists(ctx, user.FirstName, user.Surname, user.LastName)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, q, user.ID, user.FirstName, user.Surname, user.LastName, user.AddressID)
	if err != nil {
		return fmt.Errorf("%s, ExecContext: %w", op, err)
	}
	return nil
}

func (r *UserRepository) CheckNotExists(ctx context.Context, name, surname, last_name string) error {
	const op = "internal.repository.postgres.user_repo.CheckNotExists"

	const q = `
	SELECT 1 
	FROM users 
	WHERE first_name = $1 AND surname = $2 AND last_name = $3
	`
	var dummy int
	err := r.db.QueryRowContext(ctx, q, name, surname, last_name).Scan(&dummy)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return repoerrors.ErrUserExists
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	const op = "internal.repository.postgres.user_repo.GetByID"

	const q = `
	SELECT user_id, first_name, surname, last_name, address_id
	FROM users
	WHERE user_id = $1
	`

	u := &user.User{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.FirstName, &u.Surname, &u.LastName, &u.AddressID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: %w", op, repoerrors.ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return u, nil
}

func (r *UserRepository) FindByFullName(ctx context.Context, name, surname, lastName string) (*user.User, error) {
	const op = "internal.repository.postgres.user_repo.FindByFullName"

	const q = `
	SELECT user_id, first_name, surname, last_name, address_id
	FROM users
	WHERE first_name = $1 AND surname = $2 AND last_name = $3
	`

	u := &user.User{}
	err := r.db.QueryRowContext(ctx, q, name, surname, lastName).Scan(&u.ID, &u.FirstName, &u.Surname, &u.LastName, &u.AddressID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: %w", op, repoerrors.ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return u, nil
}

func (r *UserRepository) List(ctx context.Context) ([]*user.User, error) {
	const op = "internal.repository.postgres.user_repo.List"

	const q = `
	SELECT user_id, first_name, surname, last_name, address_id
	FROM users
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	users := make([]*user.User, 0)

	for rows.Next() {
		u := &user.User{}
		err := rows.Scan(&u.ID, &u.FirstName, &u.Surname, &u.LastName, &u.AddressID)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return users, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const op = "internal.repository.postgres.user_repo.DeleteUser"

	const q = `
		DELETE FROM users
		WHERE user_id = $1
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
		return fmt.Errorf("%s: %w", op, repoerrors.ErrUserNotFound)
	}

	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, u *user.User) (*user.User, error) {
	const op = "internal.repository.postgres.user_repo.UpdateUser"

	const q = `
	UPDATE users
	SET first_name = $2, surname = $3, last_name = $4, address_id = $5, updated_at = now()
	WHERE user_id = $1
	RETURNING user_id, first_name, surname, last_name, address_id
	`

	user := &user.User{}
	err := r.db.QueryRowContext(ctx, q, u.ID, u.FirstName, u.Surname, u.LastName, u.AddressID).Scan(
		&user.ID, &user.FirstName, &user.Surname, &user.LastName, &user.AddressID,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s, %w", op, repoerrors.ErrUserNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return user, nil
}

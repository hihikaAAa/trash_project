// Package accountrepo
package accountrepo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hihikaAAa/TrashProject/internal/domain/account"
	postgreserrors "github.com/hihikaAAa/TrashProject/internal/postgres/postgres_errors"
)

type AccountRepository interface{
	AddAccount(ctx context.Context, account *account.Account) error
	CheckIfExistsByEmail(ctx context.Context, email string) error
	GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error)
	List(ctx context.Context) ([]*account.Account, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
	UpdateAccount(ctx context.Context, acc *account.Account) (*account.Account, error)
	ListByUserRole(ctx context.Context, role string) ([]*account.Account, error)
	ListByWorkerRole(ctx context.Context, role string) ([]*account.Account, error)
}

type accountRepository struct{
	db *sql.DB
}

func NewAccountRepositody(db *sql.DB) AccountRepository{
	return &accountRepository{db: db}
}

func (a *accountRepository) AddAccount(ctx context.Context, account *account.Account) error{
	const op = "internal.postgres.account_repo.AddAccount"

	const q = `
	INSERT INTO accounts(account_id, email, password_hash, role)
	VALUES($1, $2, $3, $4)
	`

	_, err := a.db.ExecContext(ctx, q, account.ID, account.Email, account.PasswordHash, account.Role)
	if err != nil{
		return fmt.Errorf("%s, Exec Conetext: %w", op, err )
	}

	return nil
}

func (a *accountRepository) CheckIfExistsByEmail(ctx context.Context, email string) error{
	const op = "internal.postgres.account_repo.CheckIfExistsByEmail"

	const q = `
	SELECT account_id, email, password_hash, role
	FROM accounts 
	WHERE email = $1
	`

	acc := account.Account{}
	err := a.db.QueryRowContext(ctx, q, email).Scan(&acc.ID,&acc.Email,&acc.PasswordHash,&acc.Role)
	if err == sql.ErrNoRows{
		return nil
	}
	
	if err != nil{
		return fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}

	return postgreserrors.ErrAccountExists
}

func (a *accountRepository) GetByID(ctx context.Context, id uuid.UUID) (*account.Account, error) {
	const op = "internal.postgres.account_repo.GetByID"

	const q = `
	SELECT account_id, email, password_hash, role
	FROM accounts
	WHERE account_id = $1
	`

	acc := &account.Account{}
	err := a.db.QueryRowContext(ctx, q, id).Scan(&acc.ID, &acc.Email, &acc.PasswordHash, &acc.Role)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s: %w", op, postgreserrors.ErrAccountNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return acc, nil
}

func (a *accountRepository) List(ctx context.Context) ([]*account.Account, error) {
	const op = "internal.postgres.account_repo.List"

	const q = `
	SELECT account_id, email, password_hash, role
	FROM accounts
	`

	rows, err := a.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	accounts := make([]*account.Account, 0)

	for rows.Next() {
		acc := &account.Account{}
		err := rows.Scan(&acc.ID, &acc.Email, &acc.PasswordHash, &acc.Role)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return accounts, nil
}

func (a *accountRepository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	const op = "internal.postgres.account_repo.DeleteAccount"

	const q = `
		DELETE FROM accounts
		WHERE account_id = $1
	`

	res, err := a.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("%s, ExecContext: %w", op, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s, RowsAffected: %w", op, err)
	}
	if affected == 0 {
		return fmt.Errorf("%s: %w", op, postgreserrors.ErrAccountNotFound)
	}

	return nil
}

func (a *accountRepository) UpdateAccount(ctx context.Context, acc *account.Account) (*account.Account, error) {
	const op = "internal.postgres.account_repo.UpdateAccount"

	const q = `
	UPDATE accounts
	SET email = $2, password_hash = $3, role = $4, updated_at = now()
	WHERE account_id = $1
	RETURNING account_id, email, password_hash, role
	`

	updAcc := &account.Account{}
	err := a.db.QueryRowContext(ctx, q, acc.ID, acc.Email, acc.PasswordHash, acc.Role ).Scan(
		&updAcc.ID, &updAcc.Email, &updAcc.PasswordHash, &updAcc.Role,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%s, %w", op, postgreserrors.ErrAccountNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s, QueryRowContext: %w", op, err)
	}
	return updAcc, nil
}

func (a *accountRepository) ListByUserRole(ctx context.Context, role string) ([]*account.Account, error) {
	const op = "internal.postgres.account_repo.ListByUserRole"

	const q = `
	SELECT account_id, email, password_hash, role
	FROM accounts
	WHERER role = $1
	`

	rows, err := a.db.QueryContext(ctx, q, role)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	accounts := make([]*account.Account, 0)

	for rows.Next() {
		acc := &account.Account{}
		err := rows.Scan(&acc.ID, &acc.Email, &acc.PasswordHash, &acc.Role)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return accounts, nil
}

func (a *accountRepository) ListByWorkerRole(ctx context.Context, role string) ([]*account.Account, error) {
	const op = "internal.postgres.account_repo.ListByWorkerRole"

	const q = `
	SELECT account_id, email, password_hash, role
	FROM accounts
	WHERE role = $1
	`

	rows, err := a.db.QueryContext(ctx, q, role)
	if err != nil {
		return nil, fmt.Errorf("%s, QueryContext: %w", op, err)
	}
	defer rows.Close()

	accounts := make([]*account.Account, 0)

	for rows.Next() {
		acc := &account.Account{}
		err := rows.Scan(&acc.ID, &acc.Email, &acc.PasswordHash, &acc.Role)
		if err != nil {
			return nil, fmt.Errorf("%s, Scan: %w", op, err)
		}
		accounts = append(accounts, acc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s, RowsErr: %w", op, err)
	}

	return accounts, nil
}
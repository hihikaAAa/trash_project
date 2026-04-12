// Package account
package account

import (
	"strings"

	"github.com/google/uuid"
	"github.com/hihikaAAa/trash_project/internal/domain_errors"
)

const availableLocal = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789.-!#$%&*/=?^{|}~_+"
const availableDomain = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM0123456789.-"

type Roles string

const (
	USER   Roles = "USER"
	WORKER Roles = "WORKER"
)

type Account struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         Roles     `json:"role"`
}

func NewAccount(email, passwordHash, role string) (*Account, error) {
	id := uuid.New()
	err := validateEmail(email)

	if err != nil {
		return nil, err
	}

	err = validateRole(role)
	if err != nil {
		return nil, err
	}

	return &Account{ID: id,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         Roles(role),
	}, nil
}

func validateRole(role string) error {
	if role != string(USER) && role != string(WORKER) {
		return domainerrors.ErrWrongRole
	}
	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if len(email) == 0 || len(email) > 254 {
		return domainerrors.ErrBadEmail
	}

	for _, symb := range email {
		if symb <= 31 || symb == 127 || symb == ' ' {
			return domainerrors.ErrBadEmail
		}
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return domainerrors.ErrBadEmail
	}

	local, domain := parts[0], parts[1]
	if len(local) == 0 || len(local) > 64 {
		return domainerrors.ErrBadEmail
	}

	if len(domain) == 0 || len(domain) > 150 {
		return domainerrors.ErrBadEmail
	}

	for _, symb := range local {
		if !strings.Contains(availableLocal, string(symb)) {
			return domainerrors.ErrBadEmail
		}
	}

	for _, symb := range domain {
		if !strings.Contains(availableDomain, string(symb)) {
			return domainerrors.ErrBadEmail
		}
	}

	if strings.Contains(local, "..") || local[0] == '.' || local[len(local)-1] == '.' || local[0] == '-' || local[len(local)-1] == '-' {
		return domainerrors.ErrBadEmail
	}

	if strings.Contains(domain, "..") || domain[0] == '.' || domain[len(domain)-1] == '.' || !strings.Contains(domain, ".") || domain[0] == '-' || domain[len(domain)-1] == '-' {
		return domainerrors.ErrBadEmail
	}
	return nil
}

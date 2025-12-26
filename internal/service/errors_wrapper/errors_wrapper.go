// Package errorswrapper
package errorswrapper

import (
	"errors"
	"fmt"

	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
	postgreserrors "github.com/hihikaAAa/TrashProject/internal/postgres/postgres_errors"
)

func MapRepoErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, postgreserrors.ErrUserNotFound) {
		return domainerrors.ErrUserNotFound
	}
	if errors.Is(err, postgreserrors.ErrUserExists) {
		return domainerrors.ErrUserExists
	}

	if errors.Is(err, postgreserrors.ErrAddressNotFound) {
		return domainerrors.ErrAddressNotFound
	}
	if errors.Is(err, postgreserrors.ErrAddressExists) {
		return domainerrors.ErrAddressExists
	}

	if errors.Is(err, postgreserrors.ErrWorkerNotFound) {
		return domainerrors.ErrWorkerNotFound
	}
	if errors.Is(err, postgreserrors.ErrWorkerExists) {
		return domainerrors.ErrWorkerExists
	}

	if errors.Is(err, postgreserrors.ErrTaskNotFound) {
		return domainerrors.ErrTaskNotFound
	}

	return err
}

func WrapRepoErr(op string, err error) error {
	err = MapRepoErr(err)
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domainerrors.ErrUserNotFound),
		errors.Is(err, domainerrors.ErrUserExists),
		errors.Is(err, domainerrors.ErrAddressNotFound),
		errors.Is(err, domainerrors.ErrAddressExists),
		errors.Is(err, domainerrors.ErrTaskNotFound),
		errors.Is(err, domainerrors.ErrWorkerNotFound),
		errors.Is(err, domainerrors.ErrWorkerExists):
		return err
	default:
		return fmt.Errorf("%s: %w", op, err)
	}
}

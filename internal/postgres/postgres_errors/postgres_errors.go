package repoerrors

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrWorkerExists    = errors.New("worker already exists")
	ErrWorkerNotFound  = errors.New("worker not found")
	ErrAddressNotFound = errors.New("address not found")
	ErrAddressExists   = errors.New("address already exists")
	ErrTaskNotFound    = errors.New("task not found")
)

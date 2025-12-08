package repoerrors 

import "errors"

var(
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists = errors.New("user already exists")
	ErrWorkerExists = errors.New("worker already exists")
	ErrWorkerNotFound = errors.New("worker not found")
)
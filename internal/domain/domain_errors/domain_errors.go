// Package domainerrors
package domainerrors

import (
	"errors"
)

var(
	ErrTaskCanceled = errors.New("task was already canceled")
	ErrTaskDone = errors.New("task was already done")
	ErrTaskNotInProgress = errors.New("can complete only task, which is in progress")
	ErrTaskIsNotOpen = errors.New("can start only task, which is open")
	ErrWorkerNotActive = errors.New("worker not active")
	ErrTaskAlreadyAssigned = errors.New("task already assigned")
	ErrWorkerAlreadyActive = errors.New("worker already active")
	ErrWorkerAlreadyDeactive = errors.New("worker already deactive")
	ErrTaskIsNotFound = errors.New("task is not found")
)
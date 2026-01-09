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
	ErrTaskNotFound = errors.New("task is not found")
	ErrBadTaskInfo = errors.New("bad task info")
	ErrBadTaskTime = errors.New("bad task time")
	ErrBadTaskStatus = errors.New("bad task status")
	ErrBadTaskWorker = errors.New("bad task worker")
	ErrEmptyAccountID = errors.New("empty account id ")

	ErrWorkerNotFound = errors.New("worker not found")
	ErrWorkerExists = errors.New("worker exists")
	ErrWorkerNotActive = errors.New("worker not active")
	ErrTaskAlreadyAssigned = errors.New("task already assigned")
	ErrWorkerAlreadyActive = errors.New("worker already active")
	ErrWorkerAlreadyDeactive = errors.New("worker already deactive")

	ErrUserExists = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrEmptyID = errors.New("empty required id")

	ErrAddressNotFound = errors.New("address not found")
	ErrAddressExists = errors.New("address exists")
	ErrBadFloorNumber = errors.New("bad floor number")
	ErrBadApartmentNumber = errors.New("bad apartment number")
	ErrBadStreet = errors.New("bad street")
	ErrBadHouseNumber = errors.New("bad house number")
	ErrBadEntrance = errors.New("bad entrance")

	ErrBadNamePart = errors.New("bad name")
	ErrBadLastName = errors.New("bad last name")

	ErrWrongRole = errors.New("wrong role")
	ErrBadEmail = errors.New("bad email")
)
// Package task содержит модели и логику работы с задачами.
package task

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
)

var validate = validator.New()

type Status string

type Role string

const (
	StatusOpen       Status = "OPEN"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
	StatusCanceled   Status = "CANCELED"
	WorkerRole Role = "worker"
	UserRole Role = "user"
)

type Task struct {
	ID        uuid.UUID  `json:"id" validate:"required"`
	ClientID  uuid.UUID  `json:"client_id" validate:"required"`
	AddressID uuid.UUID  `json:"address_id" validate:"required"`
	Status    Status     `json:"status"`
	WorkerID  *uuid.UUID `json:"worker_id"`
	CreatedAt time.Time  `json:"created_at"`
	ClosedAt  *time.Time `json:"closed_at"`
}

func NewTask(clientID, addressID uuid.UUID, now time.Time, role Role) (*Task, error) {
	id := uuid.New()

	if !CheckRole(role){
		return nil, domainerrors.ErrWrongRole
	}
	
	t := Task{
		ID:        id,
		ClientID:  clientID,
		AddressID: addressID,
		Status:    StatusOpen,
		CreatedAt: now,
	}

	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &t, nil
}

func CheckRole(role Role) bool{
	return role == UserRole
}

func (t *Task) CheckPossibleStatus() error {
	if t.Status == StatusCanceled {
		return domainerrors.ErrTaskCanceled
	}
	if t.Status == StatusDone {
		return domainerrors.ErrTaskDone
	}
	return nil
}

func (t *Task) StartTask(role Role) error {
	if err := t.CheckPossibleStatus(); err != nil {
		return err
	}
	if t.Status != StatusOpen{
		return domainerrors.ErrTaskIsNotOpen
	}
	if CheckRole(role){
		return domainerrors.ErrWrongRole
	}
	t.Status = StatusInProgress
	return nil
}

func (t *Task) CompleteTask(now time.Time, role Role) error {
	if err := t.CheckPossibleStatus(); err != nil {
		return err
	}
	if t.Status != StatusInProgress{
		return domainerrors.ErrTaskNotInProgress
	}
	if CheckRole(role){
		return domainerrors.ErrWrongRole
	}
	t.Status = StatusDone
	t.ClosedAt = &now
	return nil
}

func (t *Task) CancelTask(now time.Time, role Role) error {
	if err := t.CheckPossibleStatus(); err != nil {
		return err
	}
	if !CheckRole(role){
		return domainerrors.ErrWrongRole
	}
	t.Status = StatusCanceled
	t.ClosedAt = &now
	return nil
}

func (t *Task) DropTask(role Role) error {
	if err := t.CheckPossibleStatus(); err != nil {
		return err
	}
	if t.Status != StatusInProgress{
		return domainerrors.ErrTaskNotInProgress
	}
	if CheckRole(role){
		return domainerrors.ErrWrongRole
	}
	t.Status = StatusOpen
	return nil
}

func (t *Task) Validate() error {
	return validate.Struct(t)
}

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

const (
	StatusOpen       Status = "OPEN"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
	StatusCanceled   Status = "CANCELED"
)

type Task struct {
	ID        uuid.UUID  `json:"id"`
	ClientID  uuid.UUID  `json:"client_id" validate:"required"`
	AddressID uuid.UUID  `json:"address_id" validate:"required"`
	Status    Status     `json:"status"`
	WorkerID  *uuid.UUID `json:"worker_id"`
	CreatedAt time.Time  `json:"created_at"`
	ClosedAt  *time.Time `json:"closed_at"`
}

func NewTask (clientID, addressID uuid.UUID) (*Task, error){
	id := uuid.New()

	t := Task{
		ID: id,
		ClientID: clientID,
		AddressID: addressID,
		Status: StatusOpen,
		CreatedAt: time.Now(),
	}

	if err := t.Validate(); err != nil{
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &t, nil
}

func (t *Task) Validate()error{
	return validate.Struct(t)
}

func (t *Task) CheckPossibleStatus() error{
	if t.Status == StatusCanceled{
		return domainerrors.ErrTaskCanceled
	}
	if t.Status == StatusDone{
		return domainerrors.ErrTaskDone
	}
	return nil
}

func (t *Task) StartTask() error{
	if err := t.CheckPossibleStatus(); err != nil{
		return err
	}
	t.Status = StatusInProgress
	return nil
}

func (t *Task) CompleteTask() error{
	if err := t.CheckPossibleStatus(); err != nil{
		return err
	}
	t.Status = StatusDone
	return nil
}

func (t *Task) CancelTask() error{
	if err := t.CheckPossibleStatus(); err != nil{
		return err
	}
	t.Status = StatusCanceled
	return nil
}

func (t *Task) DropTask() error{
	if err := t.CheckPossibleStatus(); err != nil{
		return err
	}
	t.Status = StatusOpen
	return nil
}
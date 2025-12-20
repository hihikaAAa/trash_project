// Package task содержит модели и логику работы с задачами.
package task

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
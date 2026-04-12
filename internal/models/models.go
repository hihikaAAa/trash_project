// Package dto
package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskInput struct {
	ID        uuid.UUID `json:"id,omitempty"`
	ClientID  uuid.UUID `json:"client_id"`
	AddressID uuid.UUID `json:"address_id"`
	Time      time.Time `json:"time"`
	Role      string    `json:"role"`
}

type TaskOutput struct {
	TaskID uuid.UUID `json:"task_id"`
}

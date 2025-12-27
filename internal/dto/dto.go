// Package dto
package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProfileInput struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Role      string    `json:"role"`
	FirstName string    `json:"first_name"`
	Surname   string    `json:"surname"`
	LastName  string    `json:"last_name,omitempty"`
}

type ProfileOutput struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

type TaskInput struct {
	ID uuid.UUID `json:"id,omitempty"`
	ClientID  uuid.UUID `json:"client_id"`
	AddressID uuid.UUID `json:"address_id"`
	Time      time.Time `json:"time"`
	Role      string    `json:"role"`
}

type TaskOutput struct {
	TaskID uuid.UUID `json:"task_id"`
}

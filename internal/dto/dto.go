// Package dto
package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserInput struct {
	ID              uuid.UUID `json:"id,omitempty"`
	FirstName       string    `json:"first_name"`
	Surname         string    `json:"surname"`
	LastName        string    `json:"last_name"`
	Street          string    `json:"street"`
	HouseNumber     string    `json:"house_number"`
	Entrance        string    `json:"entrance"`
	FloorNumber     int       `json:"floor_number"`
	ApartmentNumber int       `json:"apartment_number"`
	Email           string    `json:"email"`
}

type UserOutput struct {
	UserID uuid.UUID `json:"user_id"`
}

type TaskInput struct {
	ClientID  uuid.UUID
	AddressID uuid.UUID
	Time      time.Time
}

type TaskOutput struct {
	TaskID uuid.UUID
}

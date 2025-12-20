// Package address содержит модели и логику работы с адресами.
package address

import "github.com/google/uuid"

type Address struct {
	ID              uuid.UUID
	Street          string
	HouseNumber     string
	Entrance        string
	FloorNumber     int
	ApartmentNumber int
}
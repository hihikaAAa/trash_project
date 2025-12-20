// Package address содержит модели и логику работы с адресами.
package address

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type Address struct {
	ID              uuid.UUID `json:"id"`
	Street          string    `json:"street" validate:"required,min=1,max=20"`
	HouseNumber     string    `json:"house_number" validate:"required,min=1,max=20"`
	Entrance        string    `json:"entrance,omitempty"`
	FloorNumber     int       `json:"floor_number" validate:"gte=0"`
	ApartmentNumber int       `json:"apartment_number" validate:"gte=0"`
}

func NewAddress(street, houseNumber, entrance string, floorNumber, apartmentNumber int) (*Address, error) {
	id := uuid.New()

	a := Address{
		ID:              id,
		Street:          street,
		HouseNumber:     houseNumber,
		Entrance:        entrance,
		FloorNumber:     floorNumber,
		ApartmentNumber: apartmentNumber,
	}

	if err := a.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return &a, nil
}

func (a *Address) UpdateAddress(street, houseNumber, entrance string, floorNumber, apartmentNumber int) error{
	next := Address{
		ID: a.ID,
		Street: street,
		HouseNumber: houseNumber,
		Entrance: entrance,
		FloorNumber: floorNumber,
		ApartmentNumber: apartmentNumber,
	}
	if err := next.Validate(); err != nil{
		return fmt.Errorf("validate: %w", err)
	}

	*a = next

	return nil
}

func (a *Address) Validate() error {
	return validate.Struct(a)
}

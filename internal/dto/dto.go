// Package dto
package dto

import "github.com/google/uuid"

type Input struct {
	FirstName       string `json:"first_name"`
	Surname         string `json:"surname"`
	LastName        string `json:"last_name"`
	Street          string `json:"street"`
	HouseNumber     string `json:"house_number"`
	Entrance        string `json:"entrance"`
	FloorNumber     int    `json:"floor_number"`
	ApartmentNumber int    `json:"apartment_number"`
	Email           string `json:"email"`
}

type Output struct {
	ID uuid.UUID `json:"id"`
}

package address

import "github.com/google/uuid"

type Address struct {
	ID uuid.UUID
	UserID            uuid.UUID
	Street            string
	House_number      string
	Enterance         string
	Floor_number      int
	Apartament_number int
}
package user

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	FirstName string
	Surname   string
	LastName  string
	AddressID uuid.UUID

	// TODO : Добавить подписку
	// Телефон/email для логина
}

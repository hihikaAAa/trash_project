package user

import (
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/address"
)

type User struct{
	Id uuid.UUID
	First_name string
	Surname string
	Last_name string
	Addresses []address.Address 

	// TODO : Добавить подписку
	// Телефон/email для логина
}

package user

import (
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/address"
)

type User struct{
	id uuid.UUID
	first_name string
	surname string
	last_name string
	address address.Address 

	// TODO : Добавить подписку
}

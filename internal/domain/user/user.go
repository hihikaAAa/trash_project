// Package user содержит модели и логику работы с юзерами.
package user

import (
	"fmt"

	"github.com/hihikaAAa/TrashProject/internal/domain/person"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type User struct {
	ID        uuid.UUID     `json:"id"`
	Person    *person.Person `json:"person" validate:"required"`
	AddressID uuid.UUID     `json:"address_id" validate:"required"`

	// TODO : Добавить подписку
	// Телефон/email для логина
}

func NewUser(person *person.Person, addressID uuid.UUID) (*User, error){
	id := uuid.New()

	u := User{
		ID: id,
		Person: person,
		AddressID: addressID,
	}

	if err := u.Validate(); err != nil{
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &u, nil
}

func (u *User) UpdateUser(person person.Person) error{
	next := User{
		ID: u.ID,
		Person: &person,
		AddressID: u.AddressID,
	}

	if err := next.Validate(); err != nil{
		return fmt.Errorf("validate: %w", err)
	}
	u.Person = next.Person
	return nil
}

func (u *User) Validate() error{
	return validate.Struct(u)
}

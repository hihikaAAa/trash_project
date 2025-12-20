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
	AddressID uuid.UUID     `json:"address_id"`

	// TODO : Добавить подписку
	// Телефон/email для логина
}

func NewUser(name,surname,lastName string, addressID uuid.UUID) (*User, error){
	id := uuid.New()
	
	p, err := person.NewPerson(name,surname,lastName)
	if err != nil{
		return nil, err
	}

	u := User{
		ID: id,
		Person: p,
		AddressID: addressID,
	}

	if err := u.Validate(); err != nil{
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &u, nil
}

func (u *User) UpdateUser(person person.Person, addressID uuid.UUID) error{
	u.Person = &person
	u.AddressID = addressID

	if err := u.Validate(); err != nil{
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

func (u *User) Validate() error{
	return validate.Struct(u)
}

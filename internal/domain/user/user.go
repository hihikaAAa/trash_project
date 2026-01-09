// Package user содержит модели и логику работы с юзерами.
package user

import (
	"github.com/google/uuid"

	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
)

type User struct {
	ID        uuid.UUID     `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Person    *person.Person `json:"person"`
	AddressID uuid.UUID     `json:"address_id"`

	// TODO : Добавить подписку
	// Телефон для логина
}

func NewUser(name,surname, lastName string, addressID, accountID uuid.UUID) (*User, error){
	id := uuid.New()

	p, err := person.NewPerson(name, surname, lastName)
	if err != nil {
		return nil, err
	}

	if err := validateID(addressID); err != nil{
		return nil, err
	}

	if err := validateID(accountID); err != nil{
		return nil, err
	}

	u := User{
		ID: id,
		AccountID: accountID,
		Person: p,
		AddressID: addressID,
	}

	return &u, nil
}

func (u *User) UpdateUser(name,surname,lastName string) error{
	p, err := person.NewPerson(name,surname,lastName)
	if err != nil{
		return err
	}

	u.Person = p
	return nil
}

func validateID(id uuid.UUID) error{
	if id == uuid.Nil{
		return domainerrors.ErrEmptyID
	}
	return nil
}
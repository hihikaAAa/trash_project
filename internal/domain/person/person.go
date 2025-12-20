// Package person
package person

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Person struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=30"`
	Surname   string `json:"surname" validate:"required,min=1,max=30"`
	LastName  string `json:"last_name,omitempty"`
}

func NewPerson(name, surname, lastName string)(*Person ,error){
	p := Person{
		FirstName: name,
		Surname: surname,
		LastName: lastName,
	}

	if err := p.Validate(); err != nil{
		return nil, fmt.Errorf("validate : %w", err)
	}
	return &p, nil
}

func (p *Person) Validate() error{
	return validate.Struct(p)
}
// Package worker содержит модели и логику работы с работниками.
package worker

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/hihikaAAa/TrashProject/internal/domain/person"
	"github.com/hihikaAAa/TrashProject/internal/domain/task"
)

var validate = validator.New()

type Worker struct {
	ID       uuid.UUID      `json:"id"`
	Person   *person.Person `json:"person" validate:"required"`
	TaskList []task.Task    `json:"task_list"`
	IsActive bool           `json:"is_active"`
	// TODO : Добавить район для работы
}

func NewWorker(name, surname, lastName string) (*Worker, error){
	id := uuid.New()

	p, err := person.NewPerson(name,surname,lastName)
	if err != nil{
		return nil, err
	}

	w := Worker{
		ID: id,
		Person: p,
	}

	if err := w.Validate(); err != nil{
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &w, nil
}

func (w *Worker) Validate() error{
	return validate.Struct(w)
}

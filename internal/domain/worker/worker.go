// Package worker содержит модели и логику работы с работниками.
package worker

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
)

var validate = validator.New()

const WorkerRole = "worker"
type Worker struct {
	ID       uuid.UUID      `json:"id"`
	Person   *person.Person `json:"person" validate:"required"`
	TaskList []uuid.UUID    `json:"task_list"`
	IsActive bool           `json:"is_active"`
	// TODO : Добавить район для работы
}

func NewWorker(name,surname,lastName string) (*Worker, error) {
	id := uuid.New()

	p, err := person.NewPerson(name, surname, lastName, WorkerRole)
	if err != nil {
		return nil, err
	}

	w := Worker{
		ID:     id,
		Person: p,
	}

	if err := w.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &w, nil
}

func (w *Worker) UpdateWorker(person person.Person) error {
	next := Worker{
		ID:       w.ID,
		Person:   &person,
		TaskList: w.TaskList,
		IsActive: w.IsActive,
	}

	if err := next.Validate(); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	w.Person = next.Person
	return nil
}

func (w *Worker) AddTask(t uuid.UUID) error {
	if !w.IsActive {
		return domainerrors.ErrWorkerNotActive
	}
	for _, existing := range w.TaskList {
		if existing == t {
			return domainerrors.ErrTaskAlreadyAssigned
		}
	}
	w.TaskList = append(w.TaskList, t)
	return nil
}

func (w *Worker) Activate() error {
	if w.IsActive {
		return domainerrors.ErrWorkerAlreadyActive
	}
	w.IsActive = true
	return nil
}

func (w *Worker) Deactivate() error {
	if !w.IsActive {
		return domainerrors.ErrWorkerAlreadyDeactive
	}
	w.IsActive = false
	return nil
}

func (w *Worker) RemoveTask(taskID uuid.UUID) error {
	for i, existing := range w.TaskList {
		if existing == taskID {
			w.TaskList = append(w.TaskList[:i], w.TaskList[i+1:]...)
			return nil
		}
	}
	return domainerrors.ErrTaskNotFound
}

func (w *Worker) Validate() error {
	return validate.Struct(w)
}

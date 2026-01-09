// Package worker содержит модели и логику работы с работниками.
package worker

import (
	"github.com/google/uuid"

	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
)

type Worker struct {
	ID       uuid.UUID      `json:"id"`
	AccountID uuid.UUID 	`json:"account_id"`
	Person   *person.Person `json:"person"`
	TaskList []uuid.UUID    `json:"task_list"`
	IsActive bool           `json:"is_active"`
	// TODO : Добавить район для работы
}

func NewWorker(name,surname,lastName string, accountID uuid.UUID) (*Worker, error) {
	id := uuid.New()

	if err := validateAccountID(accountID); err != nil{
		return nil, err
	}

	p, err := person.NewPerson(name, surname, lastName)
	if err != nil {
		return nil, err
	}

	w := Worker{
		ID:     id,
		AccountID: accountID,
		Person: p,
	}

	return &w, nil
}

func (w *Worker) UpdateWorker(name,surname,lastName string) error {
	p, err := person.NewPerson(name,surname,lastName)
	if err != nil{
		return err
	}

	w.Person = p
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

func validateAccountID(id uuid.UUID) error{
	if id == uuid.Nil{
		return domainerrors.ErrEmptyAccountID
	}
	return nil
}
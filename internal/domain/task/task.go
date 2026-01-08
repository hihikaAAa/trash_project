// Package task содержит модели и логику работы с задачами.
package task

import (
	"time"

	"github.com/google/uuid"
	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
)

type Status string

type Role string

const (
	StatusOpen       Status = "OPEN"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
	StatusCanceled   Status = "CANCELED"
	WorkerRole Role = "WORKER"
	UserRole Role = "USER"
)

type Task struct {
	ID        uuid.UUID  `json:"id"`
	ClientID  uuid.UUID  `json:"client_id"`
	AddressID uuid.UUID  `json:"address_id"`
	Status    Status     `json:"status"`
	WorkerID  *uuid.UUID `json:"worker_id"`
	CreatedAt time.Time  `json:"created_at"`
	ClosedAt  *time.Time `json:"closed_at"`
}

func NewTask(clientID, addressID uuid.UUID, now time.Time, role Role) (*Task, error) {
	id := uuid.New()

	if !isUser(role){
		return nil, domainerrors.ErrWrongRole
	}
	
	if err := validateTaskInfo(clientID,addressID,now); err != nil{
		return nil, err
	}

	t := Task{
		ID:        id,
		ClientID:  clientID,
		AddressID: addressID,
		Status:    StatusOpen,
		CreatedAt: now,
	}

	return &t, nil
}

func isUser(role Role) bool   { 
	return role == UserRole 
}

func isWorker(role Role) bool { 
	return role == WorkerRole 
}

func (t *Task) StartTask(role Role) error {
	if err := validateTaskStatus(t.Status); err != nil {
		return err
	}
	if err := t.checkNotFinal(); err != nil {
		return err
	}
	if !isWorker(role) {
		return domainerrors.ErrWrongRole
	}
	if t.Status != StatusOpen {
		return domainerrors.ErrTaskIsNotOpen
	}

	if t.WorkerID != nil{
		return domainerrors.ErrBadTaskWorker
	}
	t.Status = StatusInProgress
	return nil
}


func (t *Task) CompleteTask(now time.Time, role Role) error {
	if err := validateTaskStatus(t.Status); err != nil {
		return err
	}
	if err := t.checkNotFinal(); err != nil {
		return err
	}
	if !isWorker(role) {
		return domainerrors.ErrWrongRole
	}
	if t.Status != StatusInProgress {
		return domainerrors.ErrTaskNotInProgress
	}
	if now.IsZero() {
		return domainerrors.ErrBadTaskTime
	}
	if t.WorkerID == nil{
		return domainerrors.ErrBadTaskWorker
	}
	t.Status = StatusDone
	t.ClosedAt = &now
	return nil
}


func (t *Task) CancelTask(now time.Time, role Role) error {
	if err := validateTaskStatus(t.Status); err != nil {
		return err
	}
	if err := t.checkNotFinal(); err != nil {
		return err
	}
	if !isUser(role) {
		return domainerrors.ErrWrongRole
	}
	if now.IsZero() {
		return domainerrors.ErrBadTaskTime
	}
	
	t.Status = StatusCanceled
	t.ClosedAt = &now
	return nil
}


func (t *Task) DropTask(role Role) error {
	if err := validateTaskStatus(t.Status); err != nil {
		return err
	}
	if err := t.checkNotFinal(); err != nil {
		return err
	}
	if !isWorker(role) {
		return domainerrors.ErrWrongRole
	}
	
	if t.Status != StatusInProgress {
		return domainerrors.ErrTaskNotInProgress
	}

	t.WorkerID = nil
	t.Status = StatusOpen
	return nil
}

func (t *Task) checkNotFinal() error {
	switch t.Status {
	case StatusCanceled:
		return domainerrors.ErrTaskCanceled
	case StatusDone:
		return domainerrors.ErrTaskDone
	default:
		return nil
	}
}

func (t *Task) AssignWorker(workerID uuid.UUID) error{
	if workerID == uuid.Nil{
		return domainerrors.ErrBadTaskWorker
	}
	
	if t.WorkerID != nil{
		return domainerrors.ErrBadTaskWorker
	}

	if t.Status != StatusOpen{
		return domainerrors.ErrBadTaskStatus
	}
	t.WorkerID = &workerID
	return nil
}

func validateTaskInfo(clientID, addressID uuid.UUID, now time.Time) error{
	if clientID == uuid.Nil || addressID == uuid.Nil{
		return domainerrors.ErrBadTaskInfo
	}

	if now.IsZero(){
		return domainerrors.ErrBadTaskTime
	}
	return nil
}

func validateTaskStatus(status Status) error {
	switch status {
	case StatusOpen, StatusInProgress, StatusDone, StatusCanceled:
		return nil
	default:
		return domainerrors.ErrBadTaskStatus
	}
}
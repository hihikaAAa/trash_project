// Package task contains order domain model and status transitions.
package task

import (
	"strings"
	"time"

	"github.com/google/uuid"
	domainerrors "github.com/hihikaAAa/trash_project/internal/domainerrors"
)

type Status string

type Role string

const (
	StatusCreated    Status = "created"
	StatusAssigned   Status = "assigned"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusCanceled   Status = "canceled"

	RoleAdmin  Role = "admin"
	RoleWorker Role = "worker"
	RoleUser   Role = "user"
)

type Task struct {
	ID            uuid.UUID  `json:"id"`
	ClientID      uuid.UUID  `json:"client_id"`
	Address       string     `json:"address"`
	Description   *string    `json:"description,omitempty"`
	PreferredTime *time.Time `json:"preferred_time,omitempty"`

	WorkerID *uuid.UUID `json:"worker_id,omitempty"`
	Status   Status     `json:"status"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	AssignedAt  *time.Time `json:"assigned_at,omitempty"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CanceledAt  *time.Time `json:"canceled_at,omitempty"`
}

func NewTask(clientID uuid.UUID, address string, description *string, preferredTime *time.Time, now time.Time, role Role) (*Task, error) {
	if role != RoleUser {
		return nil, domainerrors.ErrWrongRole
	}

	if err := validateCreateInput(clientID, address, description, preferredTime, now); err != nil {
		return nil, err
	}

	return &Task{
		ID:            uuid.New(),
		ClientID:      clientID,
		Address:       strings.TrimSpace(address),
		Description:   normalizeDescription(description),
		PreferredTime: preferredTime,
		Status:        StatusCreated,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (t *Task) AssignWorker(workerID uuid.UUID, now time.Time) error {
	if err := validateNow(now); err != nil {
		return err
	}
	if err := validateWorkerID(workerID); err != nil {
		return err
	}
	if err := t.ensureNotFinal(); err != nil {
		return err
	}

	switch t.Status {
	case StatusCreated, StatusAssigned, StatusInProgress:
		if t.Status == StatusInProgress {
			t.StartedAt = nil
		}
		t.WorkerID = &workerID
		t.Status = StatusAssigned
		t.AssignedAt = &now
		t.UpdatedAt = now
		return nil
	default:
		return domainerrors.ErrInvalidStatusTransition
	}
}

func (t *Task) StartByWorker(workerID uuid.UUID, now time.Time, role Role) error {
	if role != RoleWorker {
		return domainerrors.ErrWrongRole
	}
	if err := validateNow(now); err != nil {
		return err
	}
	if err := validateWorkerID(workerID); err != nil {
		return err
	}
	if t.Status != StatusAssigned {
		return domainerrors.ErrInvalidStatusTransition
	}
	if t.WorkerID == nil || *t.WorkerID != workerID {
		return domainerrors.ErrForbidden
	}

	t.Status = StatusInProgress
	t.StartedAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Task) CompleteByWorker(workerID uuid.UUID, now time.Time, role Role) error {
	if role != RoleWorker {
		return domainerrors.ErrWrongRole
	}
	if err := validateNow(now); err != nil {
		return err
	}
	if err := validateWorkerID(workerID); err != nil {
		return err
	}
	if t.Status != StatusInProgress {
		return domainerrors.ErrTaskNotInProgress
	}
	if t.WorkerID == nil || *t.WorkerID != workerID {
		return domainerrors.ErrForbidden
	}

	t.Status = StatusCompleted
	t.CompletedAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Task) CancelByAdmin(now time.Time, role Role) error {
	if role != RoleAdmin {
		return domainerrors.ErrWrongRole
	}
	if err := validateNow(now); err != nil {
		return err
	}
	if err := t.ensureNotFinal(); err != nil {
		return err
	}

	t.Status = StatusCanceled
	t.CanceledAt = &now
	t.UpdatedAt = now
	return nil
}

func (t *Task) CanBeViewedBy(actorID uuid.UUID, role Role) bool {
	switch role {
	case RoleAdmin:
		return true
	case RoleUser:
		return t.ClientID == actorID
	case RoleWorker:
		if t.WorkerID == nil {
			return false
		}
		return *t.WorkerID == actorID
	default:
		return false
	}
}

func normalizeDescription(v *string) *string {
	if v == nil {
		return nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return nil
	}
	return &s
}

func validateCreateInput(clientID uuid.UUID, address string, description *string, preferredTime *time.Time, now time.Time) error {
	if clientID == uuid.Nil {
		return domainerrors.ErrBadTaskInfo
	}
	if strings.TrimSpace(address) == "" {
		return domainerrors.ErrBadTaskAddress
	}
	if description != nil && len(strings.TrimSpace(*description)) > 2000 {
		return domainerrors.ErrBadTaskDescription
	}
	if preferredTime != nil && preferredTime.IsZero() {
		return domainerrors.ErrBadTaskTime
	}
	return validateNow(now)
}

func validateWorkerID(workerID uuid.UUID) error {
	if workerID == uuid.Nil {
		return domainerrors.ErrBadTaskWorker
	}
	return nil
}

func validateNow(now time.Time) error {
	if now.IsZero() {
		return domainerrors.ErrBadTaskTime
	}
	return nil
}

func (t *Task) ensureNotFinal() error {
	switch t.Status {
	case StatusCompleted:
		return domainerrors.ErrTaskDone
	case StatusCanceled:
		return domainerrors.ErrTaskCanceled
	default:
		return nil
	}
}

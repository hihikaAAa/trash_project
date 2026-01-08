package task

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
)

func TestCreateTask_Success(t *testing.T) {
	clientID, addressID, now := GenerateParams()

	task, err := NewTask(clientID, addressID, now, "USER")

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if task == nil {
		t.Fatal("expected task, got nil")
	}

	if task.ClientID != clientID {
		t.Fatalf("expected ClientID = %s, got %s", clientID, task.ClientID)
	}

	if task.AddressID != addressID {
		t.Fatalf("expected AddressID = %s, got %s", addressID, task.AddressID)
	}

	if task.Status != StatusOpen {
		t.Fatalf("expected Status = StatusOpen, got %s", task.Status)
	}

	if task.ID == uuid.Nil {
		t.Fatal("expected not-nil ID")
	}

	if task.CreatedAt.IsZero() {
		t.Fatal("expected CreatedAt to be set")
	}

	if task.ClosedAt != nil {
		t.Fatalf("expected ClosedAt nil, got %v", *task.ClosedAt)
	}
}

func TestCreateTask_Error_EmptyRequiredParams(t *testing.T) {
	clientID := uuid.Nil
	addressID := uuid.New()
	now := time.Now()

	_, err := NewTask(clientID, addressID, now, "USER")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	clientID = uuid.New()
	addressID = uuid.Nil
	_, err = NewTask(clientID, addressID, now, "user")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

}

func TestTask_StartTask_Success(t *testing.T) {
	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")

	err := task.StartTask("WORKER")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if task.Status != StatusInProgress {
		t.Fatalf("expected Status = %s, got %s", StatusInProgress, task.Status)
	}
}

func TestTask_CompleteTask_Success(t *testing.T) {
	complete := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")
	task.AssignWorker(uuid.New())
	task.Status = StatusInProgress

	err := task.CompleteTask(complete, "WORKER")

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if task.Status != StatusDone {
		t.Fatalf("expected Status = %s, got %s", StatusDone, task.Status)
	}
	if task.ClosedAt == nil {
		t.Fatal("expected ClosedAt not nil")
	}
	if !task.ClosedAt.Equal(complete) {
		t.Fatalf("expected ClosedAt = %v, got %v", complete, *task.ClosedAt)
	}
}

func TestTask_CompleteTask_Error_NotInProgress(t *testing.T) {
	complete := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")

	err := task.CompleteTask(complete, "WORKER")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskNotInProgress) {
		t.Fatalf("expected ErrTaskNotInProgress, got %v", err)
	}
}

func TestTask_CancelTask_Success(t *testing.T) {
	cancel := time.Date(2025, 6, 7, 8, 9, 10, 0, time.UTC)

	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")

	err := task.CancelTask(cancel, "USER")

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if task.Status != StatusCanceled {
		t.Fatalf("expected Status = %s, got %s", StatusCanceled, task.Status)
	}
	if task.ClosedAt == nil {
		t.Fatal("expected ClosedAt not nil")
	}
	if !task.ClosedAt.Equal(cancel) {
		t.Fatalf("expected ClosedAt = %v, got %v", cancel, *task.ClosedAt)
	}
}

func TestTask_DropTask_Success(t *testing.T) {
	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")
	task.Status = StatusInProgress

	err := task.DropTask("WORKER")

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if task.Status != StatusOpen {
		t.Fatalf("expected Status = %s, got %s", StatusOpen, task.Status)
	}
}

func TestTask_DropTask_Error_NotInProgress(t *testing.T) {
	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")

	err := task.DropTask("WORKER")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskNotInProgress) {
		t.Fatalf("expected ErrTaskNotInProgress, got %v", err)
	}
}

func TestTask_StartTask_Error_TaskDone(t *testing.T) {
	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")
	task.Status = StatusDone

	err := task.StartTask("USER")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskDone) {
		t.Fatalf("expected ErrTaskDone, got %v", err)
	}
}

func TestTask_CancelTask_Error_TaskDone(t *testing.T) {
	cancel := time.Date(2025, 6, 7, 8, 9, 10, 0, time.UTC)

	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")
	task.Status = StatusDone

	err := task.CancelTask(cancel, "USER")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskDone) {
		t.Fatalf("expected ErrTaskDone, got %v", err)
	}
}

func TestTask_CompleteTask_Error_TaskCanceled(t *testing.T) {
	complete := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")
	task.Status = StatusCanceled

	err := task.CompleteTask(complete, "WORKER")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskCanceled) {
		t.Fatalf("expected ErrTaskCanceled, got %v", err)
	}
}

func TestTask_DropTask_Error_TaskCanceled(t *testing.T) {
	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")
	task.Status = StatusCanceled

	err := task.DropTask("WORKER")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskCanceled) {
		t.Fatalf("expected ErrTaskCanceled, got %v", err)
	}
}

func TestTask_StartTask_Error_TaskCanceled(t *testing.T) {
	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")
	task.Status = StatusCanceled

	err := task.StartTask("USER")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskCanceled) {
		t.Fatalf("expected ErrTaskCanceled, got %v", err)
	}
}

func TestTask_StartTask_Error_TaskNotOpen(t *testing.T) {
	clientID, addressID, now := GenerateParams()
	task, _ := NewTask(clientID, addressID, now, "USER")
	task.Status = StatusInProgress

	err := task.StartTask("WORKER")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskIsNotOpen) {
		t.Fatalf("expected ErrTaskIsNotOpen, got %v", err)
	}
}

func TestTask_CreateTask_Error_PersonIsWorker(t *testing.T){
	clientID, addressID, now := GenerateParams()
	_, err := NewTask(clientID,addressID,now,"WORKER")
	if err == nil{
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err,domainerrors.ErrWrongRole){
		t.Fatalf("expected ErrWrongRole, got %v", err)
	}
}

func GenerateParams() (uuid.UUID, uuid.UUID, time.Time) {
	return uuid.New(), uuid.New(), time.Now()
}

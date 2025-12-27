package worker

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
)

func TestCreateWorker_Success(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if worker == nil {
		t.Fatal("expected worker, got nil")
	}

	if worker.ID == uuid.Nil {
		t.Fatal("expected not-nil ID")
	}
	if worker.Person == nil {
		t.Fatal("expected Person not nil")
	}

	if worker.Person.FirstName != "Ivan" {
		t.Fatalf("expected FirstName = Ivan, got %s", worker.Person.FirstName)
	}
	if worker.Person.Surname != "Ivanov" {
		t.Fatalf("expected Surname = Ivanov, got %s", worker.Person.Surname)
	}
	if worker.Person.LastName != "Ivanovich" {
		t.Fatalf("expected LastName = Ivanovich, got %s", worker.Person.LastName)
	}

	if worker.IsActive != false {
		t.Fatalf("expected IsActive = false, got %v", worker.IsActive)
	}
	if len(worker.TaskList) != 0 {
		t.Fatalf("expected TaskList empty, got %d", len(worker.TaskList))
	}
}

func TestCreateWorker_Error_EmptyRequiredParams(t *testing.T) {
	_, err := NewWorker("", "Ivanov", "Ivanovich")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	_, err = NewWorker("Ivan", "", "Ivanovich")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateWorker_Success(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	oldID := worker.ID

	newPerson := person.Person{
		FirstName: "Petr",
		Surname:   "Petrov",
		LastName:  "Petrovich",
		Role: "worker",
	}

	err = worker.UpdateWorker(newPerson)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if worker.ID != oldID {
		t.Fatalf("expected ID unchanged = %s, got %s", oldID, worker.ID)
	}

	if worker.Person == nil {
		t.Fatal("expected Person not nil")
	}
	if worker.Person.FirstName != "Petr" {
		t.Fatalf("expected FirstName = Petr, got %s", worker.Person.FirstName)
	}
	if worker.Person.Surname != "Petrov" {
		t.Fatalf("expected Surname = Petrov, got %s", worker.Person.Surname)
	}
	if worker.Person.LastName != "Petrovich" {
		t.Fatalf("expected LastName = Petrovich, got %s", worker.Person.LastName)
	}
}

func TestUpdateWorker_Error_InvalidPerson_StateNotChanged(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	oldFirstName := worker.Person.FirstName
	oldSurname := worker.Person.Surname
	oldLastName := worker.Person.LastName

	invalidPerson := person.Person{
		FirstName: "",
		Surname:   "Petrov",
		LastName:  "Petrovich",
		Role: "worker",
	}

	err = worker.UpdateWorker(invalidPerson)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if worker.Person == nil {
		t.Fatal("expected Person not nil")
	}
	if worker.Person.FirstName != oldFirstName {
		t.Fatalf("expected FirstName unchanged = %s, got %s", oldFirstName, worker.Person.FirstName)
	}
	if worker.Person.Surname != oldSurname {
		t.Fatalf("expected Surname unchanged = %s, got %s", oldSurname, worker.Person.Surname)
	}
	if worker.Person.LastName != oldLastName {
		t.Fatalf("expected LastName unchanged = %s, got %s", oldLastName, worker.Person.LastName)
	}
}

func TestWorker_Activate_Success(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if worker.IsActive != false {
		t.Fatalf("expected IsActive = false, got %v", worker.IsActive)
	}

	err = worker.Activate()

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if worker.IsActive != true {
		t.Fatalf("expected IsActive = true, got %v", worker.IsActive)
	}
}

func TestWorker_Activate_Error_AlreadyActive(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	worker.IsActive = true

	err = worker.Activate()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrWorkerAlreadyActive) {
		t.Fatalf("expected ErrWorkerAlreadyActive, got %v", err)
	}
}

func TestWorker_Deactivate_Success(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	worker.IsActive = true

	err = worker.Deactivate()

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if worker.IsActive != false {
		t.Fatalf("expected IsActive = false, got %v", worker.IsActive)
	}
}

func TestWorker_Deactivate_Error_AlreadyDeactive(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	worker.IsActive = false

	err = worker.Deactivate()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrWorkerAlreadyDeactive) {
		t.Fatalf("expected ErrWorkerAlreadyDeactive, got %v", err)
	}
}

func TestWorker_AddTask_Error_NotActive(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	worker.IsActive = false

	taskID := uuid.New()

	err = worker.AddTask(taskID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrWorkerNotActive) {
		t.Fatalf("expected ErrWorkerNotActive, got %v", err)
	}
	if len(worker.TaskList) != 0 {
		t.Fatalf("expected TaskList empty, got %d", len(worker.TaskList))
	}
}

func TestWorker_AddTask_Success(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	worker.IsActive = true

	taskID := uuid.New()

	err = worker.AddTask(taskID)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(worker.TaskList) != 1 {
		t.Fatalf("expected TaskList size = 1, got %d", len(worker.TaskList))
	}
	if worker.TaskList[0] != taskID {
		t.Fatalf("expected TaskList[0] = %s, got %s", taskID, worker.TaskList[0])
	}
}

func TestWorker_AddTask_Error_TaskAlreadyAssigned(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	worker.IsActive = true

	taskID := uuid.New()

	err = worker.AddTask(taskID)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	err = worker.AddTask(taskID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskAlreadyAssigned) {
		t.Fatalf("expected ErrTaskAlreadyAssigned, got %v", err)
	}
	if len(worker.TaskList) != 1 {
		t.Fatalf("expected TaskList size = 1, got %d", len(worker.TaskList))
	}
}

func TestWorker_RemoveTask_Success(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	worker.IsActive = true

	taskID1 := uuid.New()
	taskID2 := uuid.New()
	taskID3 := uuid.New()

	worker.TaskList = []uuid.UUID{taskID1, taskID2, taskID3}

	err = worker.RemoveTask(taskID2)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if len(worker.TaskList) != 2 {
		t.Fatalf("expected TaskList size = 2, got %d", len(worker.TaskList))
	}
	if worker.TaskList[0] != taskID1 {
		t.Fatalf("expected first task = %s, got %s", taskID1, worker.TaskList[0])
	}
	if worker.TaskList[1] != taskID3 {
		t.Fatalf("expected second task = %s, got %s", taskID3, worker.TaskList[1])
	}
}

func TestWorker_RemoveTask_Error_TaskNotFound(t *testing.T) {
	worker, err := NewWorker("Ivan", "Ivanov", "Ivanovich")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	taskID := uuid.New()
	worker.TaskList = []uuid.UUID{}

	err = worker.RemoveTask(taskID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domainerrors.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskIsNotFound, got %v", err)
	}
}

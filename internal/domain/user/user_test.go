package user

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hihikaAAa/TrashProject/internal/domain/person"
)

func TestCreateUser_Success(t *testing.T) {
	addressID := uuid.New()

	user, err := NewUser("Ivan", "Ivanov", "Ivanovich", addressID)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}

	if user.ID == uuid.Nil {
		t.Fatal("expected not-nil ID")
	}

	if user.Person == nil {
		t.Fatal("expected Person not nil")
	}

	if user.Person.FirstName != "Ivan" {
		t.Fatalf("expected FirstName = Ivan, got %s", user.Person.FirstName)
	}
	if user.Person.Surname != "Ivanov" {
		t.Fatalf("expected Surname = Ivanov, got %s", user.Person.Surname)
	}
	if user.Person.LastName != "Ivanovich" {
		t.Fatalf("expected LastName = Ivanovich, got %s", user.Person.LastName)
	}

	if user.AddressID != addressID {
		t.Fatalf("expected AddressID = %s, got %s", addressID, user.AddressID)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	addressID := uuid.New()

	user, err := NewUser("Ivan", "Ivanov", "Ivanovich", addressID)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	newPerson := person.Person{
		FirstName: "Petr",
		Surname:   "Petrov",
		LastName:  "Petrovich",
		Role: "user",
	}

	err = user.UpdateUser(newPerson)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if user.Person == nil {
		t.Fatal("expected Person not nil")
	}

	if user.Person.FirstName != "Petr" {
		t.Fatalf("expected FirstName = Petr, got %s", user.Person.FirstName)
	}
	if user.Person.Surname != "Petrov" {
		t.Fatalf("expected Surname = Petrov, got %s", user.Person.Surname)
	}
	if user.Person.LastName != "Petrovich" {
		t.Fatalf("expected LastName = Petrovich, got %s", user.Person.LastName)
	}
}

func TestUpdateUser_Error_InvalidPerson_StateNotChanged(t *testing.T) {
	addressID := uuid.New()

	user, err := NewUser("Ivan", "Ivanov", "Ivanovich", addressID)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	oldFirstName := user.Person.FirstName
	oldSurname := user.Person.Surname
	oldLastName := user.Person.LastName
	oldAddressID := user.AddressID

	invalidPerson := person.Person{
		FirstName: "",
		Surname:   "Petrov",
		LastName:  "Petrovich",
	}

	err = user.UpdateUser(invalidPerson)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if user.Person == nil {
		t.Fatal("expected Person not nil")
	}

	if user.Person.FirstName != oldFirstName {
		t.Fatalf("expected FirstName unchanged = %s, got %s", oldFirstName, user.Person.FirstName)
	}
	if user.Person.Surname != oldSurname {
		t.Fatalf("expected Surname unchanged = %s, got %s", oldSurname, user.Person.Surname)
	}
	if user.Person.LastName != oldLastName {
		t.Fatalf("expected LastName unchanged = %s, got %s", oldLastName, user.Person.LastName)
	}
	if user.AddressID != oldAddressID {
		t.Fatalf("expected AddressID unchanged = %s, got %s", oldAddressID, user.AddressID)
	}
}

func TestCreateUser_Error_NilPerson(t *testing.T) {
	addressID := uuid.New()

	_, err := NewUser("","","", addressID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateUser_Error_EmptyAddressID(t *testing.T) {
	addressID := uuid.Nil

	_, err := NewUser("Ivan", "Ivanov", "Ivanovich", addressID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

package user

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateUser_Success(t *testing.T) {
	addressID := uuid.New()
	accountID := uuid.New()

	user, err := NewUser("Иван", "Иванов", "Иванович", addressID,accountID)

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

	if user.Person.FirstName != "Иван" {
		t.Fatalf("expected FirstName = Иван, got %s", user.Person.FirstName)
	}
	if user.Person.Surname != "Иванов" {
		t.Fatalf("expected Surname =  Иванов, got %s", user.Person.Surname)
	}
	if user.Person.LastName != "Иванович" {
		t.Fatalf("expected LastName = Иванович, got %s", user.Person.LastName)
	}

	if user.AddressID != addressID {
		t.Fatalf("expected AddressID = %s, got %s", addressID, user.AddressID)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	addressID := uuid.New()
	accountID := uuid.New()

	user, err := NewUser("Иван", "Иванов", "Иванович", addressID,accountID)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	err = user.UpdateUser("Петр","Петров","Петрович")
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if user.Person == nil {
		t.Fatal("expected Person not nil")
	}

	if user.Person.FirstName != "Петр" {
		t.Fatalf("expected FirstName = Петр, got %s", user.Person.FirstName)
	}
	if user.Person.Surname != "Петров" {
		t.Fatalf("expected Surname = Петров, got %s", user.Person.Surname)
	}
	if user.Person.LastName != "Петрович" {
		t.Fatalf("expected LastName = PПетрович, got %s", user.Person.LastName)
	}
}

func TestUpdateUser_Error_InvalidPerson_StateNotChanged(t *testing.T) {
	addressID := uuid.New()
	accountID := uuid.New()

	user, err := NewUser("Иван", "Иванов", "Иванович", addressID,accountID)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	oldFirstName := user.Person.FirstName
	oldSurname := user.Person.Surname
	oldLastName := user.Person.LastName
	oldAddressID := user.AddressID

	err = user.UpdateUser("","Петров","Петрович")
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
	accountID := uuid.New()

	_, err := NewUser("","","", addressID,accountID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateUser_Error_EmptyAddressID(t *testing.T) {
	addressID := uuid.Nil
	accountID := uuid.New()

	_, err := NewUser("Иван", "Иванов", "Иванович", addressID,accountID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateUser_Error_EmptyAccountIDID(t *testing.T) {
	addressID := uuid.New()
	accountID := uuid.Nil

	_, err := NewUser("Иван", "Иванов", "Иванович", addressID,accountID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

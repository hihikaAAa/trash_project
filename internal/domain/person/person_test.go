package person

import (
	"testing"
)

func TestCreatePerson_Success(t *testing.T) {
	p, err := NewPerson("Ivan", "Ivanov", "Ivanovich", "worker")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p == nil {
		t.Fatal("expected person, got nil")
	}

	if p.FirstName != "Ivan" {
		t.Errorf("expected FirstName = Ivan, got %s", p.FirstName)
	}
	if p.Surname != "Ivanov" {
		t.Errorf("expected Surname = Ivanov, got %s", p.Surname)
	}
	if p.LastName != "Ivanovich" {
		t.Errorf("expected LastName = Ivanovich, got %s", p.LastName)
	}
}

func TestCreatePerson_Success_EmptyLastName(t *testing.T) {
	p, err := NewPerson("Ivan", "Ivanov", "","user")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p == nil {
		t.Fatal("expected person, got nil")
	}

	if p.FirstName != "Ivan" {
		t.Errorf("expected FirstName = Ivan, got %s", p.FirstName)
	}
	if p.Surname != "Ivanov" {
		t.Errorf("expected Surname = Ivanov, got %s", p.Surname)
	}
	if p.LastName != "" {
		t.Errorf("expected LastName is empty, got %s", p.LastName)
	}
}

func TestCreatePerson_Error_BadName(t *testing.T) {
	_, err := NewPerson("", "Ivanov", "Ivanovich", "user")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	_, err = NewPerson("badnamebadnamebadnamebadnamebadname", "Ivanov", "Ivanovich", "user")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestCreatePerson_Error_EmptySurname(t *testing.T) {
	_, err := NewPerson("Ivan", "", "Ivanovich", "user")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	_, err = NewPerson("Ivan", "badfamqbadfamqbadfamqbadfamqbadfamq", "Ivanovich", "user")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

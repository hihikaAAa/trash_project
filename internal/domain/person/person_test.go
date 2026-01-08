package person

import (
	"testing"
)

func TestCreatePerson_Success(t *testing.T) {
	p, err := NewPerson("Иван", "Иванов", "Иванович")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p == nil {
		t.Fatal("expected person, got nil")
	}

	if p.FirstName != "Иван" {
		t.Errorf("expected FirstName = Иван, got %s", p.FirstName)
	}
	if p.Surname != "Иванов" {
		t.Errorf("expected Surname = Иванов, got %s", p.Surname)
	}
	if p.LastName != "Иванович" {
		t.Errorf("expected LastName = Иванович, got %s", p.LastName)
	}
}

func TestCreatePerson_Success_EmptyLastName(t *testing.T) {
	p, err := NewPerson("Иван", "Иванов", "",)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p == nil {
		t.Fatal("expected person, got nil")
	}

	if p.FirstName != "Иван" {
		t.Errorf("expected FirstName = Иван, got %s", p.FirstName)
	}
	if p.Surname != "Иванов" {
		t.Errorf("expected Surname = Иванов, got %s", p.Surname)
	}
	if p.LastName != "" {
		t.Errorf("expected LastName is empty, got %s", p.LastName)
	}
}

func TestCreatePerson_Error_BadName(t *testing.T) {
	_, err := NewPerson("", "Иванов", "Иванович")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	_, err = NewPerson("плохоеимяплохоеимяплохоеимяплохоеимяплохоеимяплохоеимяплохоеимяплохоеимяплохоеимяплохоеимяплохоеимяплохоеимя", "Иванов", "Иванович")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestCreatePerson_Error_EmptySurname(t *testing.T) {
	_, err := NewPerson("Иван", "", "Иванович")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	_, err = NewPerson("Иван", "плохаяфамилияплохаяфамилияплохаяфамилияплохаяфамилияплохаяфамилияплохаяфамилияплохаяфамилияплохаяфамилия", "Иванович")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

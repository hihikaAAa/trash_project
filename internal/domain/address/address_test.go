package address

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateAddress_Success(t *testing.T) {
	street := "Lenina"
	houseNumber := "10A"
	entrance := "2"
	floorNumber := 3
	apartmentNumber := 15

	address, err := NewAddress(street, houseNumber, entrance, floorNumber, apartmentNumber)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if address == nil {
		t.Fatal("expected address, got nil")
	}

	if address.ID == uuid.Nil {
		t.Fatal("expected not-nil ID")
	}

	if address.Street != street {
		t.Fatalf("expected Street = %s, got %s", street, address.Street)
	}
	if address.HouseNumber != houseNumber {
		t.Fatalf("expected HouseNumber = %s, got %s", houseNumber, address.HouseNumber)
	}
	if address.Entrance != entrance {
		t.Fatalf("expected Entrance = %s, got %s", entrance, address.Entrance)
	}
	if address.FloorNumber != floorNumber {
		t.Fatalf("expected FloorNumber = %d, got %d", floorNumber, address.FloorNumber)
	}
	if address.ApartmentNumber != apartmentNumber {
		t.Fatalf("expected ApartmentNumber = %d, got %d", apartmentNumber, address.ApartmentNumber)
	}
}

func TestCreateAddress_Error_EmptyRequiredParams(t *testing.T) {
	street := ""
	houseNumber := "10A"
	entrance := "2"
	floorNumber := 3
	apartmentNumber := 15

	_, err := NewAddress(street, houseNumber, entrance, floorNumber, apartmentNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	street = "Lenina"
	houseNumber = ""
	_, err = NewAddress(street, houseNumber, entrance, floorNumber, apartmentNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateAddress_Error_NegativeNumbers(t *testing.T) {
	street := "Lenina"
	houseNumber := "10A"
	entrance := "2"
	floorNumber := -1
	apartmentNumber := 15

	_, err := NewAddress(street, houseNumber, entrance, floorNumber, apartmentNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	floorNumber = 3
	apartmentNumber = -1
	_, err = NewAddress(street, houseNumber, entrance, floorNumber, apartmentNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateAddress_Success(t *testing.T) {
	street := "Lenina"
	houseNumber := "10A"
	entrance := "2"
	floorNumber := 3
	apartmentNumber := 15

	address, err := NewAddress(street, houseNumber, entrance, floorNumber, apartmentNumber)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	oldID := address.ID

	newStreet := "Tverskaya"
	newHouseNumber := "5"
	newEntrance := "1"
	newFloorNumber := 0
	newApartmentNumber := 1

	err = address.UpdateAddress(newStreet, newHouseNumber, newEntrance, newFloorNumber, newApartmentNumber)

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if address.ID != oldID {
		t.Fatalf("expected ID unchanged = %s, got %s", oldID, address.ID)
	}

	if address.Street != newStreet {
		t.Fatalf("expected Street = %s, got %s", newStreet, address.Street)
	}
	if address.HouseNumber != newHouseNumber {
		t.Fatalf("expected HouseNumber = %s, got %s", newHouseNumber, address.HouseNumber)
	}
	if address.Entrance != newEntrance {
		t.Fatalf("expected Entrance = %s, got %s", newEntrance, address.Entrance)
	}
	if address.FloorNumber != newFloorNumber {
		t.Fatalf("expected FloorNumber = %d, got %d", newFloorNumber, address.FloorNumber)
	}
	if address.ApartmentNumber != newApartmentNumber {
		t.Fatalf("expected ApartmentNumber = %d, got %d", newApartmentNumber, address.ApartmentNumber)
	}
}

func TestUpdateAddress_Error_InvalidParams_StateNotChanged(t *testing.T) {
	street := "Lenina"
	houseNumber := "10A"
	entrance := "2"
	floorNumber := 3
	apartmentNumber := 15

	address, err := NewAddress(street, houseNumber, entrance, floorNumber, apartmentNumber)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	oldID := address.ID
	oldStreet := address.Street
	oldHouseNumber := address.HouseNumber
	oldEntrance := address.Entrance
	oldFloorNumber := address.FloorNumber
	oldApartmentNumber := address.ApartmentNumber

	newStreet := ""
	newHouseNumber := "5"
	newEntrance := "1"
	newFloorNumber := 0
	newApartmentNumber := 1

	err = address.UpdateAddress(newStreet, newHouseNumber, newEntrance, newFloorNumber, newApartmentNumber)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if address.ID != oldID {
		t.Fatalf("expected ID unchanged = %s, got %s", oldID, address.ID)
	}

	if address.Street != oldStreet {
		t.Fatalf("expected Street unchanged = %s, got %s", oldStreet, address.Street)
	}
	if address.HouseNumber != oldHouseNumber {
		t.Fatalf("expected HouseNumber unchanged = %s, got %s", oldHouseNumber, address.HouseNumber)
	}
	if address.Entrance != oldEntrance {
		t.Fatalf("expected Entrance unchanged = %s, got %s", oldEntrance, address.Entrance)
	}
	if address.FloorNumber != oldFloorNumber {
		t.Fatalf("expected FloorNumber unchanged = %d, got %d", oldFloorNumber, address.FloorNumber)
	}
	if address.ApartmentNumber != oldApartmentNumber {
		t.Fatalf("expected ApartmentNumber unchanged = %d, got %d", oldApartmentNumber, address.ApartmentNumber)
	}
}

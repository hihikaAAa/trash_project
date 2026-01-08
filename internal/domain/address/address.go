// Package address содержит модели и логику работы с адресами.
package address

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/google/uuid"
	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
)

const(
	minFloorNumber = 1
	maxFloorNumber = 100

	minApartmentNumber = 1
	maxApartmentNumber = 10000

	minStreetLength = 2
	maxStreetLength = 50

	minHouseNumberLength = 1
	maxHouseNumberLength = 20

	maxEntranceLength = 10
)
type Address struct {
	ID              uuid.UUID `json:"id"`
	Street          string    `json:"street"`
	HouseNumber     string    `json:"house_number"`
	Entrance        string    `json:"entrance,omitempty"`
	FloorNumber     int       `json:"floor_number"`
	ApartmentNumber int       `json:"apartment_number"`
}

func NewAddress(street, houseNumber, entrance string, floorNumber, apartmentNumber int) (*Address, error) {

	deleteSpacesAddress(&street,&houseNumber,&entrance)

	err := validateAddress(street,houseNumber,entrance, floorNumber, apartmentNumber)
	if err != nil{
		return nil, err
	}

	a := Address{
		ID:              uuid.New(),
		Street:          street,
		HouseNumber:     houseNumber,
		Entrance:        entrance,
		FloorNumber:     floorNumber,
		ApartmentNumber: apartmentNumber,
	}


	return &a, nil
}

func (a *Address) UpdateAddress(street, houseNumber, entrance string, floorNumber, apartmentNumber int) error{
	deleteSpacesAddress(&street,&houseNumber,&entrance)

	err := validateAddress(street,houseNumber,entrance, floorNumber, apartmentNumber)

	if err != nil{
		return err
	}

	*a = Address{
		ID: a.ID,
		Street: street,
		HouseNumber: houseNumber,
		Entrance: entrance,
		FloorNumber: floorNumber,
		ApartmentNumber: apartmentNumber,
	}

	return nil
}

func validateFloorNumber(floorNumber int) error{
	if floorNumber < minFloorNumber|| floorNumber > maxFloorNumber{
		return domainerrors.ErrBadFloorNumber
	}
	return nil
}

func validateApartmentNumber(apartmentNumber int) error{
	if apartmentNumber < minApartmentNumber || apartmentNumber > maxApartmentNumber{
		return domainerrors.ErrBadApartmentNumber
	}
	return nil
}

func validateStreet(street string) error{
	n := utf8.RuneCountInString(street)
	if n < minStreetLength || n > maxStreetLength{
		return domainerrors.ErrBadStreet
	}
	for _, char := range street {
		switch {
		case isCyrillicLetter(char):
			continue
		case char == ' ' || char == '-' || char == '.':
			continue
		default:
			return domainerrors.ErrBadStreet
		}
	}
	return nil
}

func validateHouseNumber(houseNumber string) error{
	n := utf8.RuneCountInString(houseNumber)
	if n < minHouseNumberLength || n > maxHouseNumberLength {
		return domainerrors.ErrBadHouseNumber
	}
	for _, r := range houseNumber {
		switch {
		case r >= '0' && r <= '9':
			continue
		case isCyrillicLetter(r):
			continue
		case r == ' ' || r == '-' || r == '/' || r == '.':
			continue
		default:
			return domainerrors.ErrBadHouseNumber
		}
	}
	if strings.HasPrefix(houseNumber, "-") || strings.HasPrefix(houseNumber, "/") || strings.HasPrefix(houseNumber, ".") {
		return domainerrors.ErrBadHouseNumber
	}
	if strings.HasSuffix(houseNumber, "-") || strings.HasSuffix(houseNumber, "/") || strings.HasSuffix(houseNumber, ".") {
		return domainerrors.ErrBadHouseNumber
	}
	return nil
}

func validateEntrance(entrance string) error {
	if entrance == "" {
		return nil
	}
	n := utf8.RuneCountInString(entrance)
	if n > maxEntranceLength {
		return domainerrors.ErrBadEntrance
	}
	for _, r := range entrance {
		switch {
		case r >= '0' && r <= '9':
			continue
		case isCyrillicLetter(r):
			continue
		case r == '-' || r == '/' || r == ' ':
			continue
		default:
			return domainerrors.ErrBadEntrance
		}
	}
	return nil
}

func validateAddress(street, houseNumber, entrance string, floorNumber, apartmentNumber int) error{
	if err := validateStreet(street); err != nil{
		return err
	}
	if err := validateHouseNumber(houseNumber); err != nil{
		return err
	}
	if err := validateEntrance(entrance); err != nil{
		return err
	}
	if err := validateFloorNumber(floorNumber); err != nil{
		return err
	}
	if err := validateApartmentNumber(apartmentNumber); err != nil{
		return err
	}
	return nil
}

func deleteSpacesAddress(street, houseNumber, entrance *string){
	*street = normalizeSpacesAddress(*street)
	*houseNumber = normalizeSpacesAddress(*houseNumber)
	*entrance = normalizeSpacesAddress(*entrance)
}

func normalizeSpacesAddress(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))

	spacePending := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			spacePending = true
			continue
		}
		if spacePending && b.Len() > 0 {
			b.WriteRune(' ')
		}
		spacePending = false
		b.WriteRune(r)
	}
	return b.String()
}

func isCyrillicLetter(r rune) bool {
	return unicode.In(r, unicode.Cyrillic)
}
// Package person
package person

import (
	"strings"
	"unicode"
	"unicode/utf8"

	domainerrors "github.com/hihikaAAa/TrashProject/internal/domain/domain_errors"
)

const (
	minLength = 2
	maxLength = 40
)

type Person struct {
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
	LastName  string `json:"last_name,omitempty"`
}

func NewPerson(name, surname, lastName string) (*Person, error) {
	deleteSpacesPerson(&name,&surname,&lastName)

	err := validatePerson(name,surname,lastName)
	if err != nil{
		return nil, err
	}

	p := Person{
		FirstName: name,
		Surname:   surname,
		LastName:  lastName,
	}

	return &p, nil
}

func (p *Person) UpdatePerson(name, surname, lastName string) error {
	deleteSpacesPerson(&name,&surname,&lastName)

	err := validatePerson(name,surname,lastName)
	if err != nil{
		return err
	}

	p.FirstName = name
	p.Surname = surname
	p.LastName = lastName

	return nil
}

func validateNamePart(part string) error{
	n := utf8.RuneCountInString(part)
	if n < minLength || n > maxLength{
		return domainerrors.ErrBadNamePart
	}
	for _, char := range part{
		switch{
		case isCyrillicLetter(char):
			continue
		case char == ' ' || char == '-':
			continue
		default:
			return domainerrors.ErrBadNamePart
		}
	}

	partSlice := []rune(part)
	if partSlice[0] == '-' || partSlice[len(partSlice)-1] == '-'{
		return domainerrors.ErrBadNamePart
	}

	if strings.Contains(part,"--") || strings.Contains(part, "- ") || strings.Contains(part," -"){
		return domainerrors.ErrBadNamePart
	}

	return nil
}

func validateLastName(lastName string) error{
	if len(lastName) == 0{
		return nil
	}
	n := utf8.RuneCountInString(lastName)
	if n > maxLength{
		return domainerrors.ErrBadNamePart
	}
	for _, char := range lastName{
		switch{
		case isCyrillicLetter(char):
			continue
		case char == ' ' || char == '-':
			continue
		default:
			return domainerrors.ErrBadNamePart
		}
	}

	partSlice := []rune(lastName)
	if partSlice[0] == '-' || partSlice[len(partSlice)-1] == '-'{
		return domainerrors.ErrBadNamePart
	}

	if strings.Contains(lastName,"--") || strings.Contains(lastName, "- ") || strings.Contains(lastName," -"){
		return domainerrors.ErrBadNamePart
	}
	return nil
}

func validatePerson(name, surname, lastName string) error{
	if err := validateNamePart(name); err != nil{
		return err
	}
	if err := validateNamePart(surname); err != nil{
		return err
	}
	if err := validateLastName(lastName); err != nil{
		return err
	}

	return nil
}

func deleteSpacesPerson(firstName, surname, lastName *string){
	*firstName = normalizeSpacesPerson(*firstName)
	*surname = normalizeSpacesPerson(*surname)
	*lastName = normalizeSpacesPerson(*lastName)
}

func normalizeSpacesPerson(s string) string {
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

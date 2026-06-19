package models

import (
	"regexp"
	"unicode/utf8"
)

type Lead struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var (
	nameRegex  = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s\-]+$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9][0-9]{6,14}$`)
)

func CreateLead(name string, phone string) Lead {
	return Lead{
		Name:  name,
		Phone: phone,
	}
}

func (l *Lead) Validate() error {
	if l.Name == "" || l.Phone == "" {
		return ErrEmptyField
	}

	nameLength := utf8.RuneCountInString(l.Name)

	if nameLength < 2 {
		return ErrTooShortUserName
	}

	if nameLength > 20 {
		return ErrTooLongUserName
	}

	if !nameRegex.MatchString(l.Name) {
		return ErrInvalidName
	}

	if !phoneRegex.MatchString(l.Phone) {
		return ErrInvalidPhone
	}

	return nil
}

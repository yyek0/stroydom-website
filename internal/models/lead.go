package models

import (
	"regexp"
	"time"
	"unicode/utf8"
)

type Lead struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"time"`
}

var (
	nameRegex  = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s\-]+$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9][0-9]{6,14}$`)
)

func CreateLead(name string, phone string) Lead {
	return Lead{
		Name:      name,
		Phone:     phone,
		CreatedAt: time.Now(),
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

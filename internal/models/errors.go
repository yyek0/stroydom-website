package models

import "errors"

var (
	ErrInvalidName      error = errors.New("invalid user name")
	ErrInvalidPhone     error = errors.New("invalid user phone number")
	ErrEmptyField       error = errors.New("required field is empty")
	ErrTooShortUserName error = errors.New("name is too short")
	ErrTooLongUserName  error = errors.New("name is too long")
)

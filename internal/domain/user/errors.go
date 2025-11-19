package user

import "errors"

var (
	ErrInvalidEmailFormat = errors.New("email is in invalid format")
	ErrPasswordTooShort   = errors.New("password is too short")
	ErrPasswordTooLong    = errors.New("password is too long")
	ErrPasswordTooWeak    = errors.New("password should contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

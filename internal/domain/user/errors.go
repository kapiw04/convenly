package user

import "errors"

var (
	ErrInvalidEmailFormat = errors.New("Email is in invalid format")
	ErrPasswordTooShort   = errors.New("Password is too short")
	ErrPasswordTooLong    = errors.New("Password is too long")
	ErrPasswordTooWeak    = errors.New("Password should contain at least one uppercase letter, one lowercase letter, one digit, and one special character")
	ErrInvalidCredentials = errors.New("Invalid email or password")
)

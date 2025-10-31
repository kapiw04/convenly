package user

import (
	"net/mail"
	"strings"
)

type Email string

func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrInvalidEmailFormat
	}

	addr, err := mail.ParseAddress(raw)
	if err != nil {
		return "", err
	}

	return Email(strings.ToLower(addr.Address)), nil
}

func (e Email) String() string {
	return string(e)
}

func (e Email) Equal(other Email) bool {
	return strings.EqualFold(string(e), string(other))
}

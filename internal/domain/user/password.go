package user

import (
	"regexp"
	"strings"
)

const MaxPasswordLength = 20
const MinPasswordLength = 8

type Password string

func NewPassword(raw string) (Password, error) {
	raw = strings.TrimSpace(raw)
	err := ValidateLength(raw)
	if err != nil {
		return "", err
	}
	err = ValidateStrength(raw)
	if err != nil {
		return "", err
	}
	return Password(raw), nil
}

func ValidateLength(raw string) error {
	if len(raw) < MinPasswordLength {
		return ErrPasswordTooShort
	}
	if len(raw) > MaxPasswordLength {
		return ErrPasswordTooLong
	}
	return nil
}

func ValidateStrength(raw string) error {
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(raw)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(raw)
	hasDigit := regexp.MustCompile(`\d`).MatchString(raw)
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*()+|_]{1,}`).MatchString(raw)
	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		return ErrPasswordTooWeak
	}
	return nil
}

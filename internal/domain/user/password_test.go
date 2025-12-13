package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPassword_Valid(t *testing.T) {
	password, err := NewPassword("Secret123!")
	require.NoError(t, err)
	require.Equal(t, Password("Secret123!"), password)
}

func TestNewPassword_TooShort(t *testing.T) {
	_, err := NewPassword("Sec1!")
	require.Error(t, err)
	require.Equal(t, ErrPasswordTooShort, err)
}

func TestNewPassword_TooLong(t *testing.T) {
	_, err := NewPassword("ExtreeeemeeelyLongPassword1234!")
	require.Error(t, err)
	require.Equal(t, ErrPasswordTooLong, err)
}

func TestNewPassword_MissingLowercase(t *testing.T) {
	_, err := NewPassword("SECRET123!")
	require.Error(t, err)
	require.Equal(t, ErrPasswordTooWeak, err)
}

func TestNewPassword_MissingUppercase(t *testing.T) {
	_, err := NewPassword("secret123!")
	require.Error(t, err)
	require.Equal(t, ErrPasswordTooWeak, err)
}

func TestNewPassword_MissingDigit(t *testing.T) {
	_, err := NewPassword("SecretPass!")
	require.Error(t, err)
	require.Equal(t, ErrPasswordTooWeak, err)
}

func TestNewPassword_MissingSpecialChar(t *testing.T) {
	_, err := NewPassword("Secret1234")
	require.Error(t, err)
	require.Equal(t, ErrPasswordTooWeak, err)
}

func TestValidateLength_BoundaryMin(t *testing.T) {
	err := ValidateLength("12345678")
	require.NoError(t, err)
}

func TestValidateLength_BoundaryMax(t *testing.T) {
	err := ValidateLength("12345678901234567890")
	require.NoError(t, err)
}

func TestValidateLength_JustBelowMin(t *testing.T) {
	err := ValidateLength("1234567")
	require.Error(t, err)
	require.Equal(t, ErrPasswordTooShort, err)
}

func TestValidateLength_JustAboveMax(t *testing.T) {
	err := ValidateLength("123456789012345678901")
	require.Error(t, err)
	require.Equal(t, ErrPasswordTooLong, err)
}

func TestValidateStrength_AllSpecialChars(t *testing.T) {
	specialChars := []string{"!", "@", "#", "~", "$", "%", "^", "&", "*", "(", ")", "+", "|", "_"}
	for _, char := range specialChars {
		err := ValidateStrength("Secret1" + char)
		require.NoError(t, err, "special char %s should be valid", char)
	}
}

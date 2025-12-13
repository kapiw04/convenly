package user

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEmail_Valid(t *testing.T) {
	email, err := NewEmail("test@example.com")
	require.NoError(t, err)
	require.Equal(t, Email("test@example.com"), email)
}

func TestNewEmail_TrimsWhitespace(t *testing.T) {
	email, err := NewEmail("  test@example.com  ")
	require.NoError(t, err)
	require.Equal(t, Email("test@example.com"), email)
}

func TestNewEmail_LowerCase(t *testing.T) {
	email, err := NewEmail("TEST@EXAMPLE.COM")
	require.NoError(t, err)
	require.Equal(t, Email("test@example.com"), email)
}

func TestNewEmail_Empty(t *testing.T) {
	_, err := NewEmail("")
	require.Error(t, err)
	require.Equal(t, ErrInvalidEmailFormat, err)
}

func TestNewEmail_OnlyWhitespace(t *testing.T) {
	_, err := NewEmail("   ")
	require.Error(t, err)
	require.Equal(t, ErrInvalidEmailFormat, err)
}

func TestNewEmail_InvalidFormat(t *testing.T) {
	_, err := NewEmail("not-an-email")
	require.Error(t, err)
}

func TestNewEmail_MissingDomain(t *testing.T) {
	_, err := NewEmail("test@")
	require.Error(t, err)
}

func TestNewEmail_MissingLocal(t *testing.T) {
	_, err := NewEmail("@example.com")
	require.Error(t, err)
}

func TestEmail_String(t *testing.T) {
	email := Email("test@example.com")
	require.Equal(t, "test@example.com", email.String())
}

func TestEmail_Equal(t *testing.T) {
	email1 := Email("test@example.com")
	email2 := Email("TEST@EXAMPLE.COM")
	require.True(t, email1.Equal(email2))
}

func TestEmail_NotEqual(t *testing.T) {
	email1 := Email("test@example.com")
	email2 := Email("other@example.com")
	require.False(t, email1.Equal(email2))
}

package security

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcryptHasher_HashAndCompare(t *testing.T) {
	hasher := &BcryptHasher{}

	password := "Secret123!"
	hash, err := hasher.Hash(password)

	require.NoError(t, err)
	require.NotEmpty(t, hash)
	require.NotEqual(t, password, hash)
	require.True(t, hasher.Compare(password, hash))
}

func TestBcryptHasher_CompareFails_WrongPassword(t *testing.T) {
	hasher := &BcryptHasher{}

	password := "Secret123!"
	hash, err := hasher.Hash(password)
	require.NoError(t, err)

	require.False(t, hasher.Compare("WrongPassword!", hash))
}

func TestBcryptHasher_CompareFails_InvalidHash(t *testing.T) {
	hasher := &BcryptHasher{}

	require.False(t, hasher.Compare("Secret123!", "invalid-hash"))
}

func TestBcryptHasher_DifferentPasswordsProduceDifferentHashes(t *testing.T) {
	hasher := &BcryptHasher{}

	hash1, err := hasher.Hash("Password1!")
	require.NoError(t, err)

	hash2, err := hasher.Hash("Password2!")
	require.NoError(t, err)

	require.NotEqual(t, hash1, hash2)
}

func TestBcryptHasher_SamePasswordProducesDifferentHashes(t *testing.T) {
	hasher := &BcryptHasher{}
	password := "Secret123!"

	hash1, err := hasher.Hash(password)
	require.NoError(t, err)

	hash2, err := hasher.Hash(password)
	require.NoError(t, err)

	require.NotEqual(t, hash1, hash2)
	require.True(t, hasher.Compare(password, hash1))
	require.True(t, hasher.Compare(password, hash2))
}

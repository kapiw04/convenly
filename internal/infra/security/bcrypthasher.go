package security

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func (bh *BcryptHasher) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func (bh *BcryptHasher) Compare(password string, hash string) bool {
	hashBytes := []byte(hash)
	passwordBytes := []byte(password)

	err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)	
	return err == nil
}

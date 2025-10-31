package security

//go:generate mockgen -destination=./mocks/mock_hasher.go . Hasher
type Hasher interface {
	Hash(string) (string, error)
	Compare(hash string, pw string) bool
}

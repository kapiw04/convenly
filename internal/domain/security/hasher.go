package security

type Hasher interface {
	Hash(string) (string, error)
	Compare(hash string, pw string) bool
}

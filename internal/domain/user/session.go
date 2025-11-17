package user

//go:generate mockgen -destination=./mocks/mock_sessionrepo.go . SessionRepo

type SessionRepo interface {
	Create(email string) (id string, err error)
	Get(id string) (User, error)
	Delete(id string) error
}

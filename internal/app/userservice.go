package app

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/kapiw04/convenly/internal/domain"
)

type UserService struct {
	repo domain.UserRepo
}

func NewUserService(repo domain.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(name string, email string) error {
	userUUID := uuid.NewString()
	slog.Info("Registering user with name: %s, email: %s", name, email)
	err := s.repo.Save(&domain.User{
		UUID:  userUUID,
		Name:  name,
		Email: email,
		Role:  domain.ATTENDEE,
	})
	if err != nil {
		slog.Error("Failed to save user: %v", "err", err)
		return err
	}
	slog.Info("User registered successfully with UUID: %s", "uuid", userUUID)
	return nil
}

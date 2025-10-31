package app

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/kapiw04/convenly/internal/domain/security"
	"github.com/kapiw04/convenly/internal/domain/user"
)

type UserService struct {
	repo user.UserRepo
	h    security.Hasher
}

func NewUserService(repo user.UserRepo, h security.Hasher) *UserService {
	return &UserService{repo: repo, h: h}
}

func (s *UserService) Register(name string, rawEmail string, rawPassword string) error {
	userUUID := uuid.New()
	email, err := user.NewEmail(rawEmail)
	if err != nil {
		return err
	}
	password, err := user.NewPassword(rawPassword)
	if err != nil {
		return err
	}

	passwordHash, err := s.h.Hash(string(password))
	if err != nil {
		return err
	}

	slog.Info("Registering user with id: %s, name: %s, rawEmail: %s", "id", userUUID, "name", name, "email", string(email))

	err = s.repo.Save(&user.User{
		UUID:         userUUID,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         user.ATTENDEE,
	})
	if err != nil {
		slog.Error("Failed to save user: %v", "err", err)
		return err
	}
	slog.Info("User registered successfully with UUID: %s", "uuid", userUUID)
	return nil
}

func (s *UserService) GetByEmail(email string) (*user.User, error) {
	slog.Info("Getting user with email: %s", "email", email)
	return s.repo.FindByEmail(email)
}

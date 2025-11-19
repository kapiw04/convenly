package user

import (
	"errors"

	"github.com/google/uuid"
)

type Role int

//go:generate mockgen -destination=./mocks/mock_userrepo.go . UserRepo

const (
	ATTENDEE Role = iota
	HOST
)

type User struct {
	UUID         uuid.UUID `json:"uuid"`
	Name         string    `json:"name"`
	Email        Email     `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
}

type UserRepo interface {
	Save(user *User) error
	FindByUUID(uuid string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll() ([]*User, error)
	DeleteByUUID(uuid string) error
	Update(user *User) error
	Count() (int, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

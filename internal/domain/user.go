package domain

import "errors"

type Role int

const (
	ATTENDEE = iota
	HOST
)

type User struct {
	UUID  string
	Name  string
	Email string
	Role  Role
}

type UserRepo interface {
	Save(user *User) error
	FindByUUID(uuid string) (*User, error)
	FindAll() ([]*User, error)
	DeleteByUUID(uuid string) error
	Update(user *User) error
	Count() (int, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

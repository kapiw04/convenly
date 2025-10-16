package db

import (
	"database/sql"

	"github.com/kapiw04/convenly/internal/domain"
)

type PostgresUserRepo struct {
	db *sql.DB
}

// Count implements domain.UserRepo.
func (r *PostgresUserRepo) Count() (int, error) {
	panic("unimplemented")
}

// DeleteByUUID implements domain.UserRepo.
func (r *PostgresUserRepo) DeleteByUUID(uuid string) error {
	panic("unimplemented")
}

// FindAll implements domain.UserRepo.
func (r *PostgresUserRepo) FindAll() ([]*domain.User, error) {
	panic("unimplemented")
}

// FindByUUID implements domain.UserRepo.
func (r *PostgresUserRepo) FindByUUID(uuid string) (*domain.User, error) {
	panic("unimplemented")
}

// Update implements domain.UserRepo.
func (r *PostgresUserRepo) Update(user *domain.User) error {
	panic("unimplemented")
}

func NewPostgresUserRepo(db *sql.DB) domain.UserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Save(user *domain.User) error {
	query := "INSERT INTO users (UUID, Name, Email, Role) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, user.UUID, user.Name, user.Email, user.Role)
	return err
}

var _ domain.UserRepo = (*PostgresUserRepo)(nil)

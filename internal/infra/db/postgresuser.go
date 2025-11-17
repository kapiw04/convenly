package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/kapiw04/convenly/internal/domain/user"
)

type PostgresUserRepo struct {
	DB *sql.DB
}

func (r *PostgresUserRepo) FindByEmail(email string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	query := "SELECT user_id, name, email, password_hash, role FROM users WHERE users.email = $1"
	rows, err := r.DB.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	var user user.User

	rows.Next()
	if err := rows.Scan(&user.UUID, &user.Name, &user.Email, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepo) Count() (int, error) {
	panic("unimplemented")
}

func (r *PostgresUserRepo) DeleteByUUID(uuid string) error {
	panic("unimplemented")
}

func (r *PostgresUserRepo) FindAll() ([]*user.User, error) {
	panic("unimplemented")
}

func (r *PostgresUserRepo) FindByUUID(uuid string) (*user.User, error) {
	panic("unimplemented")
}

func (r *PostgresUserRepo) Update(user *user.User) error {
	panic("unimplemented")
}

func NewPostgresUserRepo(db *sql.DB) user.UserRepo {
	return &PostgresUserRepo{DB: db}
}

func (r *PostgresUserRepo) Save(user *user.User) error {
	email := string(user.Email)
	query := "INSERT INTO users (user_id, name, email, password_hash, role) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.DB.Exec(query, user.UUID, user.Name, email, user.PasswordHash, user.Role)
	return err
}

var _ user.UserRepo = (*PostgresUserRepo)(nil)

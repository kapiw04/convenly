package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/kapiw04/convenly/internal/domain/user"
	"github.com/lib/pq"
)

type PostgresUserRepo struct {
	DB *sql.DB
}

func mapPgErr(err error) error {
	var pqe *pq.Error
	if !errors.As(err, &pqe) {
		return err
	}

	switch string(pqe.Code) {
	case "23505": // unique_violation
		switch pqe.Constraint {
		case "users_email_key", "users_email_lower_key":
			return user.ErrUserExists
		default:
			return err
		}

	case "23514": // check_violation
		switch pqe.Constraint {
		case "users_email_format":
			return user.ErrInvalidEmailFormat
		case "users_name_len":
			return user.ErrUsernameTooShort
		default:
			return err
		}
	}

	return err
}
func (r *PostgresUserRepo) FindByEmail(email string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	query := "SELECT user_id, name, email, password_hash, role FROM users WHERE users.email = $1"
	rows, err := r.DB.QueryContext(ctx, query, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var user user.User
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	query := "SELECT user_id, name, email, password_hash, role FROM users WHERE users.user_id = $1"
	rows, err := r.DB.QueryContext(ctx, query, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	var user user.User
	if err := rows.Scan(&user.UUID, &user.Name, &user.Email, &user.PasswordHash, &user.Role); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepo) Update(user *user.User) error {
	email := string(user.Email)
	query := "UPDATE users SET name=$1, email=$2, password_hash=$3, role=$4 WHERE user_id=$5"
	_, err := r.DB.Exec(query, user.Name, email, user.PasswordHash, user.Role, user.UUID)
	return mapPgErr(err)
}

func NewPostgresUserRepo(db *sql.DB) user.UserRepo {
	return &PostgresUserRepo{DB: db}
}

func (r *PostgresUserRepo) Save(user *user.User) error {
	query := "INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4)"
	_, err := r.DB.Exec(query, user.Name, user.Email, user.PasswordHash, user.Role)
	return mapPgErr(err)
}

var _ user.UserRepo = (*PostgresUserRepo)(nil)

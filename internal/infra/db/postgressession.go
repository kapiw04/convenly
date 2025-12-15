package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io"
	"time"

	"github.com/kapiw04/convenly/internal/domain/user"
)

type PostgresSessionRepo struct {
	DB       *sql.DB
	UserRepo user.UserRepo
}

func (p *PostgresSessionRepo) Create(email string) (id string, err error) {
	query := "INSERT INTO sessions (user_id, session_id) VALUES ($1, $2)"
	user, err := p.UserRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	sessionID := generateSessionID()
	userID := user.UUID
	_, err = p.DB.Exec(query, userID, sessionID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (p *PostgresSessionRepo) Delete(sessionID string) error {
	query := "DELETE FROM sessions WHERE session_id = $1"
	_, err := p.DB.Exec(query, sessionID)
	return err
}

func (p *PostgresSessionRepo) Get(sessionID string) (user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	query := "SELECT user_id FROM sessions WHERE sessions.session_id = $1"
	rows, err := p.DB.QueryContext(ctx, query, sessionID)

	if err != nil {
		return user.User{}, err
	}
	var userID string

	rows.Next()
	if err := rows.Scan(&userID); err != nil {
		return user.User{}, err
	}
	user, err := p.UserRepo.FindByUUID(userID)
	return *user, err
}

var _ user.SessionRepo = (*PostgresSessionRepo)(nil)

func NewPostgresSessionRepo(db *sql.DB, userRepo user.UserRepo) user.SessionRepo {
	return &PostgresSessionRepo{DB: db, UserRepo: userRepo}
}

func generateSessionID() string {
	id := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		panic("failed to generate session id")
	}

	return base64.RawURLEncoding.EncodeToString(id)
}

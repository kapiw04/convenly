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

	sessionId := generateSessionId()
	userId := user.UUID
	_, err = p.DB.Exec(query, userId, sessionId)
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (p *PostgresSessionRepo) Delete(sessionId string) error {
	query := "DELETE FROM sessions WHERE session_id = $1"
	_, err := p.DB.Exec(query, sessionId)
	return err
}

func (p *PostgresSessionRepo) Get(sessionId string) (user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	query := "SELECT user_id FROM sessions WHERE sessions.session_id = $1"
	rows, err := p.DB.QueryContext(ctx, query, sessionId)

	if err != nil {
		return user.User{}, err
	}
	var userId string

	rows.Next()
	if err := rows.Scan(&userId); err != nil {
		return user.User{}, err
	}
	user, err := p.UserRepo.FindByUUID(userId)
	return *user, err
}

var _ user.SessionRepo = (*PostgresSessionRepo)(nil)

func NewPostgresSessionRepo(db *sql.DB, userRepo user.UserRepo) user.SessionRepo {
	return &PostgresSessionRepo{DB: db, UserRepo: userRepo}
}

func generateSessionId() string {
	id := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		panic("failed to generate session id")
	}

	return base64.RawURLEncoding.EncodeToString(id)
}

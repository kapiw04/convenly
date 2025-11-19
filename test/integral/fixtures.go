package integral

import (
	"database/sql"
	"testing"

	"github.com/kapiw04/convenly/internal/app"
	"github.com/kapiw04/convenly/internal/infra/db"
	"github.com/kapiw04/convenly/internal/infra/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDb(t *testing.T) *sql.DB {
	t.Helper()

	pgConn, err := StartPostgres()
	assert.NoError(t, err)
	sqlDb, err := sql.Open("postgres", pgConn.DSN)
	assert.NoError(t, err)

	ApplyMigrations(t, sqlDb)
	return sqlDb
}

func setupUserService(t *testing.T, sqlDb *sql.DB) *app.UserService {
	t.Helper()
	hasher := &security.BcryptHasher{}
	pgUserRepo := db.NewPostgresUserRepo(sqlDb)
	pgSessionRepo := &db.PostgresSessionRepo{DB: sqlDb, UserRepo: pgUserRepo}
	return app.NewUserService(pgUserRepo, pgSessionRepo, hasher)
}

func setupEventService(t *testing.T, sqlDb *sql.DB) *app.EventService {
	t.Helper()
	pgEventRepo := db.NewPostgresEventRepo(sqlDb)
	return app.NewEventService(pgEventRepo)
}

func RegisterAndLoginUser(t *testing.T, userSrvc *app.UserService, name, email, password string) string {
	t.Helper()
	err := userSrvc.Register(name, email, password)
	require.NoError(t, err)

	sessionID, err := userSrvc.Login(email, password)
	require.NoError(t, err)
	return sessionID
}

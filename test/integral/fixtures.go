package integral

import (
	"database/sql"
	"testing"

	"github.com/kapiw04/convenly/internal/app"
	"github.com/kapiw04/convenly/internal/infra/db"
	"github.com/kapiw04/convenly/internal/infra/security"
	"github.com/kapiw04/convenly/internal/infra/webapi"
	"github.com/stretchr/testify/require"
)

func setupDb(t *testing.T) *sql.DB {
	t.Helper()
	require.NotNil(t, sqlDB)
	return sqlDB
}

func setupUserService(t *testing.T, dbConn *sql.DB) *app.UserService {
	t.Helper()

	hasher := &security.BcryptHasher{}
	pgUserRepo := db.NewPostgresUserRepo(dbConn)
	pgSessionRepo := &db.PostgresSessionRepo{
		DB:       dbConn,
		UserRepo: pgUserRepo,
	}

	return app.NewUserService(pgUserRepo, pgSessionRepo, hasher)
}

func setupEventService(t *testing.T, dbConn *sql.DB) *app.EventService {
	t.Helper()

	pgTagRepo := db.NewPostgresTagRepo(dbConn)
	pgEventRepo := db.NewPostgresEventRepo(dbConn, pgTagRepo)

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

func setupAllServices(t *testing.T) (*sql.DB, *app.UserService, *app.EventService, *webapi.Router) {
	t.Helper()

	dbConn := setupDb(t)
	userSrvc := setupUserService(t, dbConn)
	eventSrvc := setupEventService(t, dbConn)
	router := webapi.NewRouter(userSrvc, eventSrvc)

	return dbConn, userSrvc, eventSrvc, router
}

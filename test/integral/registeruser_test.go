package integral

import (
	"database/sql"
	"testing"

	"github.com/kapiw04/convenly/internal/app"
	"github.com/kapiw04/convenly/internal/infra/db"
	"github.com/kapiw04/convenly/internal/infra/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestUserRepo_SaveAndGet(t *testing.T) {
	pgConn, err := StartPostgres()
	assert.NoError(t, err)
	sqlDb, err := sql.Open("postgres", pgConn.DSN)
	assert.NoError(t, err)

	ApplyMigrations(t, sqlDb)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hasher := &security.BcryptHasher{}
		pgUserRepo := db.NewPostgresUserRepo(sqlDb, hasher)
		srvc := app.NewUserService(pgUserRepo)
		err := srvc.Register(
			"Alice",
			"alice@example.com",
			"Secret123!",
		)
		require.NoError(t, err)
		u, err := srvc.GetByEmail("alice@example.com")
		assert.NoError(t, err)
		assert.Equal(t, string(u.Email), "alice@example.com")
		assert.Equal(t, u.Name, "Alice")
		assert.Equal(t, len(u.Password), 0)
	})
}

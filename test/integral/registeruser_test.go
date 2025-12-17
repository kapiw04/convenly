package integral

import (
	"database/sql"
	"testing"

	"github.com/kapiw04/convenly/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestUserRepo_SaveAndGet(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
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
	})
}

func TestUserRepo_RegisterGeneratesUUID(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := srvc.Register(
			"TestUser",
			"testuser@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		u, err := srvc.GetByEmail("testuser@example.com")
		require.NoError(t, err)

		assert.NotEqual(t, [16]byte{}, u.UUID)
		assert.NotEmpty(t, u.UUID.String())

		uByID, err := srvc.GetByUUID(u.UUID.String())
		require.NoError(t, err)
		assert.Equal(t, u.UUID, uByID.UUID)
		assert.Equal(t, "testuser@example.com", string(uByID.Email))
	})
}

func TestUserRepo_UserIDHasDefaultValue(t *testing.T) {
	sqlDb := setupDb(t)
	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		var columnDefault sql.NullString
		err := sqlDb.QueryRow(`
			SELECT column_default 
			FROM information_schema.columns 
			WHERE table_name = 'users' AND column_name = 'user_id'
		`).Scan(&columnDefault)
		require.NoError(t, err)
		require.True(t, columnDefault.Valid, "user_id column should have a DEFAULT constraint")
		assert.Contains(t, columnDefault.String, "gen_random_uuid", "user_id DEFAULT should use gen_random_uuid()")
	})
}

func TestUserRepo_RegisterDuplicateEmail(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {

		err := srvc.Register(
			"Alice",
			"alice@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		err = srvc.Register(
			"Bobby",
			"alice@example.com",
			"AnotherSecret456!",
		)
		assert.Error(t, err)
		assert.Equal(t, user.ErrUserExists, err)
	})
}

func TestUserRepo_RegisterCaseInsensitiveDuplicateEmail(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {

		err := srvc.Register(
			"Alice",
			"alice@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		err = srvc.Register(
			"Bobby",
			"ALICE@EXAMPLE.COM",
			"AnotherSecret456!",
		)
		assert.Error(t, err)
		assert.Equal(t, user.ErrUserExists, err)

		err = srvc.Register(
			"Charlie",
			"Alice@Example.Com",
			"YetAnother789!",
		)
		assert.Error(t, err)
		assert.Equal(t, user.ErrUserExists, err)
	})
}

func TestUserRepo_RegisterEmailWithWhitespace(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {

		err := srvc.Register(
			"Alice",
			"  alice@example.com  ",
			"Secret123!",
		)
		require.NoError(t, err)

		u, err := srvc.GetByEmail("alice@example.com")
		assert.NoError(t, err)
		assert.Equal(t, string(u.Email), "alice@example.com")
		assert.Equal(t, u.Name, "Alice")

		err = srvc.Register(
			"Bobby",
			" alice@example.com ",
			"AnotherSecret456!",
		)
		assert.Error(t, err)
		assert.Equal(t, user.ErrUserExists, err)
	})
}

func TestUserRepo_RegisterUsernameTooShort(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := srvc.Register(
			"Bob",
			"bob@example.com",
			"Secret123!",
		)
		assert.Error(t, err)
		assert.Equal(t, user.ErrUsernameTooShort, err)
	})
}

func TestUserRepo_RegisterUsernameFourChars(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := srvc.Register(
			"John",
			"john@example.com",
			"Secret123!",
		)
		assert.Error(t, err)
		assert.Equal(t, user.ErrUsernameTooShort, err)
	})
}

func TestUserRepo_RegisterUsernameExactlyFiveChars(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := srvc.Register(
			"Bobby",
			"bobby@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		u, err := srvc.GetByEmail("bobby@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "Bobby", u.Name)
	})
}

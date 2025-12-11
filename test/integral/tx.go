package integral

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func WithTx(t *testing.T, db *sql.DB, fn func(t *testing.T, tx *sql.Tx)) {
	t.Helper()
	cleanDatabase(t, db)
	fn(t, nil)
	cleanDatabase(t, db)
}

func cleanDatabase(t *testing.T, db *sql.DB) {
	t.Helper()

	queries := []string{
		"DELETE FROM event_tag",
		"DELETE FROM attendance",
		"DELETE FROM events",
		"DELETE FROM sessions",
		"DELETE FROM users",
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		require.NoError(t, err, "Failed to execute cleanup query: %s", query)
	}
}

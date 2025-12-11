package integral

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	pgConn, err = StartPostgres()
	if err != nil {
		log.Fatalf("failed to start postgres: %v", err)
	}

	sqlDB, err = sql.Open("postgres", pgConn.DSN)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err := applyMigrations(sqlDB); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	code := m.Run()

	_ = pgConn.Container.Terminate(ctx)
	os.Exit(code)
}

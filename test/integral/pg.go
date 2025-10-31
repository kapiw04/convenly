package integral

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"

	"github.com/testcontainers/testcontainers-go/wait"
)

type PgConnection struct {
	Port      string
	IP        string
	DB        string
	User      string
	Password  string
	Container testcontainers.Container
	DSN       string
}

func StartPostgres() (*PgConnection, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image: "postgres",
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_DB":       "db"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort("5432/tcp"),
			wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
				return fmt.Sprintf("host=%s port=%s dbname=db user=user password=pass sslmode=disable", host, port.Port())
			}),
		),
		Cmd: []string{"postgres", "-c", "max_connections=500"},
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}
	ip, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s port=%s dbname=db user=user password=pass sslmode=disable connect_timeout=10", ip, port.Port())

	return &PgConnection{
		Port:      port.Port(),
		IP:        ip,
		DB:        "db",
		User:      "user",
		Password:  "pass",
		DSN:       dsn,
		Container: pgContainer,
	}, nil
}

func ApplyMigrations(t *testing.T, db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	require.NoError(t, err)
	m, err := migrate.NewWithDatabaseInstance(
		"file://../../internal/infra/db/migrations",
		"postgres", driver,
	)
	require.NoError(t, err)

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}
}

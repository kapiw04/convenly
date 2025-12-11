package integral

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const reusablePgName = "convenly_pg_dev"

type PgConnection struct {
	Port      string
	IP        string
	DB        string
	User      string
	Password  string
	Container testcontainers.Container
	DSN       string
}

var (
	pgConn *PgConnection
	sqlDB  *sql.DB
)

func StartPostgres() (*PgConnection, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Name:  reusablePgName,
		Image: "postgres:16-alpine",
		Env: map[string]string{
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_DB":       "db",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForSQL("5432/tcp", "postgres", func(host string, port nat.Port) string {
			return fmt.Sprintf(
				"host=%s port=%s dbname=db user=user password=pass sslmode=disable",
				host,
				port.Port(),
			)
		}),
		Cmd: []string{"postgres", "-c", "max_connections=500"},
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Reuse:            true,
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

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=db user=user password=pass sslmode=disable connect_timeout=10",
		ip,
		port.Port(),
	)

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

func applyMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../internal/infra/db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

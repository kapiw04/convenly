package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/kapiw04/convenly/internal/app"
	"github.com/kapiw04/convenly/internal/infra/db"
	logger "github.com/kapiw04/convenly/internal/infra/log"
	"github.com/kapiw04/convenly/internal/infra/security"
	"github.com/kapiw04/convenly/internal/infra/webapi"
	_ "github.com/lib/pq"
)

func main() {
	logger.InitializeLogger("./logs")

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	if user == "" || password == "" || dbName == "" {
		slog.Error("Environment variables POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB must be set")
		return
	}
	connStr := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", user, password, dbName)
	postgresDb, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("Error connecting to the database", "err", err)
		return
	}
	defer postgresDb.Close()
	slog.Info("Successfully connected to the database")
	hasher := &security.BcryptHasher{}
	userRepo := db.NewPostgresUserRepo(postgresDb)
	sessionRepo := db.NewPostgresSessionRepo(postgresDb, userRepo)
	userService := app.NewUserService(userRepo, sessionRepo, hasher)
	tagsRepo := db.NewPostgresTagRepo(postgresDb)
	eventRepo := db.NewPostgresEventRepo(postgresDb, tagsRepo)
	eventService := app.NewEventService(eventRepo)

	router := webapi.NewRouter(userService, eventService)
	server := webapi.NewServer(":8080", router.Handler)
	webapi.Start(server)
	defer webapi.Stop(context.Background(), server)
}

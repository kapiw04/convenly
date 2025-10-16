package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/kapiw04/convenly/internal/app"
	"github.com/kapiw04/convenly/internal/infra/db"
	"github.com/kapiw04/convenly/internal/infra/http"
	logger "github.com/kapiw04/convenly/internal/infra/log"
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
	userRepo := db.NewPostgresUserRepo(postgresDb)
	userService := app.NewUserService(userRepo)

	router := http.NewRouter(userService)
	server := http.NewServer(":8080", router.Handler)
	http.Start(server)
	defer http.Stop(context.Background(), server)
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	if user == "" || password == "" || dbName == "" {
		log.Fatal("Environment variables POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB must be set")
	}

	connStr := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", user, password, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Starting database cleanup...")

	tables := []string{
		"event_tag",
		"attendance",
		"events",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			log.Printf("Warning: Could not truncate table %s: %v", table, err)
		} else {
			fmt.Printf("Truncated %s\n", table)
		}
	}

	fmt.Println("Database cleanup completed")
}

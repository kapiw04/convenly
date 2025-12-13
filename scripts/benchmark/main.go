package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	numEvents := flag.Int("events", 1000, "Number of events to generate")
	flag.Parse()

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Printf("Generating %d events...\n", *numEvents)
	startTime := time.Now()

	userIDs, err := getExistingUsers(ctx, db)
	if err != nil {
		log.Fatalf("Error fetching users: %v", err)
	}

	if len(userIDs) == 0 {
		log.Fatal("No users found in database. Create users first.")
	}

	tagIDs, err := getExistingTags(ctx, db)
	if err != nil {
		log.Fatalf("Error fetching tags: %v", err)
	}

	if len(tagIDs) == 0 {
		log.Fatal("No tags found in database. Create tags first.")
	}

	err = generateEvents(ctx, db, *numEvents, userIDs, tagIDs)
	if err != nil {
		log.Fatalf("Error generating events: %v", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Generated %d events in %.2f seconds\n", *numEvents, elapsed.Seconds())
}

func getExistingUsers(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT user_id FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	return userIDs, rows.Err()
}

func getExistingTags(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT tag_id FROM tags")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tagIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		tagIDs = append(tagIDs, id)
	}
	return tagIDs, rows.Err()
}

func generateEvents(ctx context.Context, db *sql.DB, count int, userIDs []string, tagIDs []string) error {
	locations := []struct {
		lat float64
		lon float64
	}{
		{52.2297, 21.0122},
		{50.0647, 19.9450},
		{51.7592, 19.4560},
		{54.3520, 18.6466},
		{51.1079, 17.0385},
		{52.0969, 23.7880},
		{49.5891, 19.0215},
		{52.1288, 17.0022},
	}

	for i := 0; i < count; i++ {
		eventID := uuid.New().String()
		name := fmt.Sprintf("Event %d", i)
		description := fmt.Sprintf("Description for event %d", i)

		daysOffset := rand.Intn(180)
		eventDate := time.Now().AddDate(0, 0, daysOffset)

		loc := locations[rand.Intn(len(locations))]
		fee := float64(rand.Intn(100))
		organizerID := userIDs[rand.Intn(len(userIDs))]

		_, err := db.ExecContext(ctx,
			"INSERT INTO events (event_id, name, description, date, latitude, longitude, fee, organizer_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			eventID, name, description, eventDate, loc.lat, loc.lon, fee, organizerID)
		if err != nil {
			return err
		}

		numTags := 1 + rand.Intn(3)
		for j := 0; j < numTags; j++ {
			tagID := tagIDs[rand.Intn(len(tagIDs))]
			db.ExecContext(ctx,
				"INSERT INTO event_tag (event_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				eventID, tagID)
		}

		if (i+1)%100 == 0 {
			fmt.Printf("  %d/%d events\n", i+1, count)
		}
	}

	return nil
}

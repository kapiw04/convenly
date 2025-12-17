package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TrafficStats struct {
	mu                sync.Mutex
	UsersCreated      int
	SessionsCreated   int
	EventsCreated     int
	Registrations     int
	Unregistrations   int
	EventQueries      int
	UserQueries       int
	AttendanceQueries int
	Errors            int
}

func (s *TrafficStats) Print() {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println("\n=== Traffic Simulation Results ===")
	fmt.Printf("Users Created:        %d\n", s.UsersCreated)
	fmt.Printf("Sessions Created:     %d\n", s.SessionsCreated)
	fmt.Printf("Events Created:       %d\n", s.EventsCreated)
	fmt.Printf("Registrations:        %d\n", s.Registrations)
	fmt.Printf("Unregistrations:      %d\n", s.Unregistrations)
	fmt.Printf("Event Queries:        %d\n", s.EventQueries)
	fmt.Printf("User Queries:         %d\n", s.UserQueries)
	fmt.Printf("Attendance Queries:   %d\n", s.AttendanceQueries)
	fmt.Printf("Errors:               %d\n", s.Errors)
}

func main() {
	numEvents := flag.Int("events", 1000, "Number of events to generate")
	numUsers := flag.Int("users", 100, "Number of users to generate")
	numRegistrations := flag.Int("registrations", 5000, "Number of event registrations to simulate")
	concurrency := flag.Int("concurrency", 10, "Number of concurrent workers")
	simulateTraffic := flag.Bool("traffic", false, "Run full traffic simulation")
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

	db.SetMaxOpenConns(*concurrency + 5)
	db.SetMaxIdleConns(*concurrency)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	if *simulateTraffic {
		runTrafficSimulation(ctx, db, *numUsers, *numEvents, *numRegistrations, *concurrency)
	} else {
		runLegacyEventGeneration(ctx, db, *numEvents)
	}
}

func runTrafficSimulation(ctx context.Context, db *sql.DB, numUsers, numEvents, numRegistrations, concurrency int) {
	stats := &TrafficStats{}
	startTime := time.Now()

	fmt.Println("=== Starting Traffic Simulation ===")

	// Get existing tags
	tagIDs, err := getExistingTags(ctx, db)
	if err != nil {
		log.Fatalf("Error fetching tags: %v", err)
	}
	if len(tagIDs) == 0 {
		log.Fatal("No tags found in database. Create tags first.")
	}

	// Phase 1: Create users
	fmt.Printf("\nPhase 1: Creating %d users...\n", numUsers)
	userIDs, err := createUsers(ctx, db, numUsers, stats)
	if err != nil {
		log.Fatalf("Error creating users: %v", err)
	}

	// Phase 2: Create sessions for users
	fmt.Printf("\nPhase 2: Creating sessions for users...\n")
	err = createSessions(ctx, db, userIDs, stats)
	if err != nil {
		log.Fatalf("Error creating sessions: %v", err)
	}

	// Phase 3: Create events
	fmt.Printf("\nPhase 3: Creating %d events...\n", numEvents)
	eventIDs, err := createEvents(ctx, db, numEvents, userIDs, tagIDs, stats)
	if err != nil {
		log.Fatalf("Error creating events: %v", err)
	}

	// Phase 4: Simulate registrations with concurrent workers
	fmt.Printf("\nPhase 4: Simulating %d registrations with %d workers...\n", numRegistrations, concurrency)
	simulateRegistrations(ctx, db, userIDs, eventIDs, numRegistrations, concurrency, stats)

	// Phase 5: Simulate read traffic
	fmt.Printf("\nPhase 5: Simulating read queries...\n")
	simulateReadTraffic(ctx, db, userIDs, eventIDs, concurrency, stats)

	// Phase 6: Simulate some unregistrations
	numUnregistrations := numRegistrations / 10
	fmt.Printf("\nPhase 6: Simulating %d unregistrations...\n", numUnregistrations)
	simulateUnregistrations(ctx, db, numUnregistrations, stats)

	elapsed := time.Since(startTime)
	stats.Print()
	fmt.Printf("\nTotal time: %.2f seconds\n", elapsed.Seconds())
}

func runLegacyEventGeneration(ctx context.Context, db *sql.DB, numEvents int) {
	fmt.Printf("Generating %d events...\n", numEvents)
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

	err = generateEvents(ctx, db, numEvents, userIDs, tagIDs)
	if err != nil {
		log.Fatalf("Error generating events: %v", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Generated %d events in %.2f seconds\n", numEvents, elapsed.Seconds())
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

// ============ Traffic Simulation Functions ============

var firstNames = []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank", "Grace", "Henry", "Ivy", "Jack",
	"Kate", "Leo", "Mia", "Noah", "Olivia", "Paul", "Quinn", "Rose", "Sam", "Tina",
	"Uma", "Victor", "Wendy", "Xavier", "Yuki", "Zack", "Anna", "Ben", "Clara", "David"}

var lastNames = []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
	"Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin"}

var eventPrefixes = []string{"Tech", "Music", "Art", "Food", "Sports", "Business", "Science", "Health", "Education", "Travel"}
var eventTypes = []string{"Conference", "Workshop", "Meetup", "Festival", "Seminar", "Exhibition", "Tournament", "Networking", "Summit", "Gala"}

func createUsers(ctx context.Context, db *sql.DB, count int, stats *TrafficStats) ([]string, error) {
	userIDs := make([]string, 0, count)

	for i := 0; i < count; i++ {
		userID := uuid.New().String()
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		name := fmt.Sprintf("%s %s", firstName, lastName)
		email := fmt.Sprintf("%s.%s.%d@benchmark.test", firstName, lastName, i)
		// Using a pre-hashed password (bcrypt hash of "BenchmarkPass123!")
		passwordHash := "$2a$10$N9qo8uLOickgx2ZMRZoMy.MqPXfNdKWOjKMKqelAADm3GZ/A/AAQK"

		_, err := db.ExecContext(ctx,
			"INSERT INTO users (user_id, email, password_hash, name) VALUES ($1, $2, $3, $4)",
			userID, email, passwordHash, name)
		if err != nil {
			stats.mu.Lock()
			stats.Errors++
			stats.mu.Unlock()
			continue
		}

		userIDs = append(userIDs, userID)
		stats.mu.Lock()
		stats.UsersCreated++
		stats.mu.Unlock()

		if (i+1)%50 == 0 {
			fmt.Printf("  Created %d/%d users\n", i+1, count)
		}
	}

	return userIDs, nil
}

func createSessions(ctx context.Context, db *sql.DB, userIDs []string, stats *TrafficStats) error {
	for _, userID := range userIDs {
		sessionID := uuid.New().String()

		_, err := db.ExecContext(ctx,
			"INSERT INTO sessions (user_id, session_id) VALUES ($1, $2)",
			userID, sessionID)
		if err != nil {
			stats.mu.Lock()
			stats.Errors++
			stats.mu.Unlock()
			continue
		}

		stats.mu.Lock()
		stats.SessionsCreated++
		stats.mu.Unlock()
	}

	fmt.Printf("  Created %d sessions\n", stats.SessionsCreated)
	return nil
}

func createEvents(ctx context.Context, db *sql.DB, count int, userIDs []string, tagIDs []string, stats *TrafficStats) ([]string, error) {
	locations := []struct {
		lat  float64
		lon  float64
		city string
	}{
		{52.2297, 21.0122, "Warsaw"},
		{50.0647, 19.9450, "Krakow"},
		{51.7592, 19.4560, "Lodz"},
		{54.3520, 18.6466, "Gdansk"},
		{51.1079, 17.0385, "Wroclaw"},
		{52.0969, 23.7880, "Bialystok"},
		{49.5891, 19.0215, "Zakopane"},
		{52.1288, 17.0022, "Leszno"},
		{50.2649, 19.0238, "Katowice"},
		{53.1235, 18.0084, "Bydgoszcz"},
	}

	eventIDs := make([]string, 0, count)

	for i := 0; i < count; i++ {
		eventID := uuid.New().String()
		prefix := eventPrefixes[rand.Intn(len(eventPrefixes))]
		eventType := eventTypes[rand.Intn(len(eventTypes))]
		loc := locations[rand.Intn(len(locations))]
		name := fmt.Sprintf("%s %s %s %d", loc.city, prefix, eventType, i)
		description := fmt.Sprintf("Join us for an exciting %s %s in %s! This event brings together enthusiasts and professionals for an unforgettable experience.", prefix, eventType, loc.city)

		daysOffset := rand.Intn(365) - 30 // Events from 30 days ago to 335 days in future
		eventDate := time.Now().AddDate(0, 0, daysOffset)

		fee := float64(rand.Intn(200))
		organizerID := userIDs[rand.Intn(len(userIDs))]

		_, err := db.ExecContext(ctx,
			"INSERT INTO events (event_id, name, description, date, latitude, longitude, fee, organizer_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			eventID, name, description, eventDate, loc.lat, loc.lon, fee, organizerID)
		if err != nil {
			stats.mu.Lock()
			stats.Errors++
			stats.mu.Unlock()
			continue
		}

		// Add 1-4 tags per event
		numTags := 1 + rand.Intn(4)
		usedTags := make(map[string]bool)
		for j := 0; j < numTags; j++ {
			tagID := tagIDs[rand.Intn(len(tagIDs))]
			if usedTags[tagID] {
				continue
			}
			usedTags[tagID] = true
			db.ExecContext(ctx,
				"INSERT INTO event_tag (event_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				eventID, tagID)
		}

		eventIDs = append(eventIDs, eventID)
		stats.mu.Lock()
		stats.EventsCreated++
		stats.mu.Unlock()

		if (i+1)%100 == 0 {
			fmt.Printf("  Created %d/%d events\n", i+1, count)
		}
	}

	return eventIDs, nil
}

func simulateRegistrations(ctx context.Context, db *sql.DB, userIDs, eventIDs []string, count, concurrency int, stats *TrafficStats) {
	jobs := make(chan int, count)
	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				userID := userIDs[rand.Intn(len(userIDs))]
				eventID := eventIDs[rand.Intn(len(eventIDs))]

				_, err := db.ExecContext(ctx,
					"INSERT INTO attendance (user_id, event_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
					userID, eventID)
				if err != nil {
					stats.mu.Lock()
					stats.Errors++
					stats.mu.Unlock()
					continue
				}

				stats.mu.Lock()
				stats.Registrations++
				stats.mu.Unlock()
			}
		}()
	}

	// Send jobs
	for i := 0; i < count; i++ {
		jobs <- i
		if (i+1)%500 == 0 {
			stats.mu.Lock()
			regs := stats.Registrations
			stats.mu.Unlock()
			fmt.Printf("  Processed %d/%d registration attempts (%d successful)\n", i+1, count, regs)
		}
	}
	close(jobs)
	wg.Wait()

	fmt.Printf("  Completed %d registrations\n", stats.Registrations)
}

func simulateReadTraffic(ctx context.Context, db *sql.DB, userIDs, eventIDs []string, concurrency int, stats *TrafficStats) {
	var wg sync.WaitGroup

	// Simulate event queries
	numEventQueries := len(eventIDs) / 2
	if numEventQueries > 500 {
		numEventQueries = 500
	}

	jobs := make(chan string, numEventQueries)
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for eventID := range jobs {
				// Query event details
				var name, description string
				var date time.Time
				err := db.QueryRowContext(ctx,
					"SELECT name, description, date FROM events WHERE event_id = $1",
					eventID).Scan(&name, &description, &date)
				if err == nil {
					stats.mu.Lock()
					stats.EventQueries++
					stats.mu.Unlock()
				}

				// Query attendee count
				var count int
				err = db.QueryRowContext(ctx,
					"SELECT COUNT(*) FROM attendance WHERE event_id = $1",
					eventID).Scan(&count)
				if err == nil {
					stats.mu.Lock()
					stats.AttendanceQueries++
					stats.mu.Unlock()
				}
			}
		}()
	}

	for i := 0; i < numEventQueries; i++ {
		jobs <- eventIDs[rand.Intn(len(eventIDs))]
	}
	close(jobs)
	wg.Wait()

	// Simulate user queries
	numUserQueries := len(userIDs) / 2
	if numUserQueries > 200 {
		numUserQueries = 200
	}

	userJobs := make(chan string, numUserQueries)
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for userID := range userJobs {
				// Query user's registered events
				rows, err := db.QueryContext(ctx,
					"SELECT e.event_id, e.name FROM events e INNER JOIN attendance a ON e.event_id = a.event_id WHERE a.user_id = $1",
					userID)
				if err == nil {
					for rows.Next() {
						var eid, ename string
						rows.Scan(&eid, &ename)
					}
					rows.Close()
					stats.mu.Lock()
					stats.UserQueries++
					stats.mu.Unlock()
				}
			}
		}()
	}

	for i := 0; i < numUserQueries; i++ {
		userJobs <- userIDs[rand.Intn(len(userIDs))]
	}
	close(userJobs)
	wg.Wait()

	fmt.Printf("  Completed %d event queries, %d attendance queries, %d user queries\n",
		stats.EventQueries, stats.AttendanceQueries, stats.UserQueries)
}

func simulateUnregistrations(ctx context.Context, db *sql.DB, count int, stats *TrafficStats) {
	// Get some existing registrations to unregister
	rows, err := db.QueryContext(ctx, "SELECT user_id, event_id FROM attendance ORDER BY RANDOM() LIMIT $1", count)
	if err != nil {
		log.Printf("Error fetching registrations for unregistration: %v", err)
		return
	}
	defer rows.Close()

	type registration struct {
		userID  string
		eventID string
	}

	registrations := make([]registration, 0, count)
	for rows.Next() {
		var r registration
		if err := rows.Scan(&r.userID, &r.eventID); err != nil {
			continue
		}
		registrations = append(registrations, r)
	}

	for _, r := range registrations {
		_, err := db.ExecContext(ctx,
			"DELETE FROM attendance WHERE user_id = $1 AND event_id = $2",
			r.userID, r.eventID)
		if err != nil {
			stats.mu.Lock()
			stats.Errors++
			stats.mu.Unlock()
			continue
		}

		stats.mu.Lock()
		stats.Unregistrations++
		stats.mu.Unlock()
	}

	fmt.Printf("  Completed %d unregistrations\n", stats.Unregistrations)
}

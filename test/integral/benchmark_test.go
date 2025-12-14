package integral

import (
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kapiw04/convenly/internal/domain/event"
	"github.com/kapiw04/convenly/internal/domain/user"
	"github.com/kapiw04/convenly/internal/infra/db"
)

func seedBenchmarkData(b *testing.B, eventRepo event.EventRepo, userID string, count int) []string {
	b.Helper()

	eventIDs := make([]string, count)
	tags := []string{"music", "sports", "technology", "art", "food"}

	for i := 0; i < count; i++ {
		e := &event.Event{
			EventID:     uuid.New().String(),
			Name:        fmt.Sprintf("Benchmark Event %d", i),
			Description: fmt.Sprintf("Description for benchmark event %d", i),
			Date:        time.Now().AddDate(0, 0, i%180),
			Latitude:    52.2297 + float64(i%10)*0.01,
			Longitude:   21.0122 + float64(i%10)*0.01,
			Fee:         float32(i % 100),
			OrganizerID: userID,
			Tags:        []string{tags[i%len(tags)]},
		}
		eventRepo.Save(e)
		eventIDs[i] = e.EventID
	}

	return eventIDs
}

func setupBenchmarkRepos(b *testing.B) (event.EventRepo, user.UserRepo, event.TagRepo) {
	b.Helper()

	tagRepo := db.NewPostgresTagRepo(sqlDB)
	eventRepo := db.NewPostgresEventRepo(sqlDB, tagRepo)
	userRepo := db.NewPostgresUserRepo(sqlDB)

	return eventRepo, userRepo, tagRepo
}

func createBenchmarkUser(b *testing.B, userRepo user.UserRepo) string {
	b.Helper()
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	id := uuid.New()
	u := &user.User{
		UUID:         id,
		Name:         fmt.Sprintf("bench_%s", id.String()[:8]),
		Email:        user.Email(fmt.Sprintf("bench_%s@test.com", id.String()[:8])),
		PasswordHash: "hashedpass",
		Role:         user.ATTENDEE,
	}
	userRepo.Save(u)
	return id.String()
}

func BenchmarkFindAllEvents(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	eventRepo, userRepo, _ := setupBenchmarkRepos(b)
	userID := createBenchmarkUser(b, userRepo)

	for _, size := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("events_%d", size), func(b *testing.B) {
			seedBenchmarkData(b, eventRepo, userID, size)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				eventRepo.FindAll()
			}
		})
	}
}

func BenchmarkFindEventByID(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))

	eventRepo, userRepo, _ := setupBenchmarkRepos(b)
	userID := createBenchmarkUser(b, userRepo)
	eventIDs := seedBenchmarkData(b, eventRepo, userID, 1000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eventRepo.FindByID(eventIDs[i%len(eventIDs)])
	}
}

func BenchmarkFindAllWithTagFilter(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	eventRepo, userRepo, _ := setupBenchmarkRepos(b)
	userID := createBenchmarkUser(b, userRepo)
	seedBenchmarkData(b, eventRepo, userID, 1000)

	filter := &event.EventFilter{
		Tags: []string{"music", "sports"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eventRepo.FindAllWithFilters(filter)
	}
}

func BenchmarkFindAllWithDateFilter(b *testing.B) {
	eventRepo, userRepo, _ := setupBenchmarkRepos(b)
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	userID := createBenchmarkUser(b, userRepo)
	seedBenchmarkData(b, eventRepo, userID, 1000)

	from := time.Now()
	to := time.Now().AddDate(0, 1, 0)
	filter := &event.EventFilter{
		DateFrom: &from,
		DateTo:   &to,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eventRepo.FindAllWithFilters(filter)
	}
}

func BenchmarkFindAllWithCombinedFilters(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	eventRepo, userRepo, _ := setupBenchmarkRepos(b)
	userID := createBenchmarkUser(b, userRepo)
	seedBenchmarkData(b, eventRepo, userID, 1000)

	from := time.Now()
	to := time.Now().AddDate(0, 1, 0)
	minFee := float32(10)
	maxFee := float32(50)

	filter := &event.EventFilter{
		DateFrom: &from,
		DateTo:   &to,
		MinFee:   &minFee,
		MaxFee:   &maxFee,
		Tags:     []string{"music"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eventRepo.FindAllWithFilters(filter)
	}
}

func BenchmarkFindByOrganizer(b *testing.B) {
	eventRepo, userRepo, _ := setupBenchmarkRepos(b)
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	userID := createBenchmarkUser(b, userRepo)
	seedBenchmarkData(b, eventRepo, userID, 500)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eventRepo.FindByOrganizer(userID, nil)
	}
}

func BenchmarkRegisterAttendance(b *testing.B) {
	eventRepo, userRepo, _ := setupBenchmarkRepos(b)
	userID := createBenchmarkUser(b, userRepo)
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	eventIDs := seedBenchmarkData(b, eventRepo, userID, 100)

	attendees := make([]string, b.N)
	for i := range attendees {
		attendees[i] = createBenchmarkUser(b, userRepo)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eventRepo.RegisterAttendance(attendees[i], eventIDs[i%len(eventIDs)])
	}
}

func BenchmarkGetAttendees(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	eventRepo, userRepo, _ := setupBenchmarkRepos(b)
	organizerID := createBenchmarkUser(b, userRepo)
	eventIDs := seedBenchmarkData(b, eventRepo, organizerID, 10)

	for i := 0; i < 100; i++ {
		attendeeID := createBenchmarkUser(b, userRepo)
		eventRepo.RegisterAttendance(attendeeID, eventIDs[i%len(eventIDs)])
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eventRepo.GetAttendees(eventIDs[i%len(eventIDs)])
	}
}

func BenchmarkSaveEvent(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	eventRepo, userRepo, _ := setupBenchmarkRepos(b)
	userID := createBenchmarkUser(b, userRepo)

	events := make([]*event.Event, b.N)
	for i := range events {
		events[i] = &event.Event{
			EventID:     uuid.New().String(),
			Name:        fmt.Sprintf("New Event %d", i),
			Description: "Benchmark insert test",
			Date:        time.Now().AddDate(0, 0, i%30),
			Latitude:    52.2297,
			Longitude:   21.0122,
			Fee:         25.0,
			OrganizerID: userID,
			Tags:        []string{"technology"},
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		eventRepo.Save(events[i])
	}
}

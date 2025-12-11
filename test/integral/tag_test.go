package integral

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kapiw04/convenly/internal/domain/event"
	"github.com/kapiw04/convenly/internal/infra/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func setupTagRepo(t *testing.T, dbConn *sql.DB) *db.PostgresTagRepo {
	t.Helper()
	return db.NewPostgresTagRepo(dbConn)
}

func TestTagRepo_FindAll_ReturnsDefaultTags(t *testing.T) {
	sqlDb := setupDb(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		tagRepo := setupTagRepo(t, sqlDb)

		tags, err := tagRepo.FindAll()
		require.NoError(t, err)

		assert.Len(t, tags, len(event.DefaultTagNames))

		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			tagNames[i] = tag.Name
		}

		for _, defaultTag := range event.DefaultTagNames {
			assert.Contains(t, tagNames, defaultTag)
		}
	})
}

func TestTagRepo_FindByName_ExistingTag(t *testing.T) {
	sqlDb := setupDb(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		tagRepo := setupTagRepo(t, sqlDb)

		tag, err := tagRepo.FindByName("Music")
		require.NoError(t, err)

		require.NotNil(t, tag)
		assert.Equal(t, "Music", tag.Name)
		assert.Greater(t, tag.TagID, int64(0))
	})
}

func TestTagRepo_FindByName_NonExistingTag(t *testing.T) {
	sqlDb := setupDb(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		tagRepo := setupTagRepo(t, sqlDb)

		tag, err := tagRepo.FindByName("NonExistentTag")
		require.NoError(t, err)

		assert.Nil(t, tag)
	})
}

func TestTagRepo_CreateIfNotExists_NewTag(t *testing.T) {
	sqlDb := setupDb(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		tagRepo := setupTagRepo(t, sqlDb)

		tag, err := tagRepo.CreateIfNotExists("CustomTag")
		require.NoError(t, err)

		require.NotNil(t, tag)
		assert.Equal(t, "CustomTag", tag.Name)
		assert.Greater(t, tag.TagID, int64(0))

		foundTag, err := tagRepo.FindByName("CustomTag")
		require.NoError(t, err)
		require.NotNil(t, foundTag)
		assert.Equal(t, tag.TagID, foundTag.TagID)
	})
}

func TestTagRepo_CreateIfNotExists_ExistingTag(t *testing.T) {
	sqlDb := setupDb(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		tagRepo := setupTagRepo(t, sqlDb)

		existingTag, err := tagRepo.FindByName("Music")
		require.NoError(t, err)
		require.NotNil(t, existingTag)

		tag, err := tagRepo.CreateIfNotExists("Music")
		require.NoError(t, err)

		require.NotNil(t, tag)
		assert.Equal(t, existingTag.TagID, tag.TagID)
		assert.Equal(t, existingTag.Name, tag.Name)
	})
}

func TestTagRepo_SeedDefaults_Idempotent(t *testing.T) {
	sqlDb := setupDb(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		tagRepo := setupTagRepo(t, sqlDb)

		initialTags, err := tagRepo.FindAll()
		require.NoError(t, err)
		initialCount := len(initialTags)

		err = tagRepo.SeedDefaults()
		require.NoError(t, err)

		tags, err := tagRepo.FindAll()
		require.NoError(t, err)
		assert.Len(t, tags, initialCount)
	})
}

func TestTagRepo_FindAll_ReturnedTagsHaveValidIDs(t *testing.T) {
	sqlDb := setupDb(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		tagRepo := setupTagRepo(t, sqlDb)

		tags, err := tagRepo.FindAll()
		require.NoError(t, err)

		for _, tag := range tags {
			assert.Greater(t, tag.TagID, int64(0))
			assert.NotEmpty(t, tag.Name)
		}
	})
}

func TestTagRepo_CreateIfNotExists_MultipleTimes(t *testing.T) {
	sqlDb := setupDb(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		tagRepo := setupTagRepo(t, sqlDb)

		tag1, err := tagRepo.CreateIfNotExists("UniqueTag")
		require.NoError(t, err)

		tag2, err := tagRepo.CreateIfNotExists("UniqueTag")
		require.NoError(t, err)

		tag3, err := tagRepo.CreateIfNotExists("UniqueTag")
		require.NoError(t, err)

		assert.Equal(t, tag1.TagID, tag2.TagID)
		assert.Equal(t, tag2.TagID, tag3.TagID)
	})
}

func TestEventRepo_FindAllByTags_SingleTag(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, _ := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		user, err := userSrvc.GetByEmail("host@example.com")
		require.NoError(t, err)

		musicEvent := createTestEvent(t, "Music Festival", user.UUID.String(), []string{"Music"})
		err = eventSrvc.CreateEvent(musicEvent)
		require.NoError(t, err)

		sportsEvent := createTestEvent(t, "Sports Day", user.UUID.String(), []string{"Sports"})
		err = eventSrvc.CreateEvent(sportsEvent)
		require.NoError(t, err)

		events, err := eventSrvc.GetEventByTag([]string{"Music"})
		require.NoError(t, err)

		assert.Len(t, events, 1)
		assert.Equal(t, "Music Festival", events[0].Name)
	})
}

func TestEventRepo_FindAllByTags_MultipleTags(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, _ := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		user, err := userSrvc.GetByEmail("host@example.com")
		require.NoError(t, err)

		musicEvent := createTestEvent(t, "Music Festival", user.UUID.String(), []string{"Music"})
		err = eventSrvc.CreateEvent(musicEvent)
		require.NoError(t, err)

		sportsEvent := createTestEvent(t, "Sports Day", user.UUID.String(), []string{"Sports"})
		err = eventSrvc.CreateEvent(sportsEvent)
		require.NoError(t, err)

		techEvent := createTestEvent(t, "Tech Conference", user.UUID.String(), []string{"Tech"})
		err = eventSrvc.CreateEvent(techEvent)
		require.NoError(t, err)

		events, err := eventSrvc.GetEventByTag([]string{"Music", "Sports"})
		require.NoError(t, err)

		assert.Len(t, events, 2)
		eventNames := []string{events[0].Name, events[1].Name}
		assert.Contains(t, eventNames, "Music Festival")
		assert.Contains(t, eventNames, "Sports Day")
	})
}

func TestEventRepo_FindAllByTags_NoMatchingEvents(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, _ := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		user, err := userSrvc.GetByEmail("host@example.com")
		require.NoError(t, err)

		musicEvent := createTestEvent(t, "Music Festival", user.UUID.String(), []string{"Music"})
		err = eventSrvc.CreateEvent(musicEvent)
		require.NoError(t, err)

		events, err := eventSrvc.GetEventByTag([]string{"Gaming"})
		require.NoError(t, err)

		assert.Empty(t, events)
	})
}

func TestEventRepo_FindAllByTags_EventWithMultipleTags(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, _ := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		user, err := userSrvc.GetByEmail("host@example.com")
		require.NoError(t, err)

		multiTagEvent := createTestEvent(t, "Tech Music Party", user.UUID.String(), []string{"Tech", "Music", "Party"})
		err = eventSrvc.CreateEvent(multiTagEvent)
		require.NoError(t, err)

		singleTagEvent := createTestEvent(t, "Pure Music", user.UUID.String(), []string{"Music"})
		err = eventSrvc.CreateEvent(singleTagEvent)
		require.NoError(t, err)

		events, err := eventSrvc.GetEventByTag([]string{"Tech"})
		require.NoError(t, err)

		assert.Len(t, events, 1)
		assert.Equal(t, "Tech Music Party", events[0].Name)
	})
}

func TestEventRepo_FindAllByTags_EmptyTagList(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, _ := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		user, err := userSrvc.GetByEmail("host@example.com")
		require.NoError(t, err)

		musicEvent := createTestEvent(t, "Music Festival", user.UUID.String(), []string{"Music"})
		err = eventSrvc.CreateEvent(musicEvent)
		require.NoError(t, err)

		events, err := eventSrvc.GetEventByTag([]string{})
		require.NoError(t, err)

		assert.Empty(t, events)
	})
}

func TestEventRepo_FindAllByTags_NonExistentTag(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, _ := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		user, err := userSrvc.GetByEmail("host@example.com")
		require.NoError(t, err)

		musicEvent := createTestEvent(t, "Music Festival", user.UUID.String(), []string{"Music"})
		err = eventSrvc.CreateEvent(musicEvent)
		require.NoError(t, err)

		events, err := eventSrvc.GetEventByTag([]string{"NonExistentTag123"})
		require.NoError(t, err)

		assert.Empty(t, events)
	})
}

func createTestEvent(t *testing.T, name, organizerID string, tags []string) *event.Event {
	t.Helper()
	return &event.Event{
		EventID:     uuid.New().String(),
		Name:        name,
		Description: "Test description",
		Date:        time.Now().Add(24 * time.Hour),
		Latitude:    42.0,
		Longitude:   21.37,
		Fee:         10.0,
		OrganizerID: organizerID,
		Tags:        tags,
	}
}

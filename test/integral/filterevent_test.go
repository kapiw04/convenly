package integral

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapiw04/convenly/internal/domain/event"
	"github.com/kapiw04/convenly/internal/infra/webapi"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestFilterEvents_ByDateRange(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		sessionID, err := userSrvc.Login("host@example.com", "Secret123!")
		require.NoError(t, err)

		createEventWithDetails(t, router, sessionID, "January Event", "2025-01-15T10:00:00Z", 10.0, []string{})
		createEventWithDetails(t, router, sessionID, "February Event", "2025-02-15T10:00:00Z", 20.0, []string{})
		createEventWithDetails(t, router, sessionID, "March Event", "2025-03-15T10:00:00Z", 30.0, []string{})

		req := httptest.NewRequest(http.MethodGet, "/api/events?date_from=2025-02-01", nil)
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var events []*event.Event
		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 2)
		require.Equal(t, "February Event", events[0].Name)
		require.Equal(t, "March Event", events[1].Name)

		req = httptest.NewRequest(http.MethodGet, "/api/events?date_to=2025-02-28", nil)
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 2)
		require.Equal(t, "January Event", events[0].Name)
		require.Equal(t, "February Event", events[1].Name)

		req = httptest.NewRequest(http.MethodGet, "/api/events?date_from=2025-02-01&date_to=2025-02-28", nil)
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "February Event", events[0].Name)

		allEvents, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, allEvents, 3)
	})
}

func TestFilterEvents_ByFeeRange(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		sessionID, err := userSrvc.Login("host@example.com", "Secret123!")
		require.NoError(t, err)

		createEventWithDetails(t, router, sessionID, "Free Event", "2025-01-15T10:00:00Z", 0.0, []string{})
		createEventWithDetails(t, router, sessionID, "Cheap Event", "2025-02-15T10:00:00Z", 10.0, []string{})
		createEventWithDetails(t, router, sessionID, "Expensive Event", "2025-03-15T10:00:00Z", 100.0, []string{})

		req := httptest.NewRequest(http.MethodGet, "/api/events?max_fee=0", nil)
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var events []*event.Event
		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "Free Event", events[0].Name)

		req = httptest.NewRequest(http.MethodGet, "/api/events?min_fee=10", nil)
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 2)
		require.Equal(t, "Cheap Event", events[0].Name)
		require.Equal(t, "Expensive Event", events[1].Name)

		req = httptest.NewRequest(http.MethodGet, "/api/events?min_fee=5&max_fee=50", nil)
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "Cheap Event", events[0].Name)

		allEvents, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, allEvents, 3)
	})
}

func TestFilterEvents_ByTags(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		sessionID, err := userSrvc.Login("host@example.com", "Secret123!")
		require.NoError(t, err)

		createEventWithDetails(t, router, sessionID, "Music Festival", "2025-01-15T10:00:00Z", 50.0, []string{"Music"})
		createEventWithDetails(t, router, sessionID, "Tech Conference", "2025-02-15T10:00:00Z", 100.0, []string{"Tech"})
		createEventWithDetails(t, router, sessionID, "Music & Tech Event", "2025-03-15T10:00:00Z", 75.0, []string{"Music", "Tech"})

		req := httptest.NewRequest(http.MethodGet, "/api/events?tags=Music", nil)
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var events []*event.Event
		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 2)

		req = httptest.NewRequest(http.MethodGet, "/api/events?tags=Tech", nil)
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 2)

		req = httptest.NewRequest(http.MethodGet, "/api/events?tags=Music,Tech", nil)
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 3)

		allEvents, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, allEvents, 3)
	})
}

func TestFilterEvents_CombinedFilters(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		sessionID, err := userSrvc.Login("host@example.com", "Secret123!")
		require.NoError(t, err)
		createEventWithDetails(t, router, sessionID, "Cheap Music January", "2025-01-15T10:00:00Z", 10.0, []string{"Music"})
		createEventWithDetails(t, router, sessionID, "Expensive Music February", "2025-02-15T10:00:00Z", 100.0, []string{"Music"})
		createEventWithDetails(t, router, sessionID, "Cheap Tech March", "2025-03-15T10:00:00Z", 15.0, []string{"Tech"})
		createEventWithDetails(t, router, sessionID, "Expensive Tech April", "2025-04-15T10:00:00Z", 200.0, []string{"Tech"})

		req := httptest.NewRequest(http.MethodGet, "/api/events?tags=Music&max_fee=50&date_from=2025-01-01&date_to=2025-03-31", nil)
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var events []*event.Event
		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "Cheap Music January", events[0].Name)

		req = httptest.NewRequest(http.MethodGet, "/api/events?max_fee=20", nil)
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 2)

		req = httptest.NewRequest(http.MethodGet, "/api/events?tags=Tech&date_from=2025-02-01", nil)
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 2)
	})
}

func TestFilterEvents_NoResults(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		sessionID, err := userSrvc.Login("host@example.com", "Secret123!")
		require.NoError(t, err)

		createEventWithDetails(t, router, sessionID, "Test Event", "2025-06-15T10:00:00Z", 50.0, []string{"Music"})

		req := httptest.NewRequest(http.MethodGet, "/api/events?max_fee=10", nil)
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var events []*event.Event
		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 0)
	})
}

func TestFilterEvents_InvalidDateFormat(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		req := httptest.NewRequest(http.MethodGet, "/api/events?date_from=invalid-date", nil)
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestFilterEvents_InvalidFeeFormat(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		req := httptest.NewRequest(http.MethodGet, "/api/events?min_fee=not-a-number", nil)
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestFilterEvents_NoFilters_ReturnsAll(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "host@example.com", "Secret123!")
		sessionID, err := userSrvc.Login("host@example.com", "Secret123!")
		require.NoError(t, err)

		createEventWithDetails(t, router, sessionID, "Event 1", "2025-01-15T10:00:00Z", 10.0, []string{})
		createEventWithDetails(t, router, sessionID, "Event 2", "2025-02-15T10:00:00Z", 20.0, []string{})

		req := httptest.NewRequest(http.MethodGet, "/api/events", nil)
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var events []*event.Event
		err = json.Unmarshal(w.Body.Bytes(), &events)
		require.NoError(t, err)
		require.Len(t, events, 2)

		allEvents, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, allEvents, 2)
	})
}

func createEventWithDetails(t *testing.T, router *webapi.Router, sessionID, name, date string, fee float32, tags []string) {
	t.Helper()

	req := webapi.CreateEventRequest{
		Name:        name,
		Description: "Test event description",
		Latitude:    42.0,
		Longitude:   21.37,
		Fee:         fee,
		Date:        date,
		Tags:        tags,
	}
	body, err := json.Marshal(req)
	require.NoError(t, err)

	httpReq := newEventRequest(t, body, sessionID)
	w := httptest.NewRecorder()
	router.Handler.ServeHTTP(w, httpReq)
	require.Equal(t, http.StatusCreated, w.Code)
}

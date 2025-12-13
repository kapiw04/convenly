package integral

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapiw04/convenly/internal/app"
	"github.com/kapiw04/convenly/internal/domain/event"
	"github.com/kapiw04/convenly/internal/infra/webapi"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestEventDetail_Success(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")

		createTestEventViaAPI(t, router, sessionID, "Test Event", "2025-12-31T23:59:59Z", 25.0, []string{"Music"})

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
		eventID := events[0].EventID

		req := httptest.NewRequest(http.MethodGet, "/api/events/"+eventID, nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp eventDetailResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, "Test Event", resp.Name)
		require.Equal(t, float32(25.0), resp.Fee)
		require.Equal(t, 0, resp.AttendeesCount)
		require.False(t, resp.UserRegistered)
	})
}

func TestEventDetail_WithAttendees(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createTestEventViaAPI(t, router, hostSessionID, "Party Event", "2025-12-31T23:59:59Z", 10.0, []string{"Music"})

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		attendeeSessionID := RegisterAndLoginUser(t, userSrvc, "Attendee", "attendee@example.com", "Secret123!")
		registerReq := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		registerReq.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, registerReq)
		require.Equal(t, http.StatusOK, w.Code)

		req := httptest.NewRequest(http.MethodGet, "/api/events/"+eventID, nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp eventDetailResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Equal(t, 1, resp.AttendeesCount)
		require.True(t, resp.UserRegistered)
	})
}

func TestEventDetail_Unauthorized(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createTestEventViaAPI(t, router, sessionID, "Test Event", "2025-12-31T23:59:59Z", 25.0, []string{"Music"})

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		req := httptest.NewRequest(http.MethodGet, "/api/events/"+eventID, nil)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestEventDetail_NonExistent(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodGet, "/api/events/00000000-0000-0000-0000-000000000000", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestEventDetail_InvalidUUIDFormat(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodGet, "/api/events/not-a-valid-uuid", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

type eventDetailResponse struct {
	*event.Event
	AttendeesCount int  `json:"attendees_count"`
	UserRegistered bool `json:"user_registered"`
}

func registerHostAndLogin(t *testing.T, userSrvc *app.UserService, email, password string) string {
	t.Helper()
	registerAndPromoteHost(t, userSrvc, email, password)
	sessionID, err := userSrvc.Login(email, password)
	require.NoError(t, err)
	return sessionID
}

func createTestEventViaAPI(t *testing.T, router *webapi.Router, sessionID, name, date string, fee float32, tags []string) {
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

	httpReq := httptest.NewRequest(http.MethodPost, "/api/events/add", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
	w := httptest.NewRecorder()
	router.Handler.ServeHTTP(w, httpReq)
	require.Equal(t, http.StatusCreated, w.Code)
}

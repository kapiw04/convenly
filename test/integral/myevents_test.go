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

type MyEventsResponse struct {
	Hosting   []*event.Event `json:"hosting"`
	Attending []*event.Event `json:"attending"`
}

func TestMyEvents_Unauthorized(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		req := httptest.NewRequest(http.MethodGet, "/api/my-events", nil)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestMyEvents_EmptyEvents(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodGet, "/api/my-events", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp MyEventsResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Hosting, 0)
		require.Len(t, resp.Attending, 0)
	})
}

func TestMyEvents_HostingEvents(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForMyEvents(t, router, hostSessionID, "Test Event 1")
		createEventForMyEvents(t, router, hostSessionID, "Test Event 2")

		req := httptest.NewRequest(http.MethodGet, "/api/my-events", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: hostSessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp MyEventsResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Hosting, 2)
		require.Len(t, resp.Attending, 0)
	})
}

func TestMyEvents_AttendingEvents(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForMyEvents(t, router, hostSessionID, "Test Event 1")

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		attendeeSessionID := RegisterAndLoginUser(t, userSrvc, "Attendee", "attendee@example.com", "Secret123!")

		registerReq := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		registerReq.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, registerReq)
		require.Equal(t, http.StatusOK, w.Code)

		req := httptest.NewRequest(http.MethodGet, "/api/my-events", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w = httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp MyEventsResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Hosting, 0)
		require.Len(t, resp.Attending, 1)
		require.Equal(t, "Test Event 1", resp.Attending[0].Name)
	})
}

func TestMyEvents_HostingAndAttending(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		host1SessionID := registerHostAndLogin(t, userSrvc, "host1@example.com", "Secret123!")
		createEventForMyEvents(t, router, host1SessionID, "Host 1 Event")

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		host1EventID := events[0].EventID

		host2SessionID := registerHostAndLoginWithName(t, userSrvc, "Host 2", "host2@example.com", "Secret123!")
		createEventForMyEvents(t, router, host2SessionID, "Host 2 Event")

		registerReq := httptest.NewRequest(http.MethodPost, "/api/events/"+host1EventID+"/register", nil)
		registerReq.AddCookie(&http.Cookie{Name: "session-id", Value: host2SessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, registerReq)
		require.Equal(t, http.StatusOK, w.Code)

		req := httptest.NewRequest(http.MethodGet, "/api/my-events", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: host2SessionID})
		w = httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp MyEventsResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		require.Len(t, resp.Hosting, 1)
		require.Len(t, resp.Attending, 1)
		require.Equal(t, "Host 2 Event", resp.Hosting[0].Name)
		require.Equal(t, "Host 1 Event", resp.Attending[0].Name)
	})
}

func createEventForMyEvents(t *testing.T, router *webapi.Router, sessionID string, name string) {
	t.Helper()
	req := webapi.CreateEventRequest{
		Name:        name,
		Description: "Test event description",
		Latitude:    42.0,
		Longitude:   21.37,
		Fee:         10.0,
		Date:        "2025-12-31T23:59:59Z",
		Tags:        []string{"Music"},
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

func registerHostAndLoginWithName(t *testing.T, userSrvc *app.UserService, name, email, password string) string {
	t.Helper()
	err := userSrvc.Register(name, email, password)
	require.NoError(t, err)
	user, err := userSrvc.GetByEmail(email)
	require.NoError(t, err)
	err = userSrvc.PromoteToHost(user.UUID.String())
	require.NoError(t, err)
	sessionID, err := userSrvc.Login(email, password)
	require.NoError(t, err)
	return sessionID
}

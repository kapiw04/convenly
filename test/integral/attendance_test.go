package integral

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapiw04/convenly/internal/infra/webapi"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestAttendance_RegisterSuccess(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForAttendance(t, router, hostSessionID)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		attendeeSessionID := RegisterAndLoginUser(t, userSrvc, "Attendee", "attendee@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		attendees, err := eventSrvc.GetAttendees(eventID)
		require.NoError(t, err)
		require.Len(t, attendees, 1)
	})
}

func TestAttendance_RegisterUnauthorized(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForAttendance(t, router, hostSessionID)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		req := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAttendance_UnregisterSuccess(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForAttendance(t, router, hostSessionID)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		attendeeSessionID := RegisterAndLoginUser(t, userSrvc, "Attendee", "attendee@example.com", "Secret123!")

		registerReq := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		registerReq.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, registerReq)
		require.Equal(t, http.StatusOK, w.Code)

		unregisterReq := httptest.NewRequest(http.MethodDelete, "/api/events/"+eventID+"/unregister", nil)
		unregisterReq.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w = httptest.NewRecorder()

		router.Handler.ServeHTTP(w, unregisterReq)

		require.Equal(t, http.StatusOK, w.Code)

		attendees, err := eventSrvc.GetAttendees(eventID)
		require.NoError(t, err)
		require.Len(t, attendees, 0)
	})
}

func TestAttendance_UnregisterUnauthorized(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForAttendance(t, router, hostSessionID)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		req := httptest.NewRequest(http.MethodDelete, "/api/events/"+eventID+"/unregister", nil)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAttendance_MultipleAttendees(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForAttendance(t, router, hostSessionID)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		attendee1SessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")
		attendee2SessionID := RegisterAndLoginUser(t, userSrvc, "Bob", "bob@example.com", "Secret123!")

		req1 := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		req1.AddCookie(&http.Cookie{Name: "session-id", Value: attendee1SessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req1)
		require.Equal(t, http.StatusOK, w.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		req2.AddCookie(&http.Cookie{Name: "session-id", Value: attendee2SessionID})
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req2)
		require.Equal(t, http.StatusOK, w.Code)

		attendees, err := eventSrvc.GetAttendees(eventID)
		require.NoError(t, err)
		require.Len(t, attendees, 2)
	})
}

func TestAttendance_RegisterNonExistentEvent(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodPost, "/api/events/00000000-0000-0000-0000-000000000000/register", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAttendance_DoubleRegistration(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForAttendance(t, router, hostSessionID)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		attendeeSessionID := RegisterAndLoginUser(t, userSrvc, "Attendee", "attendee@example.com", "Secret123!")

		req1 := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		req1.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req1)
		require.Equal(t, http.StatusOK, w.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		req2.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req2)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAttendance_UnregisterWhenNotRegistered(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForAttendance(t, router, hostSessionID)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		attendeeSessionID := RegisterAndLoginUser(t, userSrvc, "Attendee", "attendee@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodDelete, "/api/events/"+eventID+"/unregister", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAttendance_UnregisterNonExistentEvent(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodDelete, "/api/events/00000000-0000-0000-0000-000000000000/unregister", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func createEventForAttendance(t *testing.T, router *webapi.Router, sessionID string) {
	t.Helper()
	req := webapi.CreateEventRequest{
		Name:        "Test Event",
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

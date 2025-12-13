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

func TestDeleteEvent_Success(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForDelete(t, router, hostSessionID, "Test Event")

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
		eventID := events[0].EventID

		req := httptest.NewRequest(http.MethodDelete, "/api/events/"+eventID, nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: hostSessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		events, err = eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 0)
	})
}

func TestDeleteEvent_Unauthorized(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForDelete(t, router, hostSessionID, "Test Event")

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		req := httptest.NewRequest(http.MethodDelete, "/api/events/"+eventID, nil)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)

		events, err = eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
	})
}

func TestDeleteEvent_NotOrganizer(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		host1SessionID := registerHostAndLogin(t, userSrvc, "host1@example.com", "Secret123!")
		createEventForDelete(t, router, host1SessionID, "Test Event")

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		host2SessionID := registerHostAndLoginWithName(t, userSrvc, "Host 2", "host2@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodDelete, "/api/events/"+eventID, nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: host2SessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)

		events, err = eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
	})
}

func TestDeleteEvent_AttendeeCannotDelete(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForDelete(t, router, hostSessionID, "Test Event")

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		attendeeSessionID := RegisterAndLoginUser(t, userSrvc, "Attendee", "attendee@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodDelete, "/api/events/"+eventID, nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)

		events, err = eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
	})
}

func TestDeleteEvent_NonExistent(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodDelete, "/api/events/00000000-0000-0000-0000-000000000000", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: hostSessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteEvent_WithAttendees(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		hostSessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")
		createEventForDelete(t, router, hostSessionID, "Test Event")

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		eventID := events[0].EventID

		attendeeSessionID := RegisterAndLoginUser(t, userSrvc, "Attendee", "attendee@example.com", "Secret123!")
		registerReq := httptest.NewRequest(http.MethodPost, "/api/events/"+eventID+"/register", nil)
		registerReq.AddCookie(&http.Cookie{Name: "session-id", Value: attendeeSessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, registerReq)
		require.Equal(t, http.StatusOK, w.Code)

		attendees, err := eventSrvc.GetAttendees(eventID)
		require.NoError(t, err)
		require.Len(t, attendees, 1)

		req := httptest.NewRequest(http.MethodDelete, "/api/events/"+eventID, nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: hostSessionID})
		w = httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		events, err = eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 0)
	})
}

func createEventForDelete(t *testing.T, router *webapi.Router, sessionID string, name string) {
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

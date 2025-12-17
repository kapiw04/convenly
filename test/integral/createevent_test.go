package integral

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapiw04/convenly/internal/app"
	"github.com/kapiw04/convenly/internal/infra/webapi"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestCreateEvent_Success(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "bob@example.com", "Secret123!")

		sessionID, err := userSrvc.Login("bob@example.com", "Secret123!")
		require.NoError(t, err)

		body := createEventRequest(t)
		req := newEventRequest(t, body, sessionID)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, []string{"Music"}, events[0].Tags)
	})
}

func TestCreateEvent_Unauthorized(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		body := createEventRequest(t)
		req := newEventRequest(t, body, "")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCreateEvent_InvalidJSON(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "bob@example.com", "Secret123!")

		sessionID, err := userSrvc.Login("bob@example.com", "Secret123!")
		require.NoError(t, err)

		invalidJSON := []byte(`{"name": "Event", "date": invalid}`)

		req := newEventRequest(t, invalidJSON, sessionID)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateEvent_InvalidDateFormat(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "bob@example.com", "Secret123!")

		sessionID, err := userSrvc.Login("bob@example.com", "Secret123!")
		require.NoError(t, err)

		req := webapi.CreateEventRequest{
			Name:        "Event Name",
			Description: "Event desc",
			Latitude:    42.0,
			Longitude:   21.37,
			Fee:         10.0,
			Date:        "2005-04-02",
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)
		httpReq := newEventRequest(t, body, sessionID)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, httpReq)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCreateEvent_NotHost(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := userSrvc.Register("Bobby", "bob@example.com", "Secret123!")
		require.NoError(t, err)

		sessionID, err := userSrvc.Login("bob@example.com", "Secret123!")
		require.NoError(t, err)

		body := createEventRequest(t)
		req := newEventRequest(t, body, sessionID)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestCreateEvent_WithMinimumFee(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "bob@example.com", "Secret123!")

		sessionID, err := userSrvc.Login("bob@example.com", "Secret123!")
		require.NoError(t, err)

		req := webapi.CreateEventRequest{
			Name:        "Free Event",
			Description: "This event is free",
			Latitude:    42.0,
			Longitude:   21.37,
			Fee:         0.0,
			Date:        "2025-12-31T23:59:59Z",
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)
		httpReq := newEventRequest(t, body, sessionID)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, httpReq)

		require.Equal(t, http.StatusCreated, w.Code)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, float32(0.0), events[0].Fee)
	})
}

func TestCreateEvent_MultipleEvents(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "bob@example.com", "Secret123!")

		sessionID, err := userSrvc.Login("bob@example.com", "Secret123!")
		require.NoError(t, err)

		req1 := webapi.CreateEventRequest{
			Name:        "First Event",
			Description: "First event desc",
			Latitude:    42.0,
			Longitude:   21.37,
			Fee:         10.0,
			Date:        "2025-04-02T21:37:00Z",
		}
		body1, err := json.Marshal(req1)
		require.NoError(t, err)
		httpReq1 := newEventRequest(t, body1, sessionID)
		w1 := httptest.NewRecorder()

		router.Handler.ServeHTTP(w1, httpReq1)
		require.Equal(t, http.StatusCreated, w1.Code)

		req2 := webapi.CreateEventRequest{
			Name:        "Second Event",
			Description: "Second event desc",
			Latitude:    43.0,
			Longitude:   22.37,
			Fee:         20.0,
			Date:        "2025-05-02T21:37:00Z",
		}
		body2, err := json.Marshal(req2)
		require.NoError(t, err)
		httpReq2 := newEventRequest(t, body2, sessionID)
		w2 := httptest.NewRecorder()

		router.Handler.ServeHTTP(w2, httpReq2)
		require.Equal(t, http.StatusCreated, w2.Code)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 2)
	})
}

func TestCreateEvent_WithInvalidSessionID(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)
	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		body := createEventRequest(t)
		req := newEventRequest(t, body, "invalid-session-id")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestCreateEvent_NonExistentTag(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		registerAndPromoteHost(t, userSrvc, "bob@example.com", "Secret123!")
		sessionID, err := userSrvc.Login("bob@example.com", "Secret123!")
		require.NoError(t, err)

		req := webapi.CreateEventRequest{
			Name:        "Event Name",
			Description: "Event desc",
			Latitude:    42.0,
			Longitude:   21.37,
			Fee:         10.0,
			Date:        "2025-12-31T23:59:59Z",
			Tags:        []string{"NonExistentTag123"},
		}
		body, err := json.Marshal(req)
		require.NoError(t, err)

		httpReq := newEventRequest(t, body, sessionID)
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, httpReq)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func registerAndPromoteHost(t *testing.T, userSrvc *app.UserService, email, password string) {
	t.Helper()
	err := userSrvc.Register("Bobby", email, password)
	require.NoError(t, err)
	user, err := userSrvc.GetByEmail(email)
	require.NoError(t, err)
	err = userSrvc.PromoteToHost(user.UUID.String())
	require.NoError(t, err)
}

func createEventRequest(t *testing.T) []byte {
	t.Helper()
	req := webapi.CreateEventRequest{
		Name:        "Event Name",
		Description: "Event desc",
		Latitude:    42.0,
		Longitude:   21.37,
		Fee:         10.0,
		Date:        "2005-04-02T21:37:00Z",
		Tags:        []string{"Music"},
	}
	body, err := json.Marshal(req)
	require.NoError(t, err)
	return body
}

func newEventRequest(t *testing.T, body []byte, sessionID string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/events/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if sessionID != "" {
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
	}
	return req
}

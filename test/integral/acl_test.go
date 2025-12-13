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

func TestAclMiddleware_AttendeeCannotCreateEvent(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		req := webapi.CreateEventRequest{
			Name:        "Test Event",
			Description: "Test description",
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

		require.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestAclMiddleware_HostCanCreateEvent(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := registerHostAndLogin(t, userSrvc, "host@example.com", "Secret123!")

		req := webapi.CreateEventRequest{
			Name:        "Host Event",
			Description: "Test description",
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

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
	})
}

func TestAclMiddleware_UnauthenticatedCannotAccessProtectedRoute(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		httpReq := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, httpReq)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

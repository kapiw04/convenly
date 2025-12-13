package integral

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapiw04/convenly/internal/domain/user"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestBecomeHost_Success(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodPost, "/api/become-host", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		u, err := userSrvc.GetByEmail("alice@example.com")
		require.NoError(t, err)
		require.Equal(t, user.HOST, u.Role)
	})
}

func TestBecomeHost_Unauthorized(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		req := httptest.NewRequest(http.MethodPost, "/api/become-host", nil)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestBecomeHost_VerifyCanCreateEvent(t *testing.T) {
	sqlDb, userSrvc, eventSrvc, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		becomeHostReq := httptest.NewRequest(http.MethodPost, "/api/become-host", nil)
		becomeHostReq.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, becomeHostReq)
		require.Equal(t, http.StatusOK, w.Code)

		createTestEventViaAPI(t, router, sessionID, "My Event", "2025-12-31T23:59:59Z", 15.0, []string{"Music"})

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "My Event", events[0].Name)
	})
}

func TestBecomeHost_RolePersistedAfterRelogin(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		becomeHostReq := httptest.NewRequest(http.MethodPost, "/api/become-host", nil)
		becomeHostReq.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, becomeHostReq)
		require.Equal(t, http.StatusOK, w.Code)

		logoutReq := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
		logoutReq.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, logoutReq)
		require.Equal(t, http.StatusOK, w.Code)

		newSessionID, err := userSrvc.Login("alice@example.com", "Secret123!")
		require.NoError(t, err)

		meReq := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		meReq.AddCookie(&http.Cookie{Name: "session-id", Value: newSessionID})
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, meReq)
		require.Equal(t, http.StatusOK, w.Code)

		var userResp struct {
			Role user.Role `json:"role"`
		}
		err = json.NewDecoder(w.Body).Decode(&userResp)
		require.NoError(t, err)
		require.Equal(t, user.HOST, userResp.Role)
	})
}

func TestBecomeHost_AlreadyHost(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodPost, "/api/become-host", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		req = httptest.NewRequest(http.MethodPost, "/api/become-host", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		u, err := userSrvc.GetByEmail("alice@example.com")
		require.NoError(t, err)
		require.Equal(t, user.HOST, u.Role)
	})
}

func TestBecomeHost_InvalidSession(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		req := httptest.NewRequest(http.MethodPost, "/api/become-host", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: "invalid-session-id"})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

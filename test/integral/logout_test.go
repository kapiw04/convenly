package integral

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestLogout_Success(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestLogout_SessionInvalidatedAfterLogout(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")

		logoutReq := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
		logoutReq.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()
		router.Handler.ServeHTTP(w, logoutReq)
		require.Equal(t, http.StatusOK, w.Code)

		meReq := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		meReq.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w = httptest.NewRecorder()
		router.Handler.ServeHTTP(w, meReq)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestLogout_Unauthorized(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestLogout_InvalidSession(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: "invalid-session-id"})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

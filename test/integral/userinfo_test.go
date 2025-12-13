package integral

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapiw04/convenly/internal/domain/user"
	"github.com/stretchr/testify/require"
)

func Test_GetUserInfo(t *testing.T) {
	sqlDb, userSrvc, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		sessionID := RegisterAndLoginUser(t, userSrvc, "Alice", "alice@example.com", "Secret123!")
		req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		var userResp struct {
			UUID  string    `json:"uuid"`
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Role  user.Role `json:"role"`
		}
		err := json.NewDecoder(w.Body).Decode(&userResp)
		require.NoError(t, err)
		require.Equal(t, "Alice", userResp.Name)
		require.Equal(t, "alice@example.com", userResp.Email)
		require.Equal(t, user.ATTENDEE, userResp.Role)
		require.NotEmpty(t, userResp.UUID)
	})
}

func Test_GetUserInfo_Unauthorized(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func Test_GetUserInfo_InvalidSession(t *testing.T) {
	sqlDb, _, _, router := setupAllServices(t)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		req.AddCookie(&http.Cookie{Name: "session-id", Value: "invalid-session-id"})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

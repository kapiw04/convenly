package integral

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kapiw04/convenly/internal/infra/webapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestLogin_Success(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	router := webapi.NewRouter(srvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := srvc.Register(
			"Bob",
			"bob@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		loginReq := webapi.LoginRequest{
			Email:    "bob@example.com",
			Password: "Secret123!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "sessionId")
		assert.NotEmpty(t, response["sessionId"])
	})
}

func TestLogin_InvalidPassword(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	router := webapi.NewRouter(srvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		// First, register a user
		err := srvc.Register(
			"Bob",
			"bob@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		// Try to login with wrong password
		loginReq := webapi.LoginRequest{
			Email:    "bob@example.com",
			Password: "WrongPassword!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin_NonExistentUser(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	router := webapi.NewRouter(srvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		// Try to login with a user that doesn't exist
		loginReq := webapi.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "SomePassword123!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin_EmptyEmail(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	router := webapi.NewRouter(srvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		// Try to login with empty email
		loginReq := webapi.LoginRequest{
			Email:    "",
			Password: "SomePassword123!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Equal(t, "empty fields", response["error"])
	})
}

func TestLogin_EmptyPassword(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	router := webapi.NewRouter(srvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		// Try to login with empty password
		loginReq := webapi.LoginRequest{
			Email:    "bob@example.com",
			Password: "",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
		assert.Equal(t, "empty fields", response["error"])
	})
}

func TestLogin_CaseInsensitiveEmail(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	router := webapi.NewRouter(srvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		// Register a user
		err := srvc.Register(
			"Bob",
			"bob@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		// Try to login with uppercase email
		loginReq := webapi.LoginRequest{
			Email:    "Bob@EXAMPLE.COM",
			Password: "Secret123!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "sessionId")
		assert.NotEmpty(t, response["sessionId"])
	})
}

func TestLogin_EmailWithWhitespace(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	router := webapi.NewRouter(srvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		// Register a user
		err := srvc.Register(
			"Bob",
			"bob@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		// Try to login with email that has whitespace
		loginReq := webapi.LoginRequest{
			Email:    "  bob@example.com  ",
			Password: "Secret123!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "sessionId")
		assert.NotEmpty(t, response["sessionId"])
	})
}

func TestLogin_InvalidJSON(t *testing.T) {
	sqlDb := setupDb(t)
	srvc := setupUserService(t, sqlDb)
	router := webapi.NewRouter(srvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		// Send invalid JSON
		body := []byte(`{"email": "bob@example.com", "password":`)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response, "error")
	})
}

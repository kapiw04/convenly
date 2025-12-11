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
	userSrvc := setupUserService(t, sqlDb)
	eventSrvc := setupEventService(t, sqlDb)
	router := webapi.NewRouter(userSrvc, eventSrvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := userSrvc.Register(
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

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
	})
}

func TestLogin_InvalidPassword(t *testing.T) {
	sqlDb := setupDb(t)
	userSrvc := setupUserService(t, sqlDb)
	eventSrvc := setupEventService(t, sqlDb)
	router := webapi.NewRouter(userSrvc, eventSrvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := userSrvc.Register(
			"Bob",
			"bob@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		loginReq := webapi.LoginRequest{
			Email:    "bob@example.com",
			Password: "WrongPassword!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin_NonExistentUser(t *testing.T) {
	sqlDb := setupDb(t)
	userSrvc := setupUserService(t, sqlDb)
	eventSrvc := setupEventService(t, sqlDb)
	router := webapi.NewRouter(userSrvc, eventSrvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		loginReq := webapi.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "SomePassword123!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin_EmptyEmail(t *testing.T) {
	sqlDb := setupDb(t)
	userSrvc := setupUserService(t, sqlDb)
	eventSrvc := setupEventService(t, sqlDb)
	router := webapi.NewRouter(userSrvc, eventSrvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		loginReq := webapi.LoginRequest{
			Email:    "",
			Password: "SomePassword123!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
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
	userSrvc := setupUserService(t, sqlDb)
	eventSrvc := setupEventService(t, sqlDb)
	router := webapi.NewRouter(userSrvc, eventSrvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		loginReq := webapi.LoginRequest{
			Email:    "bob@example.com",
			Password: "",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
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
	userSrvc := setupUserService(t, sqlDb)
	eventSrvc := setupEventService(t, sqlDb)
	router := webapi.NewRouter(userSrvc, eventSrvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := userSrvc.Register(
			"Bob",
			"bob@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		loginReq := webapi.LoginRequest{
			Email:    "Bob@EXAMPLE.COM",
			Password: "Secret123!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
	})
}

func TestLogin_EmailWithWhitespace(t *testing.T) {
	sqlDb := setupDb(t)
	userSrvc := setupUserService(t, sqlDb)
	eventSrvc := setupEventService(t, sqlDb)
	router := webapi.NewRouter(userSrvc, eventSrvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		err := userSrvc.Register(
			"Bob",
			"bob@example.com",
			"Secret123!",
		)
		require.NoError(t, err)

		loginReq := webapi.LoginRequest{
			Email:    "  bob@example.com  ",
			Password: "Secret123!",
		}
		body, err := json.Marshal(loginReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
	})
}

func TestLogin_InvalidJSON(t *testing.T) {
	sqlDb := setupDb(t)
	userSrvc := setupUserService(t, sqlDb)
	eventSrvc := setupEventService(t, sqlDb)
	router := webapi.NewRouter(userSrvc, eventSrvc)

	WithTx(t, sqlDb, func(t *testing.T, tx *sql.Tx) {
		body := []byte(`{"email": "bob@example.com", "password":`)

		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(body))
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

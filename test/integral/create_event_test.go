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

func TestCreateEvent_Success(t *testing.T) {
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

		sessionID, err := userSrvc.Login("bob@example.com", "Secret123!")
		require.NoError(t, err)

		addEventReq := webapi.CreateEventRequest{
			Name:        "Event Name",
			Description: "Event desc",
			Latitude:    42.0,
			Longitude:   21.37,
			Fee:         10.0,
			Date:        "2005-04-02T21:37:00Z",
		}
		body, err := json.Marshal(addEventReq)
		require.NoError(t, err)
		require.NotNil(t, body)

		req := httptest.NewRequest(http.MethodPost, "/events/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "session-id", Value: sessionID})
		w := httptest.NewRecorder()

		router.Handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		events, err := eventSrvc.GetAllEvents()
		require.NoError(t, err)
		require.Len(t, events, 1)
	})
}

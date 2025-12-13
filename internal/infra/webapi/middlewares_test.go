package webapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/kapiw04/convenly/internal/app"
	mock_security "github.com/kapiw04/convenly/internal/domain/security/mocks"
	"github.com/kapiw04/convenly/internal/domain/user"
	mock_user "github.com/kapiw04/convenly/internal/domain/user/mocks"
	"github.com/kapiw04/convenly/internal/infra/webapi"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthMiddleware_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockUserRepo := mock_user.NewMockUserRepo(ctrl)
	mockSessionRepo := mock_user.NewMockSessionRepo(ctrl)
	mockHasher := mock_security.NewMockHasher(ctrl)
	userSrvc := app.NewUserService(mockUserRepo, mockSessionRepo, mockHasher)

	testUser := user.User{
		UUID:  uuid.New(),
		Name:  "Alice",
		Email: "alice@example.com",
		Role:  user.ATTENDEE,
	}

	mockSessionRepo.
		EXPECT().
		Get("valid-session-id").
		Return(testUser, nil).
		Times(1)

	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	})

	middleware := webapi.AuthMiddleware(userSrvc)
	handler := middleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session-id", Value: "valid-session-id"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, handlerCalled)
}

func TestAuthMiddleware_MissingCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockUserRepo := mock_user.NewMockUserRepo(ctrl)
	mockSessionRepo := mock_user.NewMockSessionRepo(ctrl)
	mockHasher := mock_security.NewMockHasher(ctrl)
	userSrvc := app.NewUserService(mockUserRepo, mockSessionRepo, mockHasher)

	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	middleware := webapi.AuthMiddleware(userSrvc)
	handler := middleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.False(t, handlerCalled)
}

func TestAuthMiddleware_InvalidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockUserRepo := mock_user.NewMockUserRepo(ctrl)
	mockSessionRepo := mock_user.NewMockSessionRepo(ctrl)
	mockHasher := mock_security.NewMockHasher(ctrl)
	userSrvc := app.NewUserService(mockUserRepo, mockSessionRepo, mockHasher)

	mockSessionRepo.
		EXPECT().
		Get("invalid-session-id").
		Return(user.User{}, user.ErrUserNotFound).
		Times(1)

	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	middleware := webapi.AuthMiddleware(userSrvc)
	handler := middleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "session-id", Value: "invalid-session-id"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	require.False(t, handlerCalled)
}

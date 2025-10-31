package webapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kapiw04/convenly/internal/app"
	mock_security "github.com/kapiw04/convenly/internal/domain/security/mocks"
	mock_user "github.com/kapiw04/convenly/internal/domain/user/mocks"
	"github.com/kapiw04/convenly/internal/infra/webapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupMockController(t *testing.T) *gomock.Controller {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	return ctrl
}

func setupMockService(t *testing.T, ctrl *gomock.Controller) (*mock_user.MockUserRepo, *mock_security.MockHasher, *app.UserService) {
	t.Helper()
	mockRepo := mock_user.NewMockUserRepo(ctrl)
	mockHasher := mock_security.NewMockHasher(ctrl)
	mockSvc := app.NewUserService(mockRepo, mockHasher)
	return mockRepo, mockHasher, mockSvc
}

func setupServer(t *testing.T, srvc *app.UserService) *httptest.Server {
	t.Helper()
	mux := chi.NewRouter()
	rt := &webapi.Router{UserService: srvc, Handler: mux}
	mux.Post("/register", rt.RegisterHandler)
	srv := httptest.NewServer(rt.Handler)
	t.Cleanup(srv.Close)
	return srv
}

func TestRegister_Success(t *testing.T) {
	ctrl := setupMockController(t)
	mockRepo, mockHasher, mockSvc := setupMockService(t, ctrl)

	mockHasher.
		EXPECT().
		Hash(gomock.Any()).
		Return("hashed", nil).
		Times(1)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(1)

	srv := setupServer(t, mockSvc)

	body := `{"name":"Alice","email":"alice@example.com","password":"Secret123!"}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NotNil(t, res)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)
}

func TestRegister_EmptyFields(t *testing.T) {
	ctrl := setupMockController(t)
	mockRepo, _, mockSvc := setupMockService(t, ctrl)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(0)

	srv := setupServer(t, mockSvc)

	body := `{"name":"Alice","email":"alice@example.com","password":""}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestRegister_MissingFields(t *testing.T) {
	ctrl := setupMockController(t)
	mockRepo, _, mockSvc := setupMockService(t, ctrl)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(0)

	srv := setupServer(t, mockSvc)

	body := `{"name":"Alice","email":"alice@example.com"}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestRegister_InvalidEmail(t *testing.T) {
	ctrl := setupMockController(t)
	mockRepo, _, mockSvc := setupMockService(t, ctrl)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(0)

	srv := setupServer(t, mockSvc)

	body := `{"name":"Alice","email":"not-an-email","password":"Secret123!"}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestRegister_WeakPassword(t *testing.T) {
	ctrl := setupMockController(t)
	mockRepo, _, mockSvc := setupMockService(t, ctrl)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(0)

	srv := setupServer(t, mockSvc)

	body := `{"name":"Alice","email":"alice@example.com","password":"password"}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestRegister_PasswordTooShort(t *testing.T) {
	ctrl := setupMockController(t)
	mockRepo, _, mockSvc := setupMockService(t, ctrl)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(0)

	srv := setupServer(t, mockSvc)

	body := `{"name":"Alice","email":"alice@example.com","password":"Sec1!"}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestRegister_PasswordTooLong(t *testing.T) {
	ctrl := setupMockController(t)
	mockRepo, _, mockSvc := setupMockService(t, ctrl)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(0)

	srv := setupServer(t, mockSvc)

	body := `{"name":"Alice","email":"alice@example.com","password":"ExtreeeemeeelyLongPassword1234!"}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

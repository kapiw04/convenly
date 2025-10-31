package webapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kapiw04/convenly/internal/app"
	mock_user "github.com/kapiw04/convenly/internal/domain/user/mocks"
	"github.com/kapiw04/convenly/internal/infra/webapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func newServer(t *testing.T, srvc app.UserService) *httptest.Server {
	t.Helper()
	mux := chi.NewRouter()
	rt := &webapi.Router{UserService: &srvc, Handler: mux}
	mux.Post("/register", rt.RegisterHandler)
	return httptest.NewServer(rt.Handler)
}

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_user.NewMockUserRepo(ctrl)
	mockSvc := *app.NewUserService(mockRepo)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(1)

	srv := newServer(t, mockSvc)
	t.Cleanup(srv.Close)

	body := `{"name":"Alice","email":"alice@example.com","password":"Secret123!"}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NotNil(t, res)
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)
}

func TestRegister_EmptyFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_user.NewMockUserRepo(ctrl)
	mockSvc := *app.NewUserService(mockRepo)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(0)

	srv := newServer(t, mockSvc)
	t.Cleanup(srv.Close)

	body := `{"name":"Alice","email":"alice@example.com","password":""}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

func TestRegister_MissingFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_user.NewMockUserRepo(ctrl)
	mockSvc := *app.NewUserService(mockRepo)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil).
		Times(0)

	srv := newServer(t, mockSvc)
	t.Cleanup(srv.Close)

	body := `{"name":"Alice","email":"alice@example.com"}`
	res, err := http.Post(srv.URL+"/register", "application/json", strings.NewReader(body))
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusBadRequest)
}

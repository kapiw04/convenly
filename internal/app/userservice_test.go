package app

import (
	"errors"
	"testing"

	mock_security "github.com/kapiw04/convenly/internal/domain/security/mocks"
	"github.com/kapiw04/convenly/internal/domain/user"
	mock_user "github.com/kapiw04/convenly/internal/domain/user/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	hasher.EXPECT().Hash("Password123!").Return("hashedpassword", nil)
	userRepo.EXPECT().Save(gomock.Any()).Return(nil)

	svc := NewUserService(userRepo, sessionRepo, hasher)
	err := svc.Register("TestUser", "test@example.com", "Password123!")

	require.NoError(t, err)
}

func TestUserService_Register_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	svc := NewUserService(userRepo, sessionRepo, hasher)
	err := svc.Register("TestUser", "invalid-email", "Password123!")

	require.Error(t, err)
}

func TestUserService_Register_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	svc := NewUserService(userRepo, sessionRepo, hasher)
	err := svc.Register("TestUser", "test@example.com", "short")

	require.Error(t, err)
}

func TestUserService_Register_HashError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	hasher.EXPECT().Hash("Password123!").Return("", errors.New("hashing failed"))

	svc := NewUserService(userRepo, sessionRepo, hasher)
	err := svc.Register("TestUser", "test@example.com", "Password123!")

	require.Error(t, err)
}

func TestUserService_Register_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	hasher.EXPECT().Hash("Password123!").Return("hashedpassword", nil)
	userRepo.EXPECT().Save(gomock.Any()).Return(errors.New("database error"))

	svc := NewUserService(userRepo, sessionRepo, hasher)
	err := svc.Register("TestUser", "test@example.com", "Password123!")

	require.Error(t, err)
}

func TestUserService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	testUser := &user.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	userRepo.EXPECT().FindByEmail("test@example.com").Return(testUser, nil)
	hasher.EXPECT().Compare("Password123!", "hashedpassword").Return(true)
	sessionRepo.EXPECT().Create("test@example.com").Return("session-id", nil)

	svc := NewUserService(userRepo, sessionRepo, hasher)
	sessionID, err := svc.Login("test@example.com", "Password123!")

	require.NoError(t, err)
	require.Equal(t, "session-id", sessionID)
}

func TestUserService_Login_InvalidEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	svc := NewUserService(userRepo, sessionRepo, hasher)
	_, err := svc.Login("invalid-email", "Password123!")

	require.Error(t, err)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	svc := NewUserService(userRepo, sessionRepo, hasher)
	_, err := svc.Login("test@example.com", "short")

	require.Error(t, err)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	userRepo.EXPECT().FindByEmail("test@example.com").Return(nil, errors.New("user not found"))

	svc := NewUserService(userRepo, sessionRepo, hasher)
	_, err := svc.Login("test@example.com", "Password123!")

	require.Error(t, err)
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	testUser := &user.User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	userRepo.EXPECT().FindByEmail("test@example.com").Return(testUser, nil)
	hasher.EXPECT().Compare("WrongPassword123!", "hashedpassword").Return(false)

	svc := NewUserService(userRepo, sessionRepo, hasher)
	_, err := svc.Login("test@example.com", "WrongPassword123!")

	require.ErrorIs(t, err, user.ErrInvalidCredentials)
}

func TestUserService_Logout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	sessionRepo.EXPECT().Delete("session-id").Return(nil)

	svc := NewUserService(userRepo, sessionRepo, hasher)
	err := svc.Logout("session-id")

	require.NoError(t, err)
}

func TestUserService_PromoteToHost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	testUser := &user.User{
		Name: "TestUser",
		Role: user.ATTENDEE,
	}

	userRepo.EXPECT().FindByUUID("user-uuid").Return(testUser, nil)
	userRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(u *user.User) error {
		require.Equal(t, user.HOST, u.Role)
		return nil
	})

	svc := NewUserService(userRepo, sessionRepo, hasher)
	err := svc.PromoteToHost("user-uuid")

	require.NoError(t, err)
}

func TestUserService_PromoteToHost_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	userRepo.EXPECT().FindByUUID("nonexistent-uuid").Return(nil, errors.New("user not found"))

	svc := NewUserService(userRepo, sessionRepo, hasher)
	err := svc.PromoteToHost("nonexistent-uuid")

	require.Error(t, err)
}

func TestUserService_GetBySessionID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	testUser := user.User{
		Name:  "TestUser",
		Email: "test@example.com",
	}

	sessionRepo.EXPECT().Get("session-id").Return(testUser, nil)

	svc := NewUserService(userRepo, sessionRepo, hasher)
	u, err := svc.GetBySessionID("session-id")

	require.NoError(t, err)
	require.Equal(t, "TestUser", u.Name)
}

func TestUserService_GetBySessionID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepo := mock_user.NewMockUserRepo(ctrl)
	sessionRepo := mock_user.NewMockSessionRepo(ctrl)
	hasher := mock_security.NewMockHasher(ctrl)

	sessionRepo.EXPECT().Get("invalid-session").Return(user.User{}, errors.New("session not found"))

	svc := NewUserService(userRepo, sessionRepo, hasher)
	_, err := svc.GetBySessionID("invalid-session")

	require.Error(t, err)
}

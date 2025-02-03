package usecase

import (
	"auth_records/internal/auth/apperrors"
	"auth_records/internal/auth/domain"
	"auth_records/pkg/password"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) UserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateTokenWithID(userID uint64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) CheckJWTAndGetUserID(tokenString string) (uint64, error) {
	args := m.Called(tokenString)
	return args.Get(0).(uint64), args.Error(1)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	logger, _ := zap.NewProduction()
	usecase := New(logger, mockRepo, mockTokenService)

	email := "test@mail.com"
	pswd := "correctPassword"
	expectedToken := "someToken"

	hashedPassword, err := password.HashPassword([]byte(pswd))
	if err != nil {
		t.Fatal("Failed to hash password:", err)
	}

	mockUser := &domain.User{
		ID:           1,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	mockRepo.On("UserByEmail", mock.Anything, email).Return(mockUser, nil)
	mockTokenService.On("GenerateTokenWithID", mockUser.ID).Return(expectedToken, nil)

	ctx := context.Background()
	token, err := usecase.Login(ctx, email, pswd)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)

	mockRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	logger, _ := zap.NewProduction()
	usecase := New(logger, mockRepo, mockTokenService)

	email := "notfound@mail.com"
	pswd := "somePassword"
	expectedError := apperrors.ErrUserNotFound

	mockRepo.On("UserByEmail", mock.Anything, email).Return(nil, expectedError)

	ctx := context.Background()
	token, err := usecase.Login(ctx, email, pswd)

	assert.Empty(t, token)
	assert.ErrorIs(t, err, expectedError)

	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	logger, _ := zap.NewProduction()
	usecase := New(logger, mockRepo, mockTokenService)

	email := "test@mail.com"
	pswd := "wrongPassword"
	corpswd := "correctPassword"
	expectedError := apperrors.ErrInvalidLoginCredantials

	hashedPassword, err := password.HashPassword([]byte(corpswd))
	if err != nil {
		t.Fatal("Failed to hash password:", err)
	}

	mockUser := &domain.User{
		ID:           1,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	mockRepo.On("UserByEmail", mock.Anything, email).Return(mockUser, nil)

	ctx := context.Background()
	token, err := usecase.Login(ctx, email, pswd)

	assert.Empty(t, token)
	assert.ErrorIs(t, err, expectedError)

	mockRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

func TestLogin_TokenGenerationFailure(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenService := new(MockTokenService)
	logger, _ := zap.NewProduction()
	usecase := New(logger, mockRepo, mockTokenService)

	email := "test@mail.com"
	pswd := "correctPassword"
	expectedError := apperrors.ErrInternalServer

	hashedPassword, err := password.HashPassword([]byte(pswd))
	if err != nil {
		t.Fatal("Failed to hash password:", err)
	}

	mockUser := &domain.User{
		ID:           1,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	mockRepo.On("UserByEmail", mock.Anything, email).Return(mockUser, nil)
	mockTokenService.On("GenerateTokenWithID", mockUser.ID).Return("", expectedError)

	ctx := context.Background()
	token, err := usecase.Login(ctx, email, pswd)

	assert.Empty(t, token)
	assert.ErrorIs(t, err, expectedError)

	mockRepo.AssertExpectations(t)
	mockTokenService.AssertExpectations(t)
}

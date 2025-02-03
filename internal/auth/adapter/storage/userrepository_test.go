package storage

import (
	"auth_records/internal/auth/adapter/storage/dto"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorageProvider struct {
	mock.Mock
}

func (m *MockStorageProvider) UserByEmail(ctx context.Context, email string) (*dto.LoginUser, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginUser), args.Error(1)
}

func TestUserByEmail_Success(t *testing.T) {
	mockProvider := new(MockStorageProvider)
	repo := NewUserRepository(mockProvider)

	email := "test@mail.com"
	expectedUser := &dto.LoginUser{
		ID:       1,
		Email:    email,
		Username: "TestUser",
	}

	mockProvider.On("UserByEmail", mock.Anything, email).Return(expectedUser, nil)

	ctx := context.Background()
	user, err := repo.UserByEmail(ctx, email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Username, user.Username)

	mockProvider.AssertExpectations(t)
}

func TestUserByEmail_Failure(t *testing.T) {
	mockProvider := new(MockStorageProvider)
	repo := NewUserRepository(mockProvider)

	email := "notfound@mail.com"
	expectedError := errors.New("user not found")

	mockProvider.On("UserByEmail", mock.Anything, email).Return(nil, expectedError)

	ctx := context.Background()
	user, err := repo.UserByEmail(ctx, email)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, expectedError)

	mockProvider.AssertExpectations(t)
}

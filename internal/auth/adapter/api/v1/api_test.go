package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"auth_records/internal/auth/adapter/api/v1/dto"
	v1 "auth_records/pkg/records_grpc/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Login(ctx context.Context, email string, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

type MockRecordsClient struct {
	mock.Mock
}

func (m *MockRecordsClient) GetRandomRecords(token string) ([]*v1.Record, error) {
	args := m.Called(token)
	return args.Get(0).([]*v1.Record), args.Error(1)
}

func TestLogin_Success(t *testing.T) {
	mockAuthUseCase := new(MockAuthUseCase)
	mockRecordsClient := new(MockRecordsClient)
	api := New(zap.NewNop(), mockAuthUseCase, mockRecordsClient)

	email := "test@mail.com"
	password := "valid-password"
	token := "valid-token"
	records := []*v1.Record{
		{Id: 1, Title: "Record 1"},
		{Id: 2, Title: "Record 2"},
	}

	mockAuthUseCase.On("Login", mock.Anything, email, password).Return(token, nil)
	mockRecordsClient.On("GetRandomRecords", token).Return(records, nil)

	loginRequest := dto.LoginRequest{Email: email, Password: password}
	body, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	api.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response dto.RecordsResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Records, 2)

	mockAuthUseCase.AssertExpectations(t)
	mockRecordsClient.AssertExpectations(t)
}

func TestLogin_Failure(t *testing.T) {
	mockAuthUseCase := new(MockAuthUseCase)
	mockRecordsClient := new(MockRecordsClient)
	api := New(zap.NewNop(), mockAuthUseCase, mockRecordsClient)

	email := "wrong@mail.com"
	password := "wrong-password"

	mockAuthUseCase.On("Login", mock.Anything, email, password).Return("", errors.New("unauthorized"))

	loginRequest := dto.LoginRequest{Email: email, Password: password}
	body, _ := json.Marshal(loginRequest)
	req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	api.Login(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")

	mockAuthUseCase.AssertExpectations(t)
}

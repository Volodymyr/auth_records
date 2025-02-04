package grpcgetterhandlers

import (
	"context"
	"errors"
	"testing"

	pb "auth_records/pkg/records_grpc/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockRecordsUseCase struct {
	mock.Mock
}

func (m *MockRecordsUseCase) Records(ctx context.Context) ([]*pb.Record, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*pb.Record), args.Error(1)
}

func TestGetRandomRecords_Success(t *testing.T) {
	mockUseCase := new(MockRecordsUseCase)
	handler := New(zap.NewNop(), mockUseCase)

	expectedRecords := []*pb.Record{
		{Id: 1, Title: "Record 1"},
		{Id: 2, Title: "Record 2"},
	}

	mockUseCase.On("Records", mock.Anything).Return(expectedRecords, nil)

	req := &pb.GetRandomRecordsRequest{}
	resp, err := handler.GetRandomRecords(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Records, 2)

	mockUseCase.AssertExpectations(t)
}

func TestGetRandomRecords_Failure(t *testing.T) {
	mockUseCase := new(MockRecordsUseCase)
	handler := New(zap.NewNop(), mockUseCase)

	mockUseCase.On("Records", mock.Anything).Return([]*pb.Record{}, errors.New("database error"))

	req := &pb.GetRandomRecordsRequest{}
	resp, err := handler.GetRandomRecords(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "database error", err.Error())

	mockUseCase.AssertExpectations(t)
}

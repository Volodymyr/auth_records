package usecase

import (
	pb "auth_records/pkg/records_grpc/v1"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockRecordsRepository struct {
	mock.Mock
}

func (m *MockRecordsRepository) Records(ctx context.Context) ([]*pb.Record, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*pb.Record), args.Error(1)
}

func TestUseCase_Records_Success(t *testing.T) {
	mockRepo := new(MockRecordsRepository)
	uc := New(zap.NewNop(), mockRepo)

	expectedRecords := []*pb.Record{
		{Id: 1, Title: "Record 1"},
		{Id: 2, Title: "Record 2"},
	}

	mockRepo.On("Records", mock.Anything).Return(expectedRecords, nil)

	records, err := uc.Records(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, records)
	assert.Len(t, records, 2)
	assert.Equal(t, expectedRecords, records)

	mockRepo.AssertExpectations(t)
}

func TestUseCase_Records_Failure(t *testing.T) {
	mockRepo := new(MockRecordsRepository)
	uc := New(zap.NewNop(), mockRepo)

	mockRepo.On("Records", mock.Anything).Return(nil, errors.New("failed to fetch records"))

	records, err := uc.Records(context.Background())

	assert.Error(t, err)
	assert.Nil(t, records)
	assert.Equal(t, "failed to fetch records", err.Error())

	mockRepo.AssertExpectations(t)
}

package storage

import (
	pb "auth_records/pkg/records_grpc/v1"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRecordsProvider struct {
	mock.Mock
}

func (m *MockRecordsProvider) Records(ctx context.Context) ([]*pb.Record, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*pb.Record), args.Error(1)
}

func TestRecordsRepository_Records_Success(t *testing.T) {
	mockProvider := new(MockRecordsProvider)
	repo := NewRecordsRepository(mockProvider)

	expectedRecords := []*pb.Record{
		{Id: 1, Title: "Test Record 1"},
		{Id: 2, Title: "Test Record 2"},
	}

	mockProvider.On("Records", mock.Anything).Return(expectedRecords, nil)

	records, err := repo.Records(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, records)
	assert.Len(t, records, 2)
	assert.Equal(t, expectedRecords, records)

	mockProvider.AssertExpectations(t)
}

func TestRecordsRepository_Records_Failure(t *testing.T) {
	mockProvider := new(MockRecordsProvider)
	repo := NewRecordsRepository(mockProvider)

	mockProvider.On("Records", mock.Anything).Return([]*pb.Record{}, errors.New("database error"))

	records, err := repo.Records(context.Background())

	assert.Error(t, err)
	assert.Nil(t, records)
	assert.Equal(t, "database error", err.Error())

	mockProvider.AssertExpectations(t)
}

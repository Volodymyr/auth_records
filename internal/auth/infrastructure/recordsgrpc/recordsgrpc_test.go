package recordsgrpc

import (
	"auth_records/internal/auth/apperrors"
	pb "auth_records/pkg/records_grpc/v1"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockRecordsServiceClient struct {
	pb.RecordsServiceClient
	mock.Mock
}

func (m *MockRecordsServiceClient) GetRandomRecords(ctx context.Context, in *pb.GetRandomRecordsRequest, opts ...grpc.CallOption) (*pb.GetRandomRecordsResponse, error) {
	args := m.Called(ctx, in)
	if respons, ok := args.Get(0).(*pb.GetRandomRecordsResponse); ok {
		return respons, nil
	}

	return nil, args.Error(1)
}

func TestGetRandomRecords_Success(t *testing.T) {
	mockClient := new(MockRecordsServiceClient)

	token := "valid-token"
	expectedRecords := []*pb.Record{
		{Id: 1, Title: "Record1"},
		{Id: 2, Title: "Record2"},
	}

	mockClient.On("GetRandomRecords", mock.Anything, &pb.GetRandomRecordsRequest{}).Return(&pb.GetRandomRecordsResponse{Records: expectedRecords}, nil)

	recordsGRPC := &recordsgrpc{
		getterClient: mockClient,
	}

	result, err := recordsGRPC.GetRandomRecords(token)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecords, result)

	mockClient.AssertExpectations(t)
}

func TestGetRandomRecords_Error(t *testing.T) {
	mockClient := new(MockRecordsServiceClient)

	token := "invalid-token"
	expectedError := apperrors.ErrInvalidLoginCredantials

	mockClient.On("GetRandomRecords", mock.Anything, &pb.GetRandomRecordsRequest{}).Return(nil, expectedError)

	recordsGRPC := &recordsgrpc{
		getterClient: mockClient,
	}

	result, err := recordsGRPC.GetRandomRecords(token)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockClient.AssertExpectations(t)
}

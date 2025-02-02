package grpcgetterhandlers

import (
	pb "auth_records/pkg/records_grpc/v1"
	"context"

	"go.uber.org/zap"
)

type RecordsUseCase interface {
	Records(ctx context.Context) ([]*pb.Record, error)
}

type GetterHandlers interface {
	GetRandomRecords(ctx context.Context, req *pb.GetRandomRecordsRequest) (*pb.GetRandomRecordsResponse, error)
}

type getterHandlers struct {
	log            *zap.Logger
	recordsUseCase RecordsUseCase
	pb.UnimplementedRecordsServiceServer
}

func New(log *zap.Logger, recordsUseCase RecordsUseCase) *getterHandlers {
	return &getterHandlers{
		log:            log,
		recordsUseCase: recordsUseCase,
	}
}

func (g *getterHandlers) GetRandomRecords(ctx context.Context, req *pb.GetRandomRecordsRequest) (*pb.GetRandomRecordsResponse, error) {
	records, err := g.recordsUseCase.Records(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetRandomRecordsResponse{Records: records}, nil
}

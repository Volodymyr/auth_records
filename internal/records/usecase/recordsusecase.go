package usecase

import (
	pb "auth_records/pkg/records_grpc/v1"
	"context"

	"go.uber.org/zap"
)

type RecordsRepository interface {
	Records(ctx context.Context) ([]*pb.Record, error)
}

type usecase struct {
	log  *zap.Logger
	repo RecordsRepository
}

func New(log *zap.Logger, repo RecordsRepository) *usecase {
	return &usecase{
		log,
		repo,
	}
}

func (uc *usecase) Records(ctx context.Context) ([]*pb.Record, error) {
	return uc.repo.Records(ctx)
}

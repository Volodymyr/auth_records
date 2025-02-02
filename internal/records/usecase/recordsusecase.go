package usecase

import (
	pb "auth_records/pkg/records_grpc/v1"
	"context"

	"go.uber.org/zap"
)

type recordsRepository interface {
	Records(ctx context.Context) ([]*pb.Record, error)
}

type usecase struct {
	log  *zap.Logger
	repo recordsRepository
}

func New(log *zap.Logger, repo recordsRepository) *usecase {
	return &usecase{
		log,
		repo,
	}
}

func (uc *usecase) Records(ctx context.Context) ([]*pb.Record, error) {
	return uc.repo.Records(ctx)
}

package storage

import (
	pb "auth_records/pkg/records_grpc/v1"
	"context"
)

type RecordsProvider interface {
	Records(ctx context.Context) ([]*pb.Record, error)
}

type recordsRepository struct {
	client RecordsProvider
}

func NewRecordsRepository(client RecordsProvider) *recordsRepository {
	return &recordsRepository{
		client: client,
	}
}

func (r *recordsRepository) Records(ctx context.Context) ([]*pb.Record, error) {
	records, err := r.client.Records(ctx)
	if err != nil {
		return nil, err
	}

	return records, nil
}

package pgprovider

import (
	pb "auth_records/pkg/records_grpc/v1"
	"context"
)

const queryUserByEmail = "SELECT id, title, content FROM records ORDER BY RANDOM() LIMIT 10"

func (p *pgProvider) Records(ctx context.Context) ([]*pb.Record, error) {
	rows, err := p.db.QueryContext(ctx, queryUserByEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	defer rows.Close()

	var records []*pb.Record
	for rows.Next() {
		var record pb.Record
		if err := rows.Scan(&record.Id, &record.Title, &record.Content); err != nil {
			return nil, err
		}
		records = append(records, &record)
	}

	return records, nil
}

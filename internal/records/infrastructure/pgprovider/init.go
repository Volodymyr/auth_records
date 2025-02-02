package pgprovider

import (
	"context"
	"fmt"
	"time"
)

const (
	queryCheckRecords = "SELECT COUNT(*) FROM records"
	queryInsertRecord = "INSERT INTO records (title, content) VALUES ($1, $2)"
)

// InitRecords checks if there are any records in the 'records' table
// and creates 20 random records if the table is empty.
func (p *pgProvider) InitRecords() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := p.db.QueryRowContext(ctx, queryCheckRecords).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking records: %w", err)
	}

	if count > 0 {
		return nil
	}

	return p.seedRandomRecords(ctx, 20)
}

// seedRandomRecords inserts 20 random records into the 'records' table.
func (p *pgProvider) seedRandomRecords(ctx context.Context, numRecords int) error {
	for i := 0; i < numRecords; i++ {
		title := fmt.Sprintf("Record %d", i+1)
		content := fmt.Sprintf("Description of random record %d", i+1)

		_, err := p.db.ExecContext(ctx, queryInsertRecord, title, content)
		if err != nil {
			return fmt.Errorf("error inserting record: %w", err)
		}
	}

	return nil
}

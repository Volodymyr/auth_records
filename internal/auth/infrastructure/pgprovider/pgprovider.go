package pgprovider

import (
	"auth_records/pkg/pgclient"
	"database/sql"
)

const ServiceName = "Postgres Provider"

type pgProvider struct {
	client *pgclient.Client
	db     *sql.DB
}

func NewProvider(dsn string) *pgProvider {
	client := pgclient.New(dsn)

	return &pgProvider{
		client: client,
	}
}

func (p *pgProvider) Connect() error {
	err := p.client.Connect()
	if err != nil {
		return err
	}

	p.db = p.client.Conn()

	return nil
}

func (p *pgProvider) ServiceName() string {
	return ServiceName
}

func (p *pgProvider) Shutdown() error {
	return p.client.Shutdown()
}

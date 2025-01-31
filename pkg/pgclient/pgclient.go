package pgclient

import (
	"database/sql"
)

type Client struct {
	dsn  string
	conn *sql.DB
}

func New(dsn string) *Client {
	return &Client{
		dsn: dsn,
	}
}

func (c *Client) Connect() error {
	conn, err := sql.Open("pgx", c.dsn)
	if err != nil {
		return err
	}

	err = conn.Ping()
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

func (c *Client) Conn() *sql.DB {
	return c.conn
}

func (c *Client) Shutdown() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

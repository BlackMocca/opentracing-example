package psql

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	pg "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/qustavo/sqlhooks/v2"
)

type Client struct {
	db            *sqlx.DB
	connectionURI string
}

func NewPsqlConnection(connectionStr string) (*Client, error) {
	addr, err := pg.ParseURL(connectionStr)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect("postgres", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		db:            db,
		connectionURI: connectionStr,
	}, nil
}

func NewPsqlWithTracingConnection(connectionStr string, tracing opentracing.Tracer) (*Client, error) {
	addr, err := pg.ParseURL(connectionStr)
	if err != nil {
		return nil, err
	}

	sql.Register("ot_pg", sqlhooks.Wrap(&pg.Driver{}, NewTracingHook(tracing)))

	db, err := sqlx.Connect(opentracing_driver, addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		db:            db,
		connectionURI: connectionStr,
	}, nil
}

func (c *Client) GetClient() *sqlx.DB {
	return c.db
}

func (c *Client) GetConnectionURI() string {
	return c.connectionURI
}

func (c *Client) SetDB(db *sqlx.DB) {
	c.db = db
}

func (c *Client) IsConnect() bool {
	if err := c.db.Ping(); err == nil {
		return true
	}
	return false
}

func (c *Client) Reconnect() error {
	if c.IsConnect() {
		return nil
	}

	client, err := NewPsqlConnection(c.connectionURI)
	if err != nil {
		return err
	}
	c.db = client.GetClient()
	return nil
}

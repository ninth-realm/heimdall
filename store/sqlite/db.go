package sqlite

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/ninth-realm/heimdall/store"
	_ "modernc.org/sqlite"
)

var _ store.Repository = (*DB)(nil)

type DB struct {
	Conn          *sqlx.DB
	UUIDGenerator UUIDGenerator
}

type UUIDGenerator interface {
	GenerateUUID() uuid.UUID
}

func NewDB(dsn string) (DB, error) {
	conn, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return DB{}, err
	}

	if err := conn.Ping(); err != nil {
		return DB{}, err
	}

	return DB{Conn: conn}, nil
}

func (db DB) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return db.Conn.BeginTxx(ctx, nil)
}

type Querier interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (db DB) querier(txn *sqlx.Tx) Querier {
	if txn != nil {
		return txn
	}

	return db.Conn
}

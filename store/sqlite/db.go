package sqlite

import (
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

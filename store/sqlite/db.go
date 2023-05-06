package sqlite

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type DB struct {
	Conn *sqlx.DB
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

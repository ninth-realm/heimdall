package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattmeyers/level"
	"github.com/ninth-realm/heimdall/auth"
	"github.com/ninth-realm/heimdall/client"
	"github.com/ninth-realm/heimdall/http"
	"github.com/ninth-realm/heimdall/store"
	"github.com/ninth-realm/heimdall/store/sqlite"
	"github.com/ninth-realm/heimdall/user"
	_ "modernc.org/sqlite"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := loadConfig()

	logLevel, err := level.ParseLevel(config.logLevel)
	if err != nil {
		return err
	}

	logger, err := level.NewBasicLogger(logLevel, nil)
	if err != nil {
		return err
	}

	var db sqlite.DB
	switch config.Driver {
	case "mem":
		db, err = getSqliteDB("file::memory:", true)
		db.UUIDGenerator = uuidV4Fn(uuid.NewV4)
	case "sqlite":
		db, err = getSqliteDB(
			fmt.Sprintf("file:%s?mode=rwc", config.SQLite.Path),
			config.runMigrations,
		)
		db.UUIDGenerator = uuidV4Fn(uuid.NewV4)
	default:
		return errors.New("unknown driver")
	}

	if err != nil {
		return err
	}

	logger.Info("Using DB driver: %s", config.Driver)

	return buildServer(config, db).ListenAndServe(fmt.Sprintf(":%d", config.port))
}

func buildServer(config Config, db store.Repository) *http.Server {
	srv := http.NewServer()
	srv.Logger, _ = level.NewBasicLogger(level.Info, nil)
	srv.DisableAuth = config.setupMode
	srv.UserService = user.Service{Repo: db}
	srv.ClientService = client.Service{Repo: db}
	srv.AuthService = auth.Service{Repo: db}

	return srv
}

func getSqliteDB(dsn string, runMigrations bool) (sqlite.DB, error) {
	db, err := sqlite.NewDB(dsn)
	if err != nil {
		return sqlite.DB{}, err
	}

	if runMigrations {
		driver, err := sqlite3.WithInstance(db.Conn.DB, &sqlite3.Config{})
		if err != nil {
			return sqlite.DB{}, err
		}

		m, err := migrate.NewWithDatabaseInstance(
			"file://./db/migrations/sqlite",
			"sqlite", driver)
		if err != nil {
			return sqlite.DB{}, err
		}

		err = m.Up()
		if err != migrate.ErrNoChange && err != nil {
			return sqlite.DB{}, err
		}
	}

	return db, nil
}

type uuidV4Fn func() (uuid.UUID, error)

func (f uuidV4Fn) GenerateUUID() uuid.UUID {
	id, err := f()
	if err != nil {
		panic(err)
	}

	return id
}

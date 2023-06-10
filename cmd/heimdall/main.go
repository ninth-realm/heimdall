package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattmeyers/level"
	"github.com/ninth-realm/heimdall/auth"
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
	flags := initFlags()

	logLevel, err := level.ParseLevel(flags.logLevel)
	if err != nil {
		return err
	}

	logger, err := level.NewBasicLogger(logLevel, nil)
	if err != nil {
		return err
	}

	var db sqlite.DB
	switch flags.storeDriver {
	case "mem":
		db, err = getSqliteDB("file::memory:", true)
	case "sqlite":
		db, err = getSqliteDB("file:db/data/heimdall-dev.db?mode=rwc", flags.runMigrations)
	default:
		return errors.New("unknown driver")
	}

	if err != nil {
		return err
	}
	db.UUIDGenerator = uuidV4Fn(uuid.NewV4)

	logger.Info("Using DB driver: %s", flags.storeDriver)

	return buildServer(db).ListenAndServe(fmt.Sprintf(":%d", flags.port))
}

type flags struct {
	port          int
	logLevel      string
	storeDriver   string
	runMigrations bool
}

func initFlags() flags {
	var fs flags

	flag.IntVar(&fs.port, "port", 8080, "port to run on")
	flag.StringVar(&fs.storeDriver, "driver", "mem", "Database driver: mem, sqlite")
	flag.BoolVar(&fs.runMigrations, "migrate", false, "Prevent migrating db. Ignored for mem driver.")
	flag.StringVar(&fs.logLevel, "log-level", "info", "Min log level: debug, info, warn, error, fatal")

	flag.Parse()

	return fs
}

func buildServer(db store.Repository) *http.Server {
	srv := http.NewServer()
	srv.Logger, _ = level.NewBasicLogger(level.Info, nil)
	srv.UserService = user.Service{Repo: db}
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

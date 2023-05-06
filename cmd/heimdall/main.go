package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mattmeyers/level"
	"github.com/ninth-realm/heimdall/http"
	"github.com/ninth-realm/heimdall/store/sqlite"
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

	logger.Info("Using DB driver: %s", flags.storeDriver)

	_ = db

	return buildServer().ListenAndServe(fmt.Sprintf(":%d", flags.port))
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

func buildServer() *http.Server {
	srv := http.NewServer()
	srv.Logger, _ = level.NewBasicLogger(level.Info, nil)

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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattmeyers/level"
	"github.com/ninth-realm/heimdall/http"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

var (
	flagPort int
)

func run() error {
	initFlags()
	return buildServer().ListenAndServe(fmt.Sprintf(":%d", flagPort))
}

func initFlags() {
	flag.IntVar(&flagPort, "port", 8080, "port to run on")
	flag.Parse()
}

func buildServer() *http.Server {
	srv := http.NewServer()
	srv.Logger, _ = level.NewBasicLogger(level.Info, nil)

	return srv
}
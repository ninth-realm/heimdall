package main

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	port          int
	logLevel      string
	runMigrations bool
	configPath    string
	setupMode     bool

	Driver string        `json:"driver"`
	SQLite *SQLiteConfig `json:"sqlite"`
}

type SQLiteConfig struct {
	Path string `json:"path"`
}

func loadConfig() (Config, error) {
	config := initFlags()

	return readConfig(config, config.configPath)
}

func initFlags() Config {
	var config Config

	flag.IntVar(&config.port, "port", 8080, "Port to run on.")
	flag.BoolVar(&config.runMigrations, "migrate", false, "Run db migrations. Ignored for mem driver.")
	flag.StringVar(&config.logLevel, "log-level", "info", "Min log level: debug, info, warn, error, fatal")
	flag.StringVar(&config.configPath, "config", "./config.json", "Path to the config file.")
	flag.BoolVar(
		&config.setupMode,
		"setup-mode",
		false,
		"Removes security checks. This should only be used for initial setup when the service is not exposed to the internet.",
	)

	flag.Parse()

	return config
}

func readConfig(config Config, path string) (Config, error) {
	fileContents, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(fileContents, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

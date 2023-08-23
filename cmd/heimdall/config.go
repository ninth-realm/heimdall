package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Driver string        `json:"driver"`
	SQLite *SQLiteConfig `json:"sqlite"`
}

type SQLiteConfig struct {
	Path string `json:"path"`
}

func readConfig(path string) (Config, error) {
	fileContents, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(fileContents, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

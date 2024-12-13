package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

// FromEnv parses config from env vars into a struct
func FromEnv() (Conf, error) {
	var config Conf
	err := env.Parse(&config)
	if err != nil {
		return config, fmt.Errorf("parse env: %w", err)
	}

	return config, nil
}

type Conf struct {
	DB DB
}

type DB struct {
	Host     string `env:"DB_HOST,notEmpty"`
	Database string `env:"DB_DATABASE,notEmpty"`
	Username string `env:"DB_USERNAME,notEmpty"`
	Password string `env:"DB_PASSWORD,notEmpty"`
}

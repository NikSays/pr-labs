package config

import (
	"fmt"

	"github.com/caarlos0/env/v10"
)

func FromEnv() (Conf, error) {
	var config Conf
	err := env.Parse(&config)
	if err != nil {
		return config, fmt.Errorf("parse env: %w", err)
	}

	return config, nil
}

type Conf struct {
	DB     DB
	Upload Upload
	TCP    TCP
}

type DB struct {
	Host     string `env:"DB_HOST,notEmpty"`
	Database string `env:"DB_DATABASE,notEmpty"`
	Username string `env:"DB_USERNAME,notEmpty"`
	Password string `env:"DB_PASSWORD,notEmpty"`
}
type Upload struct {
	Directory string `env:"UP_DIRECTORY,notEmpty"`
}
type TCP struct {
	FilePath string `env:"TCP_FILEPATH,notEmpty"`
}

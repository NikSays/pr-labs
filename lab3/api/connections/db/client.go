package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Use `go generate ./...` to generate go files from SQL
//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate

type Config struct {
	Host     string
	Database string
	Username string
	Password string
}

// NewClient returns a pooled connection to the database
func NewClient(conf Config) (*Queries, error) {
	connStr := fmt.Sprintf("host=%s database=%s user=%s password=%s", conf.Host, conf.Database, conf.Username, conf.Password)
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return New(conn), nil
}

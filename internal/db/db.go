// Package db provides a counter store backed by libSQL.
package db

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Counter struct {
	Name  string
	Count int64
}

// DB is the counter store the libSQL backend implements.
type DB interface {
	GetCounter(ctx context.Context, name string) (Counter, error)
	UpsertCounter(ctx context.Context, name string, count int64) error
	Close() error
}

// Config describes the libSQL database to connect to.
type Config struct {
	URL   string
	Token string
}

// Open connects to the libSQL database selected by cfg.
func Open(cfg Config) (DB, error) {
	return openLibSQL(cfg.URL, cfg.Token)
}

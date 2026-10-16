package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tursodatabase/libsql-client-go/libsql"
)

type libsqlDB struct {
	*sql.DB
}

func openLibSQL(dbURL, authToken string) (DB, error) {
	if dbURL == "" {
		return nil, errors.New("LIBSQL_URL is not set")
	}
	if authToken == "" {
		return nil, errors.New("LIBSQL_AUTH_TOKEN is not set")
	}

	connector, err := libsql.NewConnector(dbURL, libsql.WithAuthToken(authToken))
	if err != nil {
		return nil, fmt.Errorf("libSQL connector: %w", err)
	}

	conn := sql.OpenDB(connector)
	conn.SetMaxOpenConns(1)
	return &libsqlDB{conn}, nil
}

func (db *libsqlDB) GetCounter(ctx context.Context, name string) (Counter, error) {
	var c Counter
	err := db.QueryRowContext(ctx, `SELECT name, count FROM counter WHERE name = ?`, name).Scan(&c.Name, &c.Count)
	if errors.Is(err, sql.ErrNoRows) {
		return Counter{Name: name}, ErrNotFound
	}
	return c, err
}

func (db *libsqlDB) UpsertCounter(ctx context.Context, name string, count int64) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO counter (name, count) VALUES (?, ?)
	ON CONFLICT(name) DO UPDATE SET count = excluded.count
	`, name, count)
	return err
}

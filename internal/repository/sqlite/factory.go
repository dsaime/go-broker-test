package sqlite

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DSN string
}

type RepositoryFactory struct {
	db *sqlx.DB
}

func (f *RepositoryFactory) DB() *sqlx.DB {
	return f.db
}

func InitRepositoryFactory(config Config) (*RepositoryFactory, error) {
	db, err := sqlx.Connect("sqlite3", config.DSN)
	if err != nil {
		return nil, err
	}
	// Test database connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return &RepositoryFactory{
		db: db,
	}, nil
}

func (f *RepositoryFactory) Close() error {
	return f.db.Close()
}

func (f *RepositoryFactory) Healthcheck(ctx context.Context) error {
	if err := f.db.PingContext(ctx); err != nil {
		return fmt.Errorf("db.PingContext: %w", err)
	}

	return nil
}

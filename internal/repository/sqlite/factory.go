package sqlite

import (
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

func InitRepositoryFactory(config Config) (*RepositoryFactory, error) {
	db, err := sqlx.Connect("sqlite3", config.DSN)
	if err != nil {
		return nil, err
	}
	rf := &RepositoryFactory{
		db: db,
	}
	// Test database connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return rf, nil
}

func (f *RepositoryFactory) Close() error {
	return f.db.Close()
}

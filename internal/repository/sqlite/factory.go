package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	DSN string
}

type RepositoryFactory struct {
	db *sqlx.DB
}

func InitRepositoryFactory(config Config) (*RepositoryFactory, error) {
	db, err := sqlx.Connect("sqlite", config.DSN)
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

func (r *RepositoryFactory) Close() error {
	return r.db.Close()
}

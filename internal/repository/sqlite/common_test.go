package sqlite

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func initDefaultRepoFactory(t *testing.T) *RepositoryFactory {
	repositoryFactory, err := InitRepositoryFactory(Config{
		DSN: ":memory:",
	})
	require.NoError(t, err)
	require.NotNil(t, repositoryFactory)

	// Применить миграции
	err = upFromGooseMigrations(repositoryFactory.db, migrationsDir)
	require.NoError(t, err)

	return repositoryFactory
}

const migrationsDir = "../../../migrations/sqlite"

func upFromGooseMigrations(db *sqlx.DB, dir string) error {
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		var file []byte
		if file, err = os.ReadFile(path); err != nil {
			return err
		}

		// Разделяем содержимое файла по разделителю "-- +goose Down"
		parts := strings.SplitN(string(file), "-- +goose Down", 2)
		upPart := strings.TrimSpace(parts[0]) // Берем только часть до разделителя

		if upPart == "" {
			continue // Пропускаем пустые миграции
		}

		if _, err = db.Exec(upPart); err != nil {
			return err
		}
	}

	return nil
}

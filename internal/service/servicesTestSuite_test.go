package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gitlab.com/digineat/go-broker-test/internal/model"
	"gitlab.com/digineat/go-broker-test/internal/repository/sqlite"
)

type servicesTestSuite struct {
	suite.Suite
	factory *sqlite.RepositoryFactory
	rr      struct {
		trades model.TradesRepository
	}
	ss struct {
		trades *Trades
	}
}

func Test_ServicesTestSuite(t *testing.T) {
	suite.Run(t, new(servicesTestSuite))
}

func (suite *servicesTestSuite) SetupSubTest() {

	// Инициализация SQLiteMemory
	suite.factory = initDefaultRepoFactory(suite.T())

	// Инициализация репозиториев
	suite.rr.trades = suite.factory.NewTradesRepository()

	// Создание сервисов
	suite.ss.trades = &Trades{
		TradesRepo: suite.rr.trades,
	}
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *servicesTestSuite) TearDownSubTest() {
	err := suite.factory.Close()
	suite.Require().NoError(err)
}

func initDefaultRepoFactory(t *testing.T) *sqlite.RepositoryFactory {
	repositoryFactory, err := sqlite.InitRepositoryFactory(sqlite.Config{
		DSN: ":memory:",
	})
	require.NoError(t, err)
	require.NotNil(t, repositoryFactory)

	// Применить миграции
	err = upFromGooseMigrations(repositoryFactory.DB(), migrationsDir)
	require.NoError(t, err)

	return repositoryFactory
}

const migrationsDir = "../../migrations/sqlite"

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
		if !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		var file []byte
		if file, err = os.ReadFile(path); err != nil {
			return err
		}

		if _, err = db.Exec(string(file)); err != nil {
			return err
		}
	}

	return nil
}

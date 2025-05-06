package sqlite

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/digineat/go-broker-test/internal/model"
	"gitlab.com/digineat/go-broker-test/internal/model/repository_tests"
)

func TestChatsRepository(t *testing.T) {
	repository_tests.TradesRepositoryTests(t, func() model.TradesRepository {
		repositoryFactory := initDefaultRepoFactory(t)
		repo := repositoryFactory.NewTradesRepository()
		require.NotNil(t, repo)
		return repo
	})
}

package repository_tests

import (
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/digineat/go-broker-test/internal/model"
)

func saveTrade(t *testing.T, repo model.TradesRepository, trade model.Trade) model.Trade {
	err := repo.Save(trade)
	require.NoError(t, err)
	return trade
}

func newRndTrade(t *testing.T, repo model.TradesRepository) model.Trade {
	return saveTrade(t, repo, model.Trade{
		ID:      uuid.NewString(),
		Account: uuid.NewString(),
		Symbol:  "RANDOM",
		Volume:  float64(rand.Int63()),
		Open:    float64(rand.Int63()),
		Close:   float64(rand.Int63()),
		Side:    []string{model.TradeSideByy, model.TradeSideSell}[rand.Intn(2)],
	})
}

func TradesRepositoryTests(t *testing.T, newRepository func() model.TradesRepository) {
	t.Run("List", func(t *testing.T) {
		t.Run("пустой репозиторий вернет пустой список", func(t *testing.T) {
			r := newRepository()
			trades, err := r.List(model.TradeListFilter{})
			assert.NoError(t, err)
			assert.Empty(t, trades)
		})
		t.Run("без фильтра вернутся все записи", func(t *testing.T) {
			r := newRepository()
			const expectedCount = 10
			for range expectedCount {
				newRndTrade(t, r)
			}
			trades, err := r.List(model.TradeListFilter{})
			assert.NoError(t, err)
			assert.Len(t, trades, expectedCount)
			// TODO: equality check
		})
		t.Run("с фильтром по account вернутся записи с соответствующим account", func(t *testing.T) {
			r := newRepository()
			// Случайные
			for range 10 {
				newRndTrade(t, r)
			}

			// Искомые
			account := uuid.NewString()
			expectedTrades := make([]model.Trade, 5)
			for i := range expectedTrades {
				expectedTrades[i] = saveTrade(t, r, model.Trade{
					ID:      uuid.NewString(),
					Account: account,
				})
			}

			// Получить
			trades, err := r.List(model.TradeListFilter{
				Account: account,
			})
			assert.NoError(t, err)
			assert.Len(t, trades, len(expectedTrades))
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("нельзя сохранять без ID", func(t *testing.T) {
			r := newRepository()
			err := r.Save(model.Trade{ID: ""})
			assert.Error(t, err)
		})
		t.Run("остальные поля могут быть пустыми", func(t *testing.T) {
			r := newRepository()
			err := r.Save(model.Trade{ID: "qwerty"})
			assert.NoError(t, err)
		})
		t.Run("сохраненный можно запросить и он полностью соответствует сохраняемому", func(t *testing.T) {
			r := newRepository()
			// Сохранить
			savedTrade := newRndTrade(t, r)
			tradeFromRepo, err := r.List(model.TradeListFilter{})
			require.NoError(t, err)
			require.Len(t, tradeFromRepo, 1)
			// Сравнить
			assert.Equal(t, savedTrade, tradeFromRepo[0])
		})
		t.Run("перезапись при сохранении существующего ID", func(t *testing.T) {
			r := newRepository()

			savedTrade := newRndTrade(t, r)
			for range 10 {
				saveTrade(t, r, savedTrade)
			}

			tradeFromRepo, err := r.List(model.TradeListFilter{})
			require.NoError(t, err)
			require.Len(t, tradeFromRepo, 1)
		})
	})
}

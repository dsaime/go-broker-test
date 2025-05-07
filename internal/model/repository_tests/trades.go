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
		ID:        uuid.NewString(),
		Account:   uuid.NewString(),
		Symbol:    "RANDOM",
		Volume:    float64(rand.Int63()),
		Open:      float64(rand.Int63()),
		Close:     float64(rand.Int63()),
		Side:      []string{model.TradeSideByy, model.TradeSideSell}[rand.Intn(2)],
		WorkerID:  uuid.NewString(),
		JobStatus: model.TradeJobStatusNew,
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
			savedTrades := make([]model.Trade, 10)
			for i := range savedTrades {
				savedTrades[i] = newRndTrade(t, r)
			}
			tradesFromRepo, err := r.List(model.TradeListFilter{})
			assert.NoError(t, err)
			require.Len(t, tradesFromRepo, len(savedTrades))
			// Сравнить
			for i := range tradesFromRepo {
				assert.Equal(t, savedTrades[i], tradesFromRepo[i])
			}
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
	t.Run("UpdateNobodyAndGet", func(t *testing.T) {
		t.Run("пустой репозиторий вернет пустой список", func(t *testing.T) {
			r := newRepository()
			trades, err := r.UpdateNobodyAndGet(model.UpdateNobodyAndGetInput{
				NewWorkerID:  "123",
				NewJobStatus: model.TradeJobStatusNew,
			})
			assert.NoError(t, err)
			assert.Empty(t, trades)
		})
		t.Run("параметры можно передавать пустым", func(t *testing.T) {
			r := newRepository()
			trades, err := r.UpdateNobodyAndGet(model.UpdateNobodyAndGetInput{
				NewWorkerID:  "",
				NewJobStatus: "",
			})
			assert.NoError(t, err)
			assert.Empty(t, trades)
		})
		t.Run("вернутся записи у которых не было WorkerID", func(t *testing.T) {
			r := newRepository()
			// Случайные
			for range 10 {
				newRndTrade(t, r)
			}
			// Не связаны ни с одним worker
			const nobodyTradesSavedCount = 5
			for range nobodyTradesSavedCount {
				saveTrade(t, r, model.Trade{
					ID:        uuid.NewString(),
					WorkerID:  "",
					JobStatus: "",
				})
			}
			trades, err := r.UpdateNobodyAndGet(model.UpdateNobodyAndGetInput{})
			assert.NoError(t, err)
			assert.Len(t, trades, nobodyTradesSavedCount)
		})
		t.Run("после выполнения UpdateNobodyAndGet всем записям без WorkerID будет установлен NewWorkerID и NewJobStatus", func(t *testing.T) {
			r := newRepository()
			const nobodyTradesSavedCount = 10
			for range nobodyTradesSavedCount {
				saveTrade(t, r, model.Trade{ID: uuid.NewString()})
			}
			input := model.UpdateNobodyAndGetInput{
				NewWorkerID:  uuid.NewString(),
				NewJobStatus: model.TradeJobStatusProcessing,
			}
			trades, err := r.UpdateNobodyAndGet(input)
			assert.NoError(t, err)
			require.Len(t, trades, nobodyTradesSavedCount)
			// Проверить, что все записи получили NewWorkerID и NewJobStatus
			for _, trade := range trades {
				assert.Equal(t, input.NewWorkerID, trade.WorkerID)
				assert.Equal(t, input.NewJobStatus, trade.JobStatus)
			}
		})
	})
}

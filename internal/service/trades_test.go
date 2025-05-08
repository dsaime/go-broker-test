package service

import (
	"math/rand"

	"github.com/google/uuid"

	"gitlab.com/digineat/go-broker-test/internal/model"
)

func newRndEnqueueTradeInput() EnqueueTradeInput {
	return EnqueueTradeInput{
		Account: uuid.NewString(),
		Symbol:  "RANDOM",
		Volume:  float64(rand.Int63()),
		Open:    float64(rand.Int63()),
		Close:   float64(rand.Int63()),
		Side:    []string{model.TradeSideByy, model.TradeSideSell}[rand.Intn(2)],
	}
}

func (suite *servicesTestSuite) Test_Trades_CalculateProfitOnNobodyTrades() {
	suite.Run("трейда должна содержать валидные значения", func() {
		input := CalculateProfitOnNobodyTradesInput{
			WorkerID: "",
		}
		cltrades, err := suite.ss.trades.CalculateProfitOnNobodyTrades(input)
		suite.ErrorIs(err, ErrRequiredWorkerID)
		suite.Empty(cltrades)
	})
	suite.Run("обрабатываются только трейды без воркера", func() {
		// Трейды со случайными воркерами
		for range 3 {
			suite.saveTrades(newRndTrade())
		}
		// Трейды без воркеров
		const tradesWithWorkersCount = 10
		for range tradesWithWorkersCount {
			trade := newRndTrade()
			trade.WorkerID = ""
			suite.saveTrades(trade)
		}

		// Запустить обработку
		input := CalculateProfitOnNobodyTradesInput{
			WorkerID: uuid.NewString(),
		}
		cltrades, err := suite.ss.trades.CalculateProfitOnNobodyTrades(input)
		suite.Require().NoError(err)
		suite.Require().NotEmpty(cltrades)

		// Получить трейды
		trades, err := suite.rr.trades.List(model.TradeListFilter{})
		suite.Require().NoError(err)
		suite.Require().True(len(trades) > tradesWithWorkersCount)
		// Количество обработанных равно количество трейдов без воркеров (ранее)
		var matched int
		for _, trade := range trades {
			if trade.WorkerID == input.WorkerID && trade.JobStatus == model.TradeJobStatusDone {
				matched++
			}
		}
		suite.Equal(tradesWithWorkersCount, matched)
	})
	suite.Run("все трейды должны получить новый статус и wokerID", func() {
		// Сохранить трейд без воркера и с начальным статусом
		const savedTradesCount = 10
		for range savedTradesCount {
			trade := newRndTrade()
			trade.WorkerID = ""
			trade.JobStatus = model.TradeJobStatusNew
			suite.saveTrades(trade)
		}

		// Запустить обработку
		input := CalculateProfitOnNobodyTradesInput{
			WorkerID: uuid.NewString(),
		}
		cltrades, err := suite.ss.trades.CalculateProfitOnNobodyTrades(input)
		suite.Require().NoError(err)
		suite.Require().NotEmpty(cltrades)

		// Получить трейды
		trades, err := suite.rr.trades.List(model.TradeListFilter{})
		suite.Require().NoError(err)
		suite.Require().Len(trades, savedTradesCount)
		for _, trade := range trades {
			// WorkerID установлен из параметров
			suite.Equal(input.WorkerID, trade.WorkerID)
			// Статус стал равен Done
			suite.Equal(model.TradeJobStatusDone, trade.JobStatus)
		}
	})
	suite.Run("трейды получают новое значение Profit", func() {
		testCases := []struct {
			trade          model.Trade
			expectedProfit float64
		}{
			{
				trade: model.Trade{
					ID:     uuid.NewString(),
					Volume: 40,
					Open:   20,
					Close:  300,
					Side:   model.TradeSideSell,
				},
				expectedProfit: (300 - 20) * 40 * 100000 * -1,
			},
			{
				trade: model.Trade{
					ID:     uuid.NewString(),
					Volume: 20,
					Open:   200,
					Close:  30,
					Side:   model.TradeSideSell,
				},
				expectedProfit: (30 - 200) * 20 * 100000 * -1,
			},
			{
				trade: model.Trade{
					ID:     uuid.NewString(),
					Volume: 20,
					Open:   200,
					Close:  30,
					Side:   model.TradeSideByy,
				},
				expectedProfit: (30 - 200) * 20 * 100000,
			},
			{
				trade: model.Trade{
					ID:     uuid.NewString(),
					Volume: 40,
					Open:   20,
					Close:  300,
					Side:   model.TradeSideByy,
				},
				expectedProfit: (300 - 20) * 40 * 100000,
			},
		}
		for i := range testCases {
			suite.saveTrades(testCases[i].trade)
		}

		// Запустить обработку
		input := CalculateProfitOnNobodyTradesInput{
			WorkerID: uuid.NewString(),
		}
		cltrades, err := suite.ss.trades.CalculateProfitOnNobodyTrades(input)
		suite.Require().NoError(err)
		suite.Require().NotEmpty(cltrades)

		// Получить трейды
		trades, err := suite.rr.trades.List(model.TradeListFilter{})
		suite.Require().NoError(err)
		suite.Require().Len(trades, len(testCases))
		for i := range trades {
			suite.Equal(testCases[i].expectedProfit, trades[i].Profit)
		}
	})
}

func (suite *servicesTestSuite) Test_Trades_EnqueueTrade() {
	suite.Run("трейда должен содержать валидные значения", func() {
		validInput := newRndEnqueueTradeInput()
		input := EnqueueTradeInput{
			Account: validInput.Account,
			Symbol:  validInput.Symbol,
			Volume:  validInput.Volume,
			Open:    validInput.Open,
			Close:   validInput.Close,
			Side:    validInput.Side,
		}

		input.Account = ""
		suite.ErrorIs(suite.ss.trades.EnqueueTrade(input), model.ErrInvalidAccount)
		input.Account = validInput.Account

		input.Symbol = ""
		suite.ErrorIs(suite.ss.trades.EnqueueTrade(input), model.ErrInvalidSymbol)
		input.Symbol = validInput.Symbol

		input.Volume = -1
		suite.ErrorIs(suite.ss.trades.EnqueueTrade(input), model.ErrInvalidVolume)
		input.Volume = validInput.Volume

		input.Open = -1
		suite.ErrorIs(suite.ss.trades.EnqueueTrade(input), model.ErrInvalidOpen)
		input.Open = validInput.Open

		input.Close = -1
		suite.ErrorIs(suite.ss.trades.EnqueueTrade(input), model.ErrInvalidClose)
		input.Close = validInput.Close

		input.Side = ""
		suite.ErrorIs(suite.ss.trades.EnqueueTrade(input), model.ErrInvalidSide)
		input.Side = validInput.Side
	})
	suite.Run("трейд будет сохранен с заданными значениями", func() {
		input := newRndEnqueueTradeInput()
		err := suite.ss.trades.EnqueueTrade(input)
		suite.NoError(err)

		trades, err := suite.rr.trades.List(model.TradeListFilter{})
		suite.NoError(err)
		suite.Require().Len(trades, 1)
		// Совпадает с входными параметрами
		suite.Equal(input.Account, trades[0].Account)
		suite.Equal(input.Symbol, trades[0].Symbol)
		suite.Equal(input.Volume, trades[0].Volume)
		suite.Equal(input.Open, trades[0].Open)
		suite.Equal(input.Close, trades[0].Close)
		suite.Equal(input.Side, trades[0].Side)
	})
	suite.Run("workerID и JobStatus заполнены начальными значениями", func() {
		input := newRndEnqueueTradeInput()
		err := suite.ss.trades.EnqueueTrade(input)
		suite.NoError(err)

		trades, err := suite.rr.trades.List(model.TradeListFilter{})
		suite.NoError(err)
		suite.Require().Len(trades, 1)
		// Воркер не задан
		suite.Empty(trades[0].WorkerID)
		// Статус Новый
		suite.Equal(model.TradeJobStatusNew, trades[0].JobStatus)
	})
}

func newRndTrade() model.Trade {
	return model.Trade{
		ID:        uuid.NewString(),
		Account:   uuid.NewString(),
		Symbol:    "RANDOM",
		Volume:    float64(rand.Int31()),
		Open:      float64(rand.Int31()),
		Close:     float64(rand.Int31()),
		Side:      []string{model.TradeSideByy, model.TradeSideSell}[rand.Intn(2)],
		WorkerID:  uuid.NewString(),
		JobStatus: model.TradeJobStatusNew,
	}
}

func (suite *servicesTestSuite) saveTrades(trades ...model.Trade) []model.Trade {
	err := suite.rr.trades.SaveAll(trades)
	suite.Require().NoError(err)
	return trades
}

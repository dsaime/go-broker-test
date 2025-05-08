package service

import (
	"github.com/google/uuid"

	"gitlab.com/digineat/go-broker-test/internal/model"
)

type Trades struct {
	TradesRepo model.TradesRepository
}

type EnqueueTradeInput struct {
	Account string
	Symbol  string
	Volume  float64
	Open    float64
	Close   float64
	Side    string
}

func (t *Trades) EnqueueTrade(in EnqueueTradeInput) error {
	newTrade := model.Trade{
		ID:        uuid.NewString(),
		Account:   in.Account,
		Symbol:    in.Symbol,
		Volume:    in.Volume,
		Open:      in.Open,
		Close:     in.Close,
		Side:      in.Side,
		WorkerID:  "",
		JobStatus: model.TradeJobStatusNew,
		Profit:    0,
	}
	if err := newTrade.Validate(); err != nil {
		return err
	}

	if err := t.TradesRepo.Save(newTrade); err != nil {
		return err
	}

	return nil
}

type CalculateProfitOnNobodyTradesInput struct {
	WorkerID string
}

func (t *Trades) CalculateProfitOnNobodyTrades(in CalculateProfitOnNobodyTradesInput) ([]model.Trade, error) {
	if in.WorkerID == "" {
		return nil, ErrRequiredWorkerID
	}

	// "Ничейным задачам" назначить worker и сразу установить новый статус
	updateIn := model.UpdateNobodyAndGetInput{
		NewWorkerID:  in.WorkerID,
		NewJobStatus: model.TradeJobStatusProcessing,
	}
	trades, err := t.TradesRepo.UpdateNobodyAndGet(updateIn)
	if err != nil {
		return nil, err
	}

	for i, trade := range trades {
		// Посчитать профит трейда
		lot := 100000.0
		profit := (trade.Close - trade.Open) * trade.Volume * lot
		if trade.Side == model.TradeSideSell {
			profit = -profit
		}
		// Обновить профит и статус
		trade.Profit = profit
		trade.JobStatus = model.TradeJobStatusDone
		trades[i] = trade
	}

	// Записать изменения в БД
	if err = t.TradesRepo.SaveAll(trades); err != nil {
		return nil, err
	}

	return trades, nil
}

type AccountStatisticsInput struct {
	Account string
}

func (t *Trades) AccountStatistics(in AccountStatisticsInput) (model.AccountStats, error) {
	if in.Account == "" {
		return model.AccountStats{}, ErrRequiredAccountID
	}

	trades, err := t.TradesRepo.List(model.TradeListFilter{Account: in.Account})
	if err != nil {
		return model.AccountStats{}, err
	}

	stats := model.AccountStats{
		Account: in.Account,
		Trades:  len(trades),
		Profit:  0,
	}
	for _, trade := range trades {
		stats.Profit += trade.Profit
	}

	return stats, nil
}

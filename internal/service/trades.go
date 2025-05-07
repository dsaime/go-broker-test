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

func (t *Trades) CalculateProfitOnNobodyTrades(in CalculateProfitOnNobodyTradesInput) error {
	if in.WorkerID == "" {
		return ErrRequiredWorkerID
	}

	// "Ничейным задачам" назначить worker и сразу установить новый статус
	updateIn := model.UpdateNobodyAndGetInput{
		NewWorkerID:  in.WorkerID,
		NewJobStatus: model.TradeJobStatusProcessing,
	}
	trades, err := t.TradesRepo.UpdateNobodyAndGet(updateIn)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

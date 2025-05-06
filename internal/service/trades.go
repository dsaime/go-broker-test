package service

import "gitlab.com/digineat/go-broker-test/internal/model"

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
		Account: in.Account,
		Symbol:  in.Symbol,
		Volume:  in.Volume,
		Open:    in.Open,
		Close:   in.Close,
		Side:    in.Side,
	}
	if err := newTrade.Validate(); err != nil {
		return err
	}

	if err := t.TradesRepo.Save(newTrade); err != nil {
		return err
	}

	return nil
}

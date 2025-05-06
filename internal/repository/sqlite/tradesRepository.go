package sqlite

import (
	"gitlab.com/digineat/go-broker-test/internal/model"
)

type TradesRepository struct{}

func (t *TradesRepository) Save(trade model.Trade) error {
	//TODO implement me
	panic("implement me")
}

func (t *TradesRepository) List() ([]model.Trade, error) {
	//TODO implement me
	panic("implement me")
}

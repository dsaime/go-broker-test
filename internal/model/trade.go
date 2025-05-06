package model

import (
	"errors"
	"regexp"
)

type Trade struct {
	Account string  // must not be empty
	Symbol  string  // ^[A-Z]{6}$ (e.g. EURUSD)
	Volume  float64 // must be > 0
	Open    float64 // must be > 0
	Close   float64 // must be > 0
	Side    string  //	either "buy" or "sell"
}

const (
	TradeSideByy  = "buy"
	TradeSideSell = "sell"
)

var (
	ErrInvalidAccount = errors.New("некорректное значение Account")
	ErrInvalidSymbol  = errors.New("некорректное значение Symbol")
	ErrInvalidVolume  = errors.New("некорректное значение Volume")
	ErrInvalidOpen    = errors.New("некорректное значение Open")
	ErrInvalidClose   = errors.New("некорректное значение Close")
	ErrInvalidSide    = errors.New("некорректное значение Side")
)

var TradeSymbolRegex = regexp.MustCompile(`^[A-Z]{6}$`)

func (t Trade) Validate() error {
	if t.Account == "" {
		return ErrInvalidAccount
	}

	if !TradeSymbolRegex.MatchString(t.Symbol) {
		return ErrInvalidSymbol
	}
	if t.Volume <= 0 {
		return ErrInvalidVolume
	}
	if t.Open <= 0 {
		return ErrInvalidOpen
	}
	if t.Close <= 0 {
		return ErrInvalidClose
	}
	if t.Side != TradeSideByy && t.Side != TradeSideSell {
		return ErrInvalidSide
	}

	return nil
}

type TradesRepository interface {
	Save(Trade) error
	List() ([]Trade, error)
}

type TradeListFilter struct {
	Account string
}

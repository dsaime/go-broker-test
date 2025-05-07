package model

import (
	"errors"
	"regexp"
	"slices"

	"github.com/google/uuid"
)

type Trade struct {
	ID      string  // must not be empty
	Account string  // must not be empty
	Symbol  string  // ^[A-Z]{6}$ (e.g. EURUSD)
	Volume  float64 // must be > 0
	Open    float64 // must be > 0
	Close   float64 // must be > 0
	Side    string  // either "buy" or "sell"

	WorkerID  string
	JobStatus string
	//	JobStatusUpdatedAt time.Time
	Profit float64
}

const (
	TradeJobStatusNew        = "new"
	TradeJobStatusProcessing = "processing"
	TradeJobStatusFailed     = "failed"
	TradeJobStatusDone       = "done"
)

const (
	TradeSideByy  = "buy"
	TradeSideSell = "sell"
)

var (
	ErrInvalidID        = errors.New("некорректное значение ID")
	ErrInvalidAccount   = errors.New("некорректное значение Account")
	ErrInvalidSymbol    = errors.New("некорректное значение Symbol")
	ErrInvalidVolume    = errors.New("некорректное значение Volume")
	ErrInvalidOpen      = errors.New("некорректное значение Open")
	ErrInvalidClose     = errors.New("некорректное значение Close")
	ErrInvalidSide      = errors.New("некорректное значение Side")
	ErrInvalidJobStatus = errors.New("некорректное значение JobStatus")
)

var TradeSymbolRegex = regexp.MustCompile(`^[A-Z]{6}$`)

func (t Trade) Validate() error {
	if err := uuid.Validate(t.ID); err != nil {
		return errors.Join(ErrInvalidID, err)
	}
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

	if !slices.Contains([]string{
		TradeJobStatusNew,
		TradeJobStatusProcessing,
		TradeJobStatusFailed,
		TradeJobStatusDone,
	}, t.JobStatus) {
		return ErrInvalidJobStatus
	}

	return nil
}

type TradesRepository interface {
	Save(Trade) error
	SaveAll([]Trade) error
	List(TradeListFilter) ([]Trade, error)
	UpdateNobodyAndGet(UpdateNobodyAndGetInput) ([]Trade, error)
}

type TradeListFilter struct {
	Account string
}

type UpdateNobodyAndGetInput struct {
	NewWorkerID  string
	NewJobStatus string
}

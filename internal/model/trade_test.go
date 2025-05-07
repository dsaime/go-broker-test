package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newValidTrade() Trade {
	return Trade{
		ID:        uuid.NewString(),
		Account:   "732cbdf0",
		Symbol:    "AAABBB",
		Volume:    300000.0,
		Open:      200000.0,
		Close:     100000.0,
		Side:      TradeSideByy,
		WorkerID:  "worker123",
		JobStatus: TradeJobStatusNew,
	}
}

func TestTrade_Validate(t *testing.T) {
	t.Run("ID должно быть валидной uuid-строкой", func(t *testing.T) {
		trade := newValidTrade()
		validValues := []string{
			"414adada-8e00-4a7c-9992-64f45cc771a0",
			"1d659b09-04c8-49b3-862d-dc41235f30b9",
			"f5012cc6-b08f-44e0-8a81-96caa22ccd80",
			"4cb5596d-0dea-45f6-b67a-ccb01206ccf8",
		}
		for i := range validValues {
			trade.ID = validValues[i]
			assert.NoError(t, trade.Validate())
		}

		invalidValues := []string{
			"",
			"1",
			"b",
			"123456789012345678901234567890123456789012345678901234567890123456789012345",
		}
		for i := range invalidValues {
			trade.ID = invalidValues[i]
			assert.ErrorIs(t, trade.Validate(), ErrInvalidID)
		}
	})
	t.Run("Account не должно быть пустым", func(t *testing.T) {
		trade := newValidTrade()
		trade.Account = ""
		err := trade.Validate()
		assert.ErrorIs(t, err, ErrInvalidAccount)
	})
	t.Run("С валидным значение Account ошибка не вернется", func(t *testing.T) {
		trade := newValidTrade()
		trade.Account = "asdfghjkl"
		err := trade.Validate()
		assert.NoError(t, err)
	})

	// Symbol
	t.Run("Symbol должно состоять из 6ти символов", func(t *testing.T) {
		trade := newValidTrade()
		trade.Symbol = "1234567"
		err := trade.Validate()
		assert.ErrorIs(t, err, ErrInvalidSymbol)
	})
	t.Run("Symbol должно состоять только из букв верхнего регистра", func(t *testing.T) {
		trade := newValidTrade()
		trade.Symbol = "aaa333"
		err := trade.Validate()
		assert.ErrorIs(t, err, ErrInvalidSymbol)
	})
	t.Run("Symbol должно содержать только A-Z", func(t *testing.T) {
		trade := newValidTrade()
		trade.Symbol = "a ф1@/"
		err := trade.Validate()
		assert.ErrorIs(t, err, ErrInvalidSymbol)
	})
	t.Run("С валидным значение Symbol ошибка не вернется", func(t *testing.T) {
		trade := newValidTrade()
		validSymbols := []string{
			"ABCDEF",
			"GHIJKL",
			"MNOPQR",
			"STUVWX",
			"YYYZZZ",
		}
		for _, symbol := range validSymbols {
			trade.Symbol = symbol
			err := trade.Validate()
			assert.NoError(t, err)
		}
	})

	// Volume
	t.Run("Volume не может быть меньше 1", func(t *testing.T) {
		trade := newValidTrade()
		trade.Volume = -1
		err := trade.Validate()
		assert.ErrorIs(t, err, ErrInvalidVolume)
	})
	t.Run("С валидным значение Volume ошибка не вернется", func(t *testing.T) {
		trade := newValidTrade()
		trade.Volume = 31289
		err := trade.Validate()
		assert.NoError(t, err)
	})

	// Open
	t.Run("Open не может быть меньше 1", func(t *testing.T) {
		trade := newValidTrade()
		trade.Open = -1
		err := trade.Validate()
		assert.ErrorIs(t, err, ErrInvalidOpen)
	})
	t.Run("С валидным значение Open ошибка не вернется", func(t *testing.T) {
		trade := newValidTrade()
		trade.Open = 31289
		err := trade.Validate()
		assert.NoError(t, err)
	})

	// Close
	t.Run("Close не может быть меньше 1", func(t *testing.T) {
		trade := newValidTrade()
		trade.Close = -1
		err := trade.Validate()
		assert.ErrorIs(t, err, ErrInvalidClose)
	})
	t.Run("С валидным значение Close ошибка не вернется", func(t *testing.T) {
		trade := newValidTrade()
		trade.Close = 31289
		err := trade.Validate()
		assert.NoError(t, err)
	})

	// Side
	t.Run("Side может быть равным только определенным значениям", func(t *testing.T) {
		trade := newValidTrade()
		invalidValues := []string{
			"",
			"1",
			"a",
			"BUY",
			"SELL",
		}
		for _, side := range invalidValues {
			trade.Side = side
			assert.ErrorIs(t, trade.Validate(), ErrInvalidSide)
		}

		validValues := []string{
			"buy",
			"sell",
		}
		for _, side := range validValues {
			trade.Side = side
			assert.NoError(t, trade.Validate())
		}
	})

	t.Run("WorkerID может быть любым значением", func(t *testing.T) {
		trade := newValidTrade()
		validValues := []string{
			"",
			"0",
			"1",
			"w",
			"worker",
			"worker123",
			"4cb5596d-0dea-45f6-b67a-ccb01206ccf8",
		}
		for _, workerID := range validValues {
			trade.WorkerID = workerID
			assert.NoError(t, trade.Validate())
		}
	})

	t.Run("JobStatus может быть равным только определенным значениям", func(t *testing.T) {
		trade := newValidTrade()
		invalidValues := []string{
			"",
			"0",
			"unknown",
			"a",
		}
		for _, jobStatus := range invalidValues {
			trade.JobStatus = jobStatus
			assert.ErrorIs(t, trade.Validate(), ErrInvalidJobStatus)
		}

		validValues := []string{
			"new",
			"processing",
			"failed",
			"done",
		}
		for _, jobStatus := range validValues {
			trade.JobStatus = jobStatus
			assert.NoError(t, trade.Validate())
		}
	})
}

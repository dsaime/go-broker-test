package sqlite

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"gitlab.com/digineat/go-broker-test/internal/model"
)

func (f *RepositoryFactory) NewTradesRepository() *TradesRepository {
	return &TradesRepository{
		db: f.db,
	}
}

type TradesRepository struct {
	db *sqlx.DB
}

func (r *TradesRepository) SaveAll(trades []model.Trade) error {
	if len(trades) == 0 {
		return nil
	}
	for _, trade := range trades {
		if trade.ID == "" {
			return errors.New("invalid ID in trade")
		}
	}

	//if _, err := r.db.NamedExec(`
	//	INSERT OR REPLACE INTO trades_q (id, account, symbol, volume, open, close, side, worker_id, job_status)
	//	VALUES (:id, :account, :symbol, :volume, :open, :close, :side, :worker_id, :job_status)
	//`, tradesFromModels(trades)); err != nil {
	//	return err
	//}
	//
	//return nil

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for _, trade := range trades {
		if _, err = tx.NamedExec(`
			INSERT OR REPLACE INTO trades_q (id, account, symbol, volume, open, close, side, worker_id, job_status, profit)
			VALUES (:id, :account, :symbol, :volume, :open, :close, :side, :worker_id, :job_status, :profit)
		`, tradeFromModel(trade)); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *TradesRepository) UpdateNobodyAndGet(in model.UpdateNobodyAndGetInput) ([]model.Trade, error) {
	var trades []dbTrade

	err := r.db.Select(&trades, `
		UPDATE trades_q 
		SET worker_id = $1,
		    job_status = $2
		WHERE worker_id = ''
		RETURNING *
	`, in.NewWorkerID, in.NewJobStatus)
	if err != nil {
		return nil, err
	}

	return tradesToModels(trades), nil
}

func (r *TradesRepository) List(filter model.TradeListFilter) ([]model.Trade, error) {
	var trades []dbTrade

	err := r.db.Select(&trades, `
		SELECT * 
		FROM trades_q
		WHERE ($1 = '' OR account = $1)
	`, filter.Account)
	if err != nil {
		return nil, fmt.Errorf("db.Select: %w", err)
	}

	return tradesToModels(trades), nil
}

func (r *TradesRepository) Save(trade model.Trade) error {
	if trade.ID == "" {
		return fmt.Errorf("invalid ID")
	}

	if _, err := r.db.NamedExec(`
		INSERT OR REPLACE INTO trades_q (id, account, symbol, volume, open, close, side, worker_id, job_status, profit)
		VALUES (:id, :account, :symbol, :volume, :open, :close, :side, :worker_id, :job_status, :profit)
	`, tradeFromModel(trade)); err != nil {
		return err
	}

	return nil
}

type dbTrade struct {
	ID        string  `db:"id"`
	Account   string  `db:"account"`
	Symbol    string  `db:"symbol"`
	Volume    float64 `db:"volume"`
	Open      float64 `db:"open"`
	Close     float64 `db:"close"`
	Side      string  `db:"side"`
	WorkerID  string  `db:"worker_id"`
	JobStatus string  `db:"job_status"`
	Profit    float64 `db:"profit"`
}

func tradeFromModel(modelTrade model.Trade) dbTrade {
	return dbTrade{
		ID:        modelTrade.ID,
		Account:   modelTrade.Account,
		Symbol:    modelTrade.Symbol,
		Volume:    modelTrade.Volume,
		Open:      modelTrade.Open,
		Close:     modelTrade.Close,
		Side:      modelTrade.Side,
		WorkerID:  modelTrade.WorkerID,
		JobStatus: modelTrade.JobStatus,
		Profit:    modelTrade.Profit,
	}
}

func tradeToModel(t dbTrade) model.Trade {
	return model.Trade{
		ID:        t.ID,
		Account:   t.Account,
		Symbol:    t.Symbol,
		Volume:    t.Volume,
		Open:      t.Open,
		Close:     t.Close,
		Side:      t.Side,
		WorkerID:  t.WorkerID,
		JobStatus: t.JobStatus,
		Profit:    t.Profit,
	}
}

// tradesFromModels преобразует слайс model.Trade в слайс *dbTrade
func tradesFromModels(modelTrades []model.Trade) []dbTrade {
	trades := make([]dbTrade, len(modelTrades))
	for i := range modelTrades {
		trades[i] = tradeFromModel(modelTrades[i])
	}
	return trades
}

// tradesToModels преобразует слайс *dbTrade в слайс model.Trade
func tradesToModels(trades []dbTrade) []model.Trade {
	modelTrades := make([]model.Trade, len(trades))
	for i := range trades {
		modelTrades[i] = tradeToModel(trades[i])
	}
	return modelTrades
}

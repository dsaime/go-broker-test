package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slog"

	"gitlab.com/digineat/go-broker-test/internal/repository/sqlite"
	"gitlab.com/digineat/go-broker-test/internal/service"
)

func main() {
	cfg := initConfig()

	deps, closeDependencies, err := initDependencies(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to init dependencies: " + err.Error())
		os.Exit(1)
	}
	defer closeDependencies()

	log.Printf("Worker started with polling interval: %v", cfg.PollInterval)

	// Main worker loop
	err = runMainWorkerLoop(context.Background(), deps.trades, cfg.WorkerID, cfg.PollInterval)
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("runMainWorkerLoop: " + err.Error())
	}
}

func initConfig() (cfg struct {
	DBPath       string
	PollInterval time.Duration
	WorkerID     string
}) {
	// Command line flags
	flag.StringVar(&cfg.DBPath, "db", "data.db", "path to SQLite database")
	flag.DurationVar(&cfg.PollInterval, "poll", 100*time.Millisecond, "polling interval")
	flag.StringVar(&cfg.WorkerID, "worker", uuid.NewString(), "worker id")
	flag.Parse()

	return cfg
}

func runMainWorkerLoop(ctx context.Context, trades *service.Trades, workerID string, pollInterval time.Duration) error {
	for {
		in := service.CalculateProfitOnNobodyTradesInput{WorkerID: workerID}
		if err := trades.CalculateProfitOnNobodyTrades(in); err != nil {
			slog.Error("trades.CalculateProfitOnNobodyTrades: " + err.Error())
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(pollInterval):
		}
	}
}

func initDependencies(dbPath string) (deps struct {
	repositoryFactory *sqlite.RepositoryFactory
	trades            *service.Trades
}, closeDependencies func(), err error) {
	deps.repositoryFactory, err = sqlite.InitRepositoryFactory(sqlite.Config{DSN: dbPath})
	if err != nil {
		return deps, nil, fmt.Errorf("sqlite.InitRepositoryFactory: %w", err)
	}

	deps.trades = &service.Trades{
		TradesRepo: deps.repositoryFactory.NewTradesRepository(),
	}

	return deps, func() {
		if err := deps.repositoryFactory.Close(); err != nil {
			slog.Error("Failed to close repository factory: " + err.Error())
		}
	}, nil
}

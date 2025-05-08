package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slog"

	httpController "gitlab.com/digineat/go-broker-test/internal/controller/http"
	"gitlab.com/digineat/go-broker-test/internal/health"
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

	controller := httpController.InitController(deps.trades, deps.healthz)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ListenAddr)
	log.Printf("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, controller); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func initConfig() (cfg struct {
	DBPath     string
	ListenAddr string
}) {
	// Command line flags
	flag.StringVar(&cfg.DBPath, "db", "data.db", "path to SQLite database")
	flag.StringVar(&cfg.ListenAddr, "listen", "8080", "HTTP server listen address")
	flag.Parse()

	return cfg
}

func initDependencies(dbPath string) (deps struct {
	repositoryFactory *sqlite.RepositoryFactory
	trades            *service.Trades
	healthz           health.Healthchecking
}, closeDependencies func(), err error) {
	deps.repositoryFactory, err = sqlite.InitRepositoryFactory(sqlite.Config{DSN: dbPath})
	if err != nil {
		return deps, nil, fmt.Errorf("sqlite.InitRepositoryFactory: %w", err)
	}

	deps.trades = &service.Trades{
		TradesRepo: deps.repositoryFactory.NewTradesRepository(),
	}

	deps.healthz = health.Health{Components: map[string]health.Healthchecking{
		"repositoryFactory": deps.repositoryFactory,
	}}

	return deps, func() {
		if err := deps.repositoryFactory.Close(); err != nil {
			slog.Error("Failed to close repository factory: " + err.Error())
		}
	}, nil
}

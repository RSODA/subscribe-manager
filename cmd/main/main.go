package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RSODA/subscribe-manager/internal/config"
	"github.com/RSODA/subscribe-manager/internal/db"
	"github.com/RSODA/subscribe-manager/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	ctx := context.Background()

	err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		panic(err)
	}

	logConfig := config.NewLoggerLevel()

	logger, err := logger.NewLogger(string(logConfig))
	if err != nil {
		fmt.Printf("Error initializing logger: %v\n", err)
		panic(err)
	}

	defer logger.Sync()

	logger.Infow("Logger initialized", "level", logConfig)

	dbCfg, err := config.NewPGConfig()
	if err != nil {
		logger.Fatalw("Error initializing database", "error", err)
	}

	database, err := pgxpool.New(ctx, dbCfg.DSN())
	if err != nil {
		logger.Fatalw("Error connecting to database", "error", err)
	}

	dbClient := db.NewDB(database, logger)
	err = dbClient.Ping(ctx)
	if err != nil {
		logger.Fatalw("Error pinging database", "error", err)
	}

	<-shutdownCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	logger.Infow("Shutting down gracefully")

	dbClient.Close(shutdownCtx)
}

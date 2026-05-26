package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RSODA/subscribe-manager/internal/config"
	"github.com/RSODA/subscribe-manager/internal/db"
	"github.com/RSODA/subscribe-manager/internal/handler"
	"github.com/RSODA/subscribe-manager/internal/logger"
	postgresrepository "github.com/RSODA/subscribe-manager/internal/repository/postgres"
	subscriptionservice "github.com/RSODA/subscribe-manager/internal/service/subscription"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const migrationsPath = "file://migrations"

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

	logger.Infow("Database connection established")

	dbCtx, dbCancel := context.WithTimeout(ctx, time.Second*20)
	defer dbCancel()

	dbClient := db.NewDB(database, logger)
	err = dbClient.Ping(dbCtx)
	if err != nil {
		logger.Fatalw("Error pinging database", "error", err)
	}

	if err := runMigrations(database, logger); err != nil {
		logger.Fatalw("Error applying migrations", "error", err)
	}

	subscriptionRepository := postgresrepository.NewPostgresRepository(dbClient, logger)
	subscriptionService := subscriptionservice.NewSubscriptionService(subscriptionRepository, logger)
	httpRouter := handler.NewRouter(subscriptionService, logger)

	server := &http.Server{
		Addr:    ":8080",
		Handler: httpRouter.InitRoutes(),
	}

	go func() {
		logger.Infow("HTTP server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalw("Error starting HTTP server", "error", err)
		}
	}()

	<-shutdownCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	logger.Infow("Shutting down gracefully")

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorw("Error shutting down HTTP server", "error", err)
	}

	dbClient.Close(shutdownCtx)
}

func runMigrations(database *pgxpool.Pool, logger *zap.SugaredLogger) error {
	sqlDB := stdlib.OpenDBFromPool(database)
	defer sqlDB.Close()

	driver, err := migratepgx.WithInstance(sqlDB, &migratepgx.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance(migrationsPath, "pgx5", driver)
	if err != nil {
		return err
	}
	defer func() {
		sourceErr, databaseErr := migration.Close()
		if sourceErr != nil {
			logger.Warnw("Error closing migration source", "error", sourceErr)
		}
		if databaseErr != nil {
			logger.Warnw("Error closing migration database", "error", databaseErr)
		}
	}()

	if err = migration.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Infow("Database migrations are up to date")
			return nil
		}

		return err
	}

	logger.Infow("Database migrations applied")
	return nil
}

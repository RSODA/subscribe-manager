package app

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
	applogger "github.com/RSODA/subscribe-manager/internal/logger"
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

func Run() error {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	ctx := context.Background()

	if err := config.Load(); err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	logConfig := config.NewLoggerLevel()
	logger, err := applogger.NewLogger(string(logConfig))
	if err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}
	defer logger.Sync()

	httpCfg := config.NewHTTPConfig()

	logger.Infow("Logger initialized", "level", logConfig)

	dbCfg, err := config.NewPGConfig()
	if err != nil {
		return fmt.Errorf("error initializing database config: %w", err)
	}

	database, err := pgxpool.New(ctx, dbCfg.DSN())
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer database.Close()

	logger.Infow("Database connection established")

	dbCtx, dbCancel := context.WithTimeout(ctx, 60*time.Second)
	defer dbCancel()

	dbClient := db.NewDB(database, logger)
	err = dbClient.Ping(dbCtx)

	if err := runMigrations(database, logger); err != nil {
		return fmt.Errorf("error applying migrations: %w", err)
	}

	subscriptionRepository := postgresrepository.NewPostgresRepository(dbClient, logger)
	subscriptionService := subscriptionservice.NewSubscriptionService(subscriptionRepository, logger)
	httpRouter := handler.NewRouter(subscriptionService, logger)

	server := &http.Server{
		Addr:    httpCfg.Address(),
		Handler: httpRouter.InitRoutes(),
	}

	go func() {
		logger.Infow("HTTP server started", "addr", server.Addr)
		if serveErr := server.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
			logger.Fatalw("Error starting HTTP server", "error", serveErr)
		}
	}()

	<-shutdownCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	logger.Infow("Shutting down gracefully")

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorw("Error shutting down HTTP server", "error", err)
	}

	dbClient.Close(shutdownCtx)
	return nil
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

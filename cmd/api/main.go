package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ivmello/go-api-template/internal/app"
	"github.com/ivmello/go-api-template/internal/config"
	"github.com/ivmello/go-api-template/internal/infrastructure/database/postgres"
	"github.com/ivmello/go-api-template/internal/infrastructure/cache"
	"github.com/ivmello/go-api-template/internal/infrastructure/telemetry"
	"golang.org/x/sync/errgroup"
)

// @title         Go API Template
// @version       1.0
// @description   A RESTful API template built with Go.
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   https://github.com/ivmello/go-api-template/issues
// @contact.email your-email@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := telemetry.NewLogger(cfg)
	defer logger.Sync()

	// Setup OpenTelemetry
	tp, err := telemetry.SetupTracing(ctx, cfg)
	if err != nil {
		logger.Error("Failed to set up tracing", "error", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error("Error shutting down tracer provider", "error", err)
		}
	}()

	// Setup database connection
	db, err := postgres.NewClient(ctx, cfg)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	if err := postgres.RunMigrations(cfg); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Setup Redis connection
	redisClient, err := cache.NewRedisClient(ctx, cfg)
	if err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()

	// Create application
	application := app.New(ctx, cfg, db, redisClient, logger)

	// Start the application
	g, gCtx := errgroup.WithContext(ctx)

	// Start HTTP server
	g.Go(func() error {
		return application.StartHTTPServer(gCtx)
	})

	// Start gRPC server
	g.Go(func() error {
		return application.StartGRPCServer(gCtx)
	})

	// Handle shutdown signals
	g.Go(func() error {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-signalChan:
			logger.Info("Received signal", "signal", sig)
			cancel()
		case <-gCtx.Done():
			return gCtx.Err()
		}

		return nil
	})

	// Wait for all goroutines to finish
	if err := g.Wait(); err != nil {
		logger.Error("Error during shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("Application shutdown completed")
}
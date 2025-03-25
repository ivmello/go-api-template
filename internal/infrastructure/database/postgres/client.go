package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ivmello/go-api-template/internal/config"
)

// NewClient creates a new PostgreSQL client
func NewClient(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	// Create connection pool configuration
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.GetDSN())
	if err != nil {
		return nil, err
	}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

// RunMigrations runs all database migrations
func RunMigrations(cfg *config.Config) error {
	migrator, err := NewMigrator(cfg.Database)
	if err != nil {
		return err
	}
	defer migrator.Close()

	return migrator.Up()
}
package postgres

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ivmello/go-api-template/internal/config"
)

// Migrator handles database migrations
type Migrator struct {
	migrate *migrate.Migrate
}

// NewMigrator creates a new Migrator
func NewMigrator(dbConfig config.DatabaseConfig) (*Migrator, error) {
	// Create DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SSLMode)

	// Create migrate instance
	m, err := migrate.New(dbConfig.MigrationSource, dsn)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		migrate: m,
	}, nil
}

// Up runs all up migrations
func (m *Migrator) Up() error {
	err := m.migrate.Up()
	if err == migrate.ErrNoChange {
		return nil // No migrations to apply is not an error
	}
	return err
}

// Down runs one down migration
func (m *Migrator) Down() error {
	return m.migrate.Steps(-1)
}

// Version gets current migration version
func (m *Migrator) Version() (uint, bool, error) {
	return m.migrate.Version()
}

// Close closes the migrator
func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return dbErr
}
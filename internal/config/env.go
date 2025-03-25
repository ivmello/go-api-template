package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	JWT        JWTConfig
	Telemetry  TelemetryConfig
	ExternalAPI ExternalAPIConfig
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string
	Environment string
	Port        int
	GRPCPort    int
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	Name           string
	SSLMode        string
	MigrationSource string
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret        string
	ExpirationHours int
}

// TelemetryConfig holds telemetry configuration
type TelemetryConfig struct {
	ServiceName      string
	ExporterEndpoint string
}

// ExternalAPIConfig holds configuration for external API calls
type ExternalAPIConfig struct {
	Timeout time.Duration
}

// Load reads environment variables and returns a Config struct
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	return &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "go-api-template"),
			Environment: getEnv("ENVIRONMENT", "development"),
			Port:        getEnvAsInt("PORT", 8080),
			GRPCPort:    getEnvAsInt("GRPC_PORT", 9090),
		},
		Database: DatabaseConfig{
			Host:           getEnv("DB_HOST", "localhost"),
			Port:           getEnvAsInt("DB_PORT", 5432),
			User:           getEnv("DB_USER", "postgres"),
			Password:       getEnv("DB_PASSWORD", "postgres"),
			Name:           getEnv("DB_NAME", "api_db"),
			SSLMode:        getEnv("DB_SSL_MODE", "disable"),
			MigrationSource: getEnv("DB_MIGRATION_SOURCE", "file://internal/infrastructure/database/migrations/postgres"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "default-secret-key-for-development-only"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		},
		Telemetry: TelemetryConfig{
			ServiceName:      getEnv("OTEL_SERVICE_NAME", "go-api-template"),
			ExporterEndpoint: getEnv("OTEL_EXPORTER_ENDPOINT", "localhost:4317"),
		},
		ExternalAPI: ExternalAPIConfig{
			Timeout: getEnvAsDuration("EXTERNAL_API_TIMEOUT", 5*time.Second),
		},
	}, nil
}

// Helper functions to get environment variables
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// GetRedisAddr returns the Redis connection string
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
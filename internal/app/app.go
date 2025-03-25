package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log/slog"

	"github.com/ivmello/go-api-template/internal/config"
	"github.com/ivmello/go-api-template/internal/core/auth"
	"github.com/ivmello/go-api-template/internal/core/message"
	"github.com/ivmello/go-api-template/internal/infrastructure/http_client"
)

// Application holds all dependencies of the application
type Application struct {
	config      *config.Config
	db          *pgxpool.Pool
	redisClient *redis.Client
	logger      *slog.Logger
	httpClient  *http_client.Client

	// Services
	authService    *auth.Service
	messageService *message.Service
}

// New creates a new Application with all dependencies
func New(ctx context.Context, cfg *config.Config, db *pgxpool.Pool, redisClient *redis.Client, logger *slog.Logger) *Application {
	// Initialize HTTP client for external APIs
	httpClient := http_client.NewClient(cfg.ExternalAPI.Timeout)

	// Initialize repositories
	authRepo := auth.NewRepository(db)
	messageRepo := message.NewRepository(db)

	// Initialize services
	authService := auth.NewService(authRepo, cfg.JWT)
	messageService := message.NewService(messageRepo)

	return &Application{
		config:         cfg,
		db:             db,
		redisClient:    redisClient,
		logger:         logger,
		httpClient:     httpClient,
		authService:    authService,
		messageService: messageService,
	}
}

// Services returns all application services
func (a *Application) Services() struct {
	Auth    *auth.Service
	Message *message.Service
} {
	return struct {
		Auth    *auth.Service
		Message *message.Service
	}{
		Auth:    a.authService,
		Message: a.messageService,
	}
}
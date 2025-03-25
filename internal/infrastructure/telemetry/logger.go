package telemetry

import (
	"log/slog"
	"os"

	"github.com/ivmello/go-api-template/internal/config"
)

// NewLogger creates a new structured logger
func NewLogger(cfg *config.Config) *slog.Logger {
	// Configure JSON handler for production and text handler for development
	var handler slog.Handler
	if cfg.App.Environment == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	// Create logger with service name attribute
	logger := slog.New(handler).With(
		"service", cfg.App.Name,
		"environment", cfg.App.Environment,
	)

	return logger
}
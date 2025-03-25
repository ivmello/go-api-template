package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ivmello/go-api-template/internal/handlers/http/auth"
	"github.com/ivmello/go-api-template/internal/handlers/http/healthcheck"
	"github.com/ivmello/go-api-template/internal/handlers/http/message"
	"github.com/ivmello/go-api-template/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// StartHTTPServer starts the HTTP server
func (a *Application) StartHTTPServer(ctx context.Context) error {
	// Set Gin mode
	if a.config.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router with middleware
	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.LoggerMiddleware(a.logger),
		otelgin.Middleware(a.config.Telemetry.ServiceName),
		middleware.CORSMiddleware(),
	)

	// Register routes
	a.registerHTTPRoutes(router)

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.App.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		a.logger.Info("Starting HTTP server", "port", a.config.App.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("HTTP server failed", "error", err)
		}
	}()

	// Wait for context cancelation (shutdown signal)
	<-ctx.Done()
	a.logger.Info("Shutting down HTTP server")

	// Create a timeout context for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		a.logger.Error("HTTP server shutdown failed", "error", err)
		return err
	}

	a.logger.Info("HTTP server shutdown completed")
	return nil
}

// registerHTTPRoutes registers all HTTP routes
func (a *Application) registerHTTPRoutes(router *gin.Engine) {
	// Root route for health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"name":   a.config.App.Name,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Health check
		healthcheckHandler := healthcheck.NewHandler(a.db, a.redisClient)
		healthGroup := v1.Group("/health")
		{
			healthGroup.GET("", healthcheckHandler.Check)
			healthGroup.GET("/liveness", healthcheckHandler.Liveness)
			healthGroup.GET("/readiness", healthcheckHandler.Readiness)
		}

		// Auth routes
		authHandler := auth.NewHandler(a.Services().Auth)
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", middleware.AuthMiddleware(), authHandler.Me)
		}

		// Message routes
		messageHandler := message.NewHandler(a.Services().Message)
		messageGroup := v1.Group("/messages")
		{
			messageGroup.GET("", messageHandler.GetAll)                                // Public
			messageGroup.GET("/:id", middleware.AuthMiddleware(), messageHandler.Get)  // Protected
			messageGroup.POST("", middleware.AuthMiddleware(), messageHandler.Create)  // Protected
			messageGroup.PUT("/:id", middleware.AuthMiddleware(), messageHandler.Update) // Protected
			messageGroup.DELETE("/:id", middleware.AuthMiddleware(), messageHandler.Delete) // Protected
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
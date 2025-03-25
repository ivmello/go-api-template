package healthcheck

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	httpTransport "github.com/ivmello/go-api-template/internal/transport/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Handler handles health check HTTP requests
type Handler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

// NewHandler creates a new health check handler
func NewHandler(db *pgxpool.Pool, redis *redis.Client) *Handler {
	return &Handler{
		db:    db,
		redis: redis,
	}
}

// Check performs a full health check
// @Summary Health check
// @Description Check health of all services
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} httpTransport.HealthCheckResponse
// @Failure 500 {object} httpTransport.ErrorResponse
// @Router /api/v1/health [get]
func (h *Handler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Check database
	dbErr := h.db.Ping(ctx)
	dbStatus := "up"
	if dbErr != nil {
		dbStatus = "down"
	}

	// Check Redis
	redisErr := h.redis.Ping(ctx).Err()
	redisStatus := "up"
	if redisErr != nil {
		redisStatus = "down"
	}

	// Determine overall status
	overallStatus := "ok"
	if dbErr != nil || redisErr != nil {
		overallStatus = "degraded"
	}

	// Return status
	c.JSON(http.StatusOK, httpTransport.HealthCheckResponse{
		Status: overallStatus,
		Services: map[string]string{
			"api":      "up",
			"database": dbStatus,
			"cache":    redisStatus,
		},
		Timestamp: time.Now(),
	})
}

// Liveness checks if the application is running
// @Summary Liveness probe
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} httpTransport.SuccessResponse
// @Router /api/v1/health/liveness [get]
func (h *Handler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, httpTransport.SuccessResponse{
		Message: "Service is alive",
	})
}

// Readiness checks if the application is ready to serve traffic
// @Summary Readiness probe
// @Description Check if the API is ready to serve traffic
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} httpTransport.SuccessResponse
// @Failure 503 {object} httpTransport.ErrorResponse
// @Router /api/v1/health/readiness [get]
func (h *Handler) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Check database connection
	if err := h.db.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, httpTransport.ErrorResponse{
			Error: "Database not ready",
		})
		return
	}

	// Check Redis connection
	if err := h.redis.Ping(ctx).Err(); err != nil {
		c.JSON(http.StatusServiceUnavailable, httpTransport.ErrorResponse{
			Error: "Cache not ready",
		})
		return
	}

	c.JSON(http.StatusOK, httpTransport.SuccessResponse{
		Message: "Service is ready",
	})
}

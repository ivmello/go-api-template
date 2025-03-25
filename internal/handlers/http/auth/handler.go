package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivmello/go-api-template/internal/core/auth"
	"github.com/ivmello/go-api-template/internal/transport/http"
)

// Handler handles authentication HTTP requests
type Handler struct {
	service *auth.Service
}

// NewHandler creates a new auth handler
func NewHandler(service *auth.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email, password and name
// @Tags auth
// @Accept json
// @Produce json
// @Param request body http.RegisterRequest true "User registration data"
// @Success 201 {object} http.UserResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 409 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req http.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, http.ErrorResponse{Error: "Invalid request format"})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, http.ErrorResponse{Error: err.Error()})
		return
	}

	// Create user
	user, err := h.service.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, auth.ErrEmailAlreadyExists) {
			status = http.StatusConflict
		}
		c.JSON(status, http.ErrorResponse{Error: err.Error()})
		return
	}

	// Return user
	c.JSON(http.StatusCreated, http.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	})
}

// Login handles user login
// @Summary Login user
// @Description Login with email and password to get a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body http.LoginRequest true "User login data"
// @Success 200 {object} http.TokenResponse
// @Failure 400 {object} http.ErrorResponse
// @Failure 401 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req http.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, http.ErrorResponse{Error: "Invalid request format"})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, http.ErrorResponse{Error: err.Error()})
		return
	}

	// Login user
	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, auth.ErrInvalidCredentials) {
			status = http.StatusUnauthorized
		}
		c.JSON(status, http.ErrorResponse{Error: err.Error()})
		return
	}

	// Return token
	c.JSON(http.StatusOK, http.TokenResponse{
		Token: token,
	})
}

// Me returns the current user
// @Summary Get current user
// @Description Get details of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} http.UserResponse
// @Failure 401 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /api/v1/auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, http.ErrorResponse{Error: "Not authenticated"})
		return
	}

	// Get user details
	user, err := h.service.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, auth.ErrUserNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, http.ErrorResponse{Error: err.Error()})
		return
	}

	// Return user
	c.JSON(http.StatusOK, http.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	})
}
package message

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ivmello/go-api-template/internal/core/message"
	httpTransport "github.com/ivmello/go-api-template/internal/transport/http"
)

// Handler handles message HTTP requests
type Handler struct {
	service *message.Service
}

// NewHandler creates a new message handler
func NewHandler(service *message.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetAll retrieves all messages
// @Summary Get all messages
// @Description Get a list of all messages
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {array} httpTransport.MessageResponse
// @Failure 500 {object} httpTransport.ErrorResponse
// @Router /api/v1/messages [get]
func (h *Handler) GetAll(c *gin.Context) {
	messages, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpTransport.ErrorResponse{Error: err.Error()})
		return
	}

	// Map domain objects to response objects
	response := make([]httpTransport.MessageResponse, len(messages))
	for i, msg := range messages {
		response[i] = httpTransport.MessageResponse{
			ID:        msg.ID,
			UserID:    msg.UserID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, response)
}

// Get retrieves a single message by ID
// @Summary Get message by ID
// @Description Get a message by its ID
// @Tags messages
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Message ID"
// @Success 200 {object} httpTransport.MessageResponse
// @Failure 400 {object} httpTransport.ErrorResponse
// @Failure 401 {object} httpTransport.ErrorResponse
// @Failure 404 {object} httpTransport.ErrorResponse
// @Failure 500 {object} httpTransport.ErrorResponse
// @Router /api/v1/messages/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")

	msg, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, message.ErrMessageNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, httpTransport.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, httpTransport.MessageResponse{
		ID:        msg.ID,
		UserID:    msg.UserID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	})
}

// Create creates a new message
// @Summary Create message
// @Description Create a new message
// @Tags messages
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body httpTransport.CreateMessageRequest true "Message content"
// @Success 201 {object} httpTransport.MessageResponse
// @Failure 400 {object} httpTransport.ErrorResponse
// @Failure 401 {object} httpTransport.ErrorResponse
// @Failure 500 {object} httpTransport.ErrorResponse
// @Router /api/v1/messages [post]
func (h *Handler) Create(c *gin.Context) {
	var req httpTransport.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httpTransport.ErrorResponse{Error: "Invalid request format"})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, httpTransport.ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, httpTransport.ErrorResponse{Error: "Not authenticated"})
		return
	}

	// Create message
	msg, err := h.service.Create(c.Request.Context(), userID.(string), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, httpTransport.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, httpTransport.MessageResponse{
		ID:        msg.ID,
		UserID:    msg.UserID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	})
}

// Update updates a message
// @Summary Update message
// @Description Update an existing message
// @Tags messages
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Message ID"
// @Param request body httpTransport.UpdateMessageRequest true "Updated message content"
// @Success 200 {object} httpTransport.SuccessResponse
// @Failure 400 {object} httpTransport.ErrorResponse
// @Failure 401 {object} httpTransport.ErrorResponse
// @Failure 403 {object} httpTransport.ErrorResponse
// @Failure 404 {object} httpTransport.ErrorResponse
// @Failure 500 {object} httpTransport.ErrorResponse
// @Router /api/v1/messages/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req httpTransport.UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, httpTransport.ErrorResponse{Error: "Invalid request format"})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, httpTransport.ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, httpTransport.ErrorResponse{Error: "Not authenticated"})
		return
	}

	// Update message
	err := h.service.Update(c.Request.Context(), id, userID.(string), req.Content)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, message.ErrMessageNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, message.ErrForbidden) {
			status = http.StatusForbidden
		}
		c.JSON(status, httpTransport.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, httpTransport.SuccessResponse{
		Message: "Message updated successfully",
	})
}

// Delete deletes a message
// @Summary Delete message
// @Description Delete an existing message
// @Tags messages
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path string true "Message ID"
// @Success 200 {object} httpTransport.SuccessResponse
// @Failure 401 {object} httpTransport.ErrorResponse
// @Failure 403 {object} httpTransport.ErrorResponse
// @Failure 404 {object} httpTransport.ErrorResponse
// @Failure 500 {object} httpTransport.ErrorResponse
// @Router /api/v1/messages/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, httpTransport.ErrorResponse{Error: "Not authenticated"})
		return
	}

	// Delete message
	err := h.service.Delete(c.Request.Context(), id, userID.(string))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, message.ErrMessageNotFound) {
			status = http.StatusNotFound
		} else if errors.Is(err, message.ErrForbidden) {
			status = http.StatusForbidden
		}
		c.JSON(status, httpTransport.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, httpTransport.SuccessResponse{
		Message: "Message deleted successfully",
	})
}
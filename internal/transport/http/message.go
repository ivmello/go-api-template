package http

import (
	"errors"
	"time"
)

// CreateMessageRequest represents a request to create a message
type CreateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// Validate validates the create message request
func (r *CreateMessageRequest) Validate() error {
	if r.Content == "" {
		return errors.New("content is required")
	}
	if len(r.Content) > 1000 {
		return errors.New("content must be less than 1000 characters")
	}
	return nil
}

// UpdateMessageRequest represents a request to update a message
type UpdateMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// Validate validates the update message request
func (r *UpdateMessageRequest) Validate() error {
	if r.Content == "" {
		return errors.New("content is required")
	}
	if len(r.Content) > 1000 {
		return errors.New("content must be less than 1000 characters")
	}
	return nil
}

// MessageResponse represents a message response
type MessageResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
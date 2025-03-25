package http

import (
	"errors"
	"time"
)

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// Validate validates the register request
func (r *RegisterRequest) Validate() error {
	if len(r.Email) < 5 {
		return errors.New("email must be at least 5 characters")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if len(r.Name) < 2 {
		return errors.New("name must be at least 2 characters")
	}
	return nil
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Validate validates the login request
func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// UserResponse represents a user response
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	Token string `json:"token"`
}
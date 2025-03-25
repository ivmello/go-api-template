package auth

import (
	"context"
	"errors"

	"github.com/ivmello/go-api-template/internal/config"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// Service provides authentication operations
type Service struct {
	repo *Repository
	jwt  config.JWTConfig
}

// NewService creates a new authentication service
func NewService(repo *Repository, jwtConfig config.JWTConfig) *Service {
	return &Service{
		repo: repo,
		jwt:  jwtConfig,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, email, password, name string) (*User, error) {
	// Create user with hashed password
	user, err := NewUser(email, password, name)
	if err != nil {
		return nil, err
	}

	// Save user to database
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	// Get user by email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	// Check password
	if !user.ComparePassword(password) {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := user.GenerateToken(s.jwt.ExpirationHours)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID gets a user by ID
func (s *Service) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}
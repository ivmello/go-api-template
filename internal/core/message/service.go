package message

import (
	"context"
)

// Service provides message operations
type Service struct {
	repo *Repository
}

// NewService creates a new message service
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a new message
func (s *Service) Create(ctx context.Context, userID, content string) (*Message, error) {
	// Create message
	message := NewMessage(userID, content)

	// Save message to database
	if err := s.repo.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

// GetAll retrieves all messages
func (s *Service) GetAll(ctx context.Context) ([]*Message, error) {
	return s.repo.GetAll(ctx)
}

// GetByID retrieves a message by ID
func (s *Service) GetByID(ctx context.Context, id string) (*Message, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates a message
func (s *Service) Update(ctx context.Context, id, userID, content string) error {
	return s.repo.Update(ctx, id, userID, content)
}

// Delete deletes a message
func (s *Service) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
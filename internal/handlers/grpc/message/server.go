package message

import (
	"context"
	"errors"
	"time"

	"github.com/ivmello/go-api-template/internal/core/message"
	"github.com/ivmello/go-api-template/internal/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the MessageService gRPC server
type Server struct {
	UnimplementedMessageServiceServer
	service *message.Service
}

// NewServer creates a new message gRPC server
func NewServer(service *message.Service) *Server {
	return &Server{
		service: service,
	}
}

// CreateMessage creates a new message
func (s *Server) CreateMessage(ctx context.Context, req *CreateMessageRequest) (*MessageResponse, error) {
	// Validate request
	if req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}

	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	// Create message
	msg, err := s.service.Create(ctx, userID, req.Content)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Return message
	return &MessageResponse{
		Id:        msg.ID,
		UserId:    msg.UserID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt.Format(time.RFC3339),
		UpdatedAt: msg.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// GetMessage returns a message by ID
func (s *Server) GetMessage(ctx context.Context, req *GetMessageRequest) (*MessageResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get message
	msg, err := s.service.GetByID(ctx, req.Id)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, message.ErrMessageNotFound) {
			code = codes.NotFound
		}
		return nil, status.Error(code, err.Error())
	}

	// Return message
	return &MessageResponse{
		Id:        msg.ID,
		UserId:    msg.UserID,
		Content:   msg.Content,
		CreatedAt: msg.CreatedAt.Format(time.RFC3339),
		UpdatedAt: msg.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// UpdateMessage updates a message
func (s *Server) UpdateMessage(ctx context.Context, req *UpdateMessageRequest) (*EmptyResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}

	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	// Update message
	err = s.service.Update(ctx, req.Id, userID, req.Content)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, message.ErrMessageNotFound) {
			code = codes.NotFound
		} else if errors.Is(err, message.ErrForbidden) {
			code = codes.PermissionDenied
		}
		return nil, status.Error(code, err.Error())
	}

	return &EmptyResponse{}, nil
}

// DeleteMessage deletes a message
func (s *Server) DeleteMessage(ctx context.Context, req *DeleteMessageRequest) (*EmptyResponse, error) {
	// Validate request
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	// Delete message
	err = s.service.Delete(ctx, req.Id, userID)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, message.ErrMessageNotFound) {
			code = codes.NotFound
		} else if errors.Is(err, message.ErrForbidden) {
			code = codes.PermissionDenied
		}
		return nil, status.Error(code, err.Error())
	}

	return &EmptyResponse{}, nil
}

// ListMessages lists all messages
func (s *Server) ListMessages(ctx context.Context, req *ListMessagesRequest) (*ListMessagesResponse, error) {
	// Get messages
	msgs, err := s.service.GetAll(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert messages to response format
	responses := make([]*MessageResponse, len(msgs))
	for i, msg := range msgs {
		responses[i] = &MessageResponse{
			Id:        msg.ID,
			UserId:    msg.UserID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Format(time.RFC3339),
			UpdatedAt: msg.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &ListMessagesResponse{
		Messages: responses,
	}, nil
}
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/ivmello/go-api-template/internal/core/auth"
	"github.com/ivmello/go-api-template/internal/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the AuthService gRPC server
type Server struct {
	UnimplementedAuthServiceServer
	service *auth.Service
}

// NewServer creates a new auth gRPC server
func NewServer(service *auth.Service) *Server {
	return &Server{
		service: service,
	}
}

// Register registers a new user
func (s *Server) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "email, password, and name are required")
	}

	// Register user
	user, err := s.service.Register(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, auth.ErrEmailAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Error(code, err.Error())
	}

	// Return user
	return &UserResponse{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

// Login authenticates a user
func (s *Server) Login(ctx context.Context, req *LoginRequest) (*TokenResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Login user
	token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, auth.ErrInvalidCredentials) {
			code = codes.Unauthenticated
		}
		return nil, status.Error(code, err.Error())
	}

	// Return token
	return &TokenResponse{
		Token: token,
	}, nil
}

// GetCurrentUser returns the current user
func (s *Server) GetCurrentUser(ctx context.Context, req *GetCurrentUserRequest) (*UserResponse, error) {
	// Get user ID from context (set by auth middleware)
	userID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not authenticated")
	}

	// Get user details
	user, err := s.service.GetUserByID(ctx, userID)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, auth.ErrUserNotFound) {
			code = codes.NotFound
		}
		return nil, status.Error(code, err.Error())
	}

	// Return user
	return &UserResponse{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
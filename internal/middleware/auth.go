package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ivmello/go-api-template/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Key for user ID in context
type contextKey string
const UserIDKey contextKey = "user_id"

// AuthMiddleware validates JWT tokens for HTTP requests
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Check if header starts with Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			return
		}

		// Extract token
		tokenString := parts[1]

		// Parse and validate token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set user ID in context
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// GRPCAuth returns a unary server interceptor for authenticating gRPC requests
func GRPCAuth() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip authentication for certain methods
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// Get metadata from context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is required")
		}

		// Get authorization token
		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is required")
		}

		// Check token format
		authHeader := values[0]
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format, expected 'Bearer {token}'")
		}

		// Validate token
		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Add user ID to context
		newCtx := context.WithValue(ctx, UserIDKey, claims.UserID)

		// Call the handler with the new context
		return handler(newCtx, req)
	}
}

// GRPCStreamAuth returns a stream server interceptor for authenticating gRPC stream requests
func GRPCStreamAuth() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Skip authentication for certain methods
		if isPublicMethod(info.FullMethod) {
			return handler(srv, ss)
		}

		// Get metadata from context
		ctx := ss.Context()
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return status.Errorf(codes.Unauthenticated, "metadata is required")
		}

		// Get authorization token
		values := md.Get("authorization")
		if len(values) == 0 {
			return status.Errorf(codes.Unauthenticated, "authorization token is required")
		}

		// Check token format
		authHeader := values[0]
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return status.Errorf(codes.Unauthenticated, "invalid authorization format, expected 'Bearer {token}'")
		}

		// Validate token
		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Create a new context with the user ID
		newCtx := context.WithValue(ctx, UserIDKey, claims.UserID)

		// Wrap the server stream with the new context
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          newCtx,
		}

		// Call the handler with the wrapped stream
		return handler(srv, wrappedStream)
	}
}

// GetUserIDFromContext extracts the user ID from the context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

// Helper function to check if a method is public
func isPublicMethod(method string) bool {
	// List of methods that don't require authentication
	publicMethods := []string{
		"/auth.AuthService/Login",
		"/auth.AuthService/Register",
	}

	for _, m := range publicMethods {
		if m == method {
			return true
		}
	}
	return false
}

// wrappedServerStream wraps a grpc.ServerStream with a new context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the wrapped context
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
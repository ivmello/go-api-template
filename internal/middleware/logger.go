package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggerMiddleware creates a middleware function for logging HTTP requests
func LoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// Log request details
		logger.Info("HTTP request",
			"method", method,
			"path", path,
			"status", statusCode,
			"latency", latency,
			"client_ip", clientIP,
			"user_agent", c.Request.UserAgent(),
		)
	}
}

// GRPCLogger returns a unary server interceptor for logging gRPC requests
func GRPCLogger(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Start timer
		start := time.Now()
		method := info.FullMethod

		// Process request
		resp, err := handler(ctx, req)

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		code := codes.OK
		if err != nil {
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			} else {
				code = codes.Internal
			}
		}

		// Log request details
		logger.Info("gRPC request",
			"method", method,
			"status", code.String(),
			"duration", duration,
		)

		return resp, err
	}
}

// GRPCStreamLogger returns a stream server interceptor for logging gRPC stream requests
func GRPCStreamLogger(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Start timer
		start := time.Now()
		method := info.FullMethod

		// Process stream
		err := handler(srv, ss)

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		code := codes.OK
		if err != nil {
			if s, ok := status.FromError(err); ok {
				code = s.Code()
			} else {
				code = codes.Internal
			}
		}

		// Log stream details
		logger.Info("gRPC stream",
			"method", method,
			"status", code.String(),
			"duration", duration,
			"is_client_stream", info.IsClientStream,
			"is_server_stream", info.IsServerStream,
		)

		return err
	}
}
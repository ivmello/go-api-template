package app

import (
	"context"
	"fmt"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2"
	"github.com/ivmello/go-api-template/internal/handlers/grpc/auth"
	"github.com/ivmello/go-api-template/internal/handlers/grpc/message"
	"github.com/ivmello/go-api-template/internal/middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// StartGRPCServer starts the gRPC server
func (a *Application) StartGRPCServer(ctx context.Context) error {
	// Create listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.App.GRPCPort))
	if err != nil {
		a.logger.Error("Failed to listen for gRPC", "error", err)
		return err
	}

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				middleware.GRPCLogger(a.logger),
				otelgrpc.UnaryServerInterceptor(),
				middleware.GRPCAuth(),
			),
		),
		grpc.StreamInterceptor(
			grpcMiddleware.ChainStreamServer(
				middleware.GRPCStreamLogger(a.logger),
				otelgrpc.StreamServerInterceptor(),
				middleware.GRPCStreamAuth(),
			),
		),
	)

	// Register services
	a.registerGRPCServices(grpcServer)

	// Start server in a goroutine
	go func() {
		a.logger.Info("Starting gRPC server", "port", a.config.App.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			a.logger.Error("gRPC server failed", "error", err)
		}
	}()

	// Wait for context cancelation (shutdown signal)
	<-ctx.Done()
	a.logger.Info("Shutting down gRPC server")

	// Gracefully stop the server
	grpcServer.GracefulStop()

	a.logger.Info("gRPC server shutdown completed")
	return nil
}

// registerGRPCServices registers all gRPC services
func (a *Application) registerGRPCServices(server *grpc.Server) {
	// Register Auth service
	authServer := auth.NewServer(a.Services().Auth)
	auth.RegisterAuthServiceServer(server, authServer)

	// Register Message service
	messageServer := message.NewServer(a.Services().Message)
	message.RegisterMessageServiceServer(server, messageServer)
}
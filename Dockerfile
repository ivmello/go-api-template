# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go.mod and go.sum first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app ./cmd/api

# Development stage with hot reload
FROM golang:1.24-alpine AS development

WORKDIR /app

# Install air for hot reloading and git for dependencies
RUN apk add --no-cache git \
    && go install github.com/air-verse/air@latest

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy entire codebase for development
COPY . .

# Expose HTTP and gRPC ports
EXPOSE 8080 9090

# Start with air for hot reloading
CMD ["air", "-c", ".air.toml"]

# Final production stage
FROM alpine:3.19 AS production

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /go/bin/app /app/

# Copy migration files
COPY --from=builder /app/internal/infrastructure/database/migrations /app/internal/infrastructure/database/migrations

# Expose HTTP and gRPC ports
EXPOSE 8080 9090

# Run the application
CMD ["/app/app"]
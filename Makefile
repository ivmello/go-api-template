# Variables
APP_NAME=go-api-template
MAIN_PATH=./cmd/api
BUILD_DIR=./bin

.PHONY: all build clean test coverage lint run run-dev docker-build docker-run docker-compose-up docker-compose-down help migrations-up migrations-down proto swagger

all: clean lint test build

build: # Build the application
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

clean: # Clean build directory
	rm -rf $(BUILD_DIR)
	go clean

test: # Run tests
	go test -v ./...

test-coverage: # Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint: # Lint the code
	golangci-lint run ./...

run: build # Build and run the application
	$(BUILD_DIR)/$(APP_NAME)

run-dev: # Run with hot reload using air
	air -c .air.toml

docker-build: # Build Docker image
	docker build -t $(APP_NAME) .

docker-run: docker-build # Run Docker container
	docker run -p 8080:8080 -p 9090:9090 --name $(APP_NAME) $(APP_NAME)

docker-compose-up: # Start all services with docker-compose
	docker-compose up -d

docker-compose-down: # Stop all services
	docker-compose down

migrations-create: # Create new migration files
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/infrastructure/database/migrations/postgres -seq $$name

migrations-up: # Run migrations up
	migrate -path internal/infrastructure/database/migrations/postgres -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrations-down: # Run migrations down
	migrate -path internal/infrastructure/database/migrations/postgres -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down 1

proto: # Generate gRPC code from proto files
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/*.proto

swagger: # Generate Swagger documentation
	swag init -g cmd/api/main.go -o ./docs

help: # Show help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?# .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
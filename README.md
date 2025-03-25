# Go API Template

A robust and scalable template for building RESTful APIs and gRPC services with Go, following best practices for clean architecture and modern Go development.

## Features

- **Clean Architecture**: Easy to understand, maintain, and extend
- **Multiple Protocol Support**: Both HTTP REST API and gRPC
- **Authentication**: JWT-based authentication system
- **Database Integration**: PostgreSQL with migrations
- **Caching**: Redis integration
- **Hot Reloading**: For efficient development workflow
- **OpenTelemetry**: Integrated monitoring with Grafana and Prometheus
- **Swagger Documentation**: Auto-generated API documentation
- **Graceful Shutdown**: Proper handling of shutdown signals
- **Structured Logging**: Using slog with JSON format for production
- **Docker Ready**: Multi-stage Docker builds and docker-compose for local development
- **HTTP Client**: Example implementation for external API calls with timeout and parallel requests

## Technology Stack

- Go 1.24
- PostgreSQL
- Redis
- gRPC
- Gin Web Framework
- OpenTelemetry
- Swagger
- Docker / Docker Compose
- Grafana / Prometheus / Jaeger

## Project Structure

The project follows a clean architecture approach with the following structure:

```
├── cmd/                         # Application entry points
│   └── api/                     # Main API application
├── internal/                    # Private application code
│   ├── app/                     # Application initialization
│   ├── config/                  # Configuration management
│   ├── core/                    # Business logic (domain)
│   │   ├── auth/                # Authentication domain
│   │   └── message/             # Message domain (example)
│   ├── handlers/                # Request handlers
│   │   ├── http/                # HTTP handlers
│   │   └── grpc/                # gRPC handlers
│   ├── middleware/              # Middleware components
│   ├── transport/               # DTOs and validation
│   └── infrastructure/          # External systems integration
│       ├── database/            # Database connections
│       ├── cache/               # Cache client
│       ├── http_client/         # HTTP client for external APIs
│       └── telemetry/           # Logging and tracing
├── pkg/                         # Reusable packages
├── proto/                       # Protocol Buffers definitions
└── test/                        # Test utilities and e2e tests
```

## Getting Started

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Make (optional but recommended)

### Installation

1. Clone the repository or use it as a template:

```bash
# Using as template
git clone https://github.com/ivmello/go-api-template.git your-project-name
cd your-project-name

# Rename the project (optional)
./scripts/rename-project.sh github.com/your-username/your-project-name
```

2. Start the development environment:

```bash
# Run the setup script
./scripts/setup.sh

# Start all services
make docker-compose-up
```

3. Run the API with hot reloading:

```bash
make run-dev
```

The API will be available at http://localhost:8080 and the gRPC server at localhost:9090.

## API Documentation

Swagger documentation is available at http://localhost:8080/swagger/index.html when running the API.

## Environment Variables

Configuration is done through environment variables. See `.env.example` for a list of all variables.

## Available Commands

```bash
# Show available commands
make help

# Run the application
make run

# Run with hot reloading
make run-dev

# Build the application
make build

# Run tests
make test

# Generate test coverage
make test-coverage

# Lint the code
make lint

# Create new migration
make migrations-create

# Run migrations up
make migrations-up

# Generate gRPC code
make proto

# Generate Swagger documentation
make swagger

# Docker operations
make docker-build
make docker-run
make docker-compose-up
make docker-compose-down
```

## Testing

The project includes both unit tests and end-to-end tests:

```bash
# Run all tests
make test

# Run with coverage report
make test-coverage
```

## Monitoring

- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Jaeger UI**: http://localhost:16686

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Credits

Created by [Ivan Mello](https://github.com/ivmello)
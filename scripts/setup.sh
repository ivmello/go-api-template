#!/bin/bash

set -e

echo "Setting up the Go API Template..."

# Create .env file from example if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
fi

# Make scripts executable
chmod +x scripts/*.sh

# Start services with Docker Compose
echo "Starting services with Docker Compose..."
docker-compose up -d postgres redis

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
sleep 5

# Run database migrations
echo "Running database migrations..."
docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE api_db;" || echo "Database already exists"
make migrations-up

echo "Setup completed successfully!"
echo "Run 'docker-compose up -d' to start all services"
echo "Run 'make run-dev' to start the API with hot reloading"
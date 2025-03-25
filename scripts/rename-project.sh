#!/bin/bash

set -e

# Current name in go.mod
CURRENT_NAME="github.com/ivmello/go-api-template"

# Check if a new name is provided
if [ -z "$1" ]; then
    echo "Error: Please provide a new module name as argument"
    echo "Usage: $0 github.com/username/project-name"
    exit 1
fi

NEW_NAME=$1

echo "Renaming project from $CURRENT_NAME to $NEW_NAME"

# Change module name in go.mod
sed -i "s|module $CURRENT_NAME|module $NEW_NAME|g" go.mod

# Update all imports in the codebase
find . -type f -name "*.go" -exec sed -i "s|$CURRENT_NAME|$NEW_NAME|g" {} \;

# Update docker-compose service names
sed -i "s|container_name: go-api-template|container_name: ${NEW_NAME##*/}|g" docker-compose.yml
sed -i "s|APP_NAME=go-api-template|APP_NAME=${NEW_NAME##*/}|g" docker-compose.yml
sed -i "s|OTEL_SERVICE_NAME=go-api-template|OTEL_SERVICE_NAME=${NEW_NAME##*/}|g" docker-compose.yml

# Update .env.example
sed -i "s|APP_NAME=go-api-template|APP_NAME=${NEW_NAME##*/}|g" .env.example
sed -i "s|OTEL_SERVICE_NAME=go-api-template|OTEL_SERVICE_NAME=${NEW_NAME##*/}|g" .env.example

# Update Makefile
sed -i "s|APP_NAME=go-api-template|APP_NAME=${NEW_NAME##*/}|g" Makefile

echo "Project successfully renamed to $NEW_NAME"
echo "You may need to run 'go mod tidy' to update dependencies"
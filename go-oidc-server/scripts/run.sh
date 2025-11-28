#!/bin/bash

# Navigate to the project directory
cd "$(dirname "$0")/.."

# Build the Go application
go build -o oidc-server ./cmd/server/main.go

# Run the OIDC server
./oidc-server
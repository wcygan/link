# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Service to build/run (hardcoded)
SERVICE_NAME=server
SERVICE_DIR=./$(SERVICE_NAME)
CMD_DIR=cmd/main.go
OUTPUT_BINARY=bin/$(SERVICE_NAME)

# Linting and Formatting
GOLINT=golangci-lint run
GOFMT=gofmt -w

# Docker parameters
DOCKER_COMPOSE=docker-compose

.PHONY: all build clean run lint fmt docker-up docker-down help

all: help

# Build the application
build:
	@echo "Building $(SERVICE_NAME)..."
	cd $(SERVICE_DIR) && $(GOBUILD) -o ../$(OUTPUT_BINARY) $(CMD_DIR)

# Clean the binary
clean:
	@echo "Cleaning $(SERVICE_NAME)..."
	rm -f $(OUTPUT_BINARY)
	cd $(SERVICE_DIR) && $(GOCLEAN)

# Run the application (requires building first)
run: build
	@echo "Running $(SERVICE_NAME)..."
	./$(OUTPUT_BINARY)

# Lint the code
lint:
	@echo "Linting $(SERVICE_NAME)..."
	cd $(SERVICE_DIR) && $(GOLINT) ./...

# Format the code
fmt:
	@echo "Formatting $(SERVICE_NAME)..."
	cd $(SERVICE_DIR) && $(GOFMT) .

# Create bin directory if it doesn't exist (needed for builds)
$(shell mkdir -p bin)

# Start Docker dependencies
docker-up:
	@echo "Starting Docker services..."
	$(DOCKER_COMPOSE) -f docker-compose.yaml up -d

# Stop Docker dependencies
docker-down:
	@echo "Stopping Docker services..."
	$(DOCKER_COMPOSE) -f docker-compose.yaml down

# Show help
help:
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build the $(SERVICE_NAME) application binary"
	@echo "  clean        Remove the $(SERVICE_NAME) application binary"
	@echo "  run          Build and run the $(SERVICE_NAME) application"
	@echo "  lint         Lint the Go code for $(SERVICE_NAME)"
	@echo "  fmt          Format the Go code for $(SERVICE_NAME)"
	@echo "  docker-up    Start required Docker services in the background"
	@echo "  docker-down  Stop required Docker services"
	@echo "  help         Show this help message"
	@echo ""

# --- Removed multi-service logic and related targets --- 
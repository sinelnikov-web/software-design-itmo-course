
.PHONY: build run test lint

# Variables
BINARY_NAME=gocli-cli
BIN_DIR=.

# Default target
all: build

# Build the application
build:
	go build -buildvcs=false -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/gocli

# Run the application
run: build
	$(BIN_DIR)/$(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run


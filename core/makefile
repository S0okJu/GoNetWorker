# Variables
APP_NAME := main
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")
VERSION := $(shell git describe --tags --always --dirty)
CMD_DIR := ./cmd

# Commands
build:
	@echo "Building the application..."
	@go build -o $(APP_NAME) $(CMD_DIR)

run:
	@echo "Running the application..."
	@go run $(CMD_DIR)

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f $(APP_NAME)

fmt:
	@echo "Formatting Go code..."
	@go fmt ./...

vet:
	@echo "Running Go vet..."
	@go vet ./...

install:
	@echo "Installing the application..."
	@go install $(CMD_DIR)

lint:
	@echo "Linting Go code..."
	@golangci-lint run ./...

# Default target
.PHONY: all
all: fmt vet build

# Build and add version information
build-version:
	@echo "Building application with version info..."
	@go build -ldflags "-X main.Version=$(VERSION)" -o $(APP_NAME) $(CMD_DIR)

# Targets to make your life easier
.PHONY: build run test clean fmt vet install lint
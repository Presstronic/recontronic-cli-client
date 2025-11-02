# Makefile for Recontronic CLI
# Build automation for Go project

# Variables
BINARY_NAME=recon-cli
MAIN_FILE=main.go
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GORUN=$(GOCMD) run

# Directories
PKG_DIR=./pkg/...
CMD_DIR=./cmd/...
TEST_DIR=./test/...

# Build output
BUILD_DIR=build
DIST_DIR=dist

.PHONY: all build install clean test test-verbose test-coverage fmt vet lint run help deps tidy check build-all build-linux build-darwin build-windows

# Default target
all: build

## help: Display this help message
help:
	@echo "Recontronic CLI - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@echo "  build           Build the binary for current platform"
	@echo "  install         Build and install to GOPATH/bin"
	@echo "  run             Run without building"
	@echo "  clean           Remove build artifacts"
	@echo "  test            Run all tests"
	@echo "  test-verbose    Run tests with verbose output"
	@echo "  test-coverage   Run tests with coverage report"
	@echo "  fmt             Format all Go source files"
	@echo "  vet             Run go vet on all packages"
	@echo "  lint            Run golangci-lint (if installed)"
	@echo "  check           Run fmt, vet, and test"
	@echo "  deps            Download dependencies"
	@echo "  tidy            Tidy and verify dependencies"
	@echo "  build-all       Build for all platforms"
	@echo "  build-linux     Build for Linux (amd64)"
	@echo "  build-darwin    Build for macOS (amd64 and arm64)"
	@echo "  build-windows   Build for Windows (amd64)"
	@echo "  help            Display this help message"
	@echo ""
	@echo "Version: $(VERSION)"

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BINARY_NAME)"

## install: Build and install to GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

## run: Run the CLI without building
run:
	$(GORUN) $(MAIN_FILE)

## clean: Remove build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	@echo "Clean complete"

## test: Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-verbose: Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	$(GOTEST) -v -count=1 ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## fmt: Format all Go source files
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Format complete"

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet complete"

## lint: Run golangci-lint (if installed)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  brew install golangci-lint  # macOS"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## check: Run fmt, vet, and test
check: fmt vet test
	@echo "All checks passed"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOGET) -v ./...
	@echo "Dependencies downloaded"

## tidy: Tidy and verify dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "Dependencies tidied"

## build-all: Build for all platforms
build-all: build-linux build-darwin build-windows
	@echo "All platform builds complete"

## build-linux: Build for Linux (amd64)
build-linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@echo "Linux build complete: $(DIST_DIR)/$(BINARY_NAME)-linux-amd64"

## build-darwin: Build for macOS (amd64 and arm64)
build-darwin:
	@echo "Building for macOS (amd64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	@echo "macOS builds complete"

## build-windows: Build for Windows (amd64)
build-windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(DIST_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	@echo "Windows build complete: $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe"

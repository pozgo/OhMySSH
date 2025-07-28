# OhMySSH Makefile
# Simple build system for SSH manager

.PHONY: all build build-all clean test help install

# Variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
APP_NAME = OhMySSH
BUILD_DIR = build
CMD_DIR = cmd/ohmyssh

# Default target
all: build

# Help target
help:
	@echo "OhMySSH Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build     - Build for current platform (stored in root directory)"
	@echo "  build-all - Build for all platforms (stored in build/ directory)"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  install   - Install to local system"
	@echo "  help      - Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION   - Build version (default: git tag or 'dev')"

# Build for current platform (for testing) - stored in root directory
build:
	@echo "Building $(APP_NAME) for current platform..."
	go build -o $(APP_NAME) $(CMD_DIR)/main.go
	@echo "✓ Built: ./$(APP_NAME) (ready for testing)"

# Build for all platforms - stored in build/ directory
build-all:
	@echo "Building $(APP_NAME) for all platforms..."
	@./build.sh $(VERSION)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests completed."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(APP_NAME)
	rm -rf $(BUILD_DIR)
	@echo "Clean completed."

# Install to local system
install: build
	@echo "Installing $(APP_NAME)..."
	sudo cp $(APP_NAME) /usr/local/bin/
	@echo "$(APP_NAME) installed to /usr/local/bin/"

# Show build information
info:
	@echo "Build Information:"
	@echo "  App Name: $(APP_NAME)"
	@echo "  Version:  $(VERSION)"
	@echo "  Go Version: $(shell go version)"
	@echo "  Build Dir: $(BUILD_DIR)"

# Run the application (for development)
run: build
	@echo "Running $(APP_NAME)..."
	@./$(APP_NAME)
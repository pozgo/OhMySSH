#!/bin/bash

# OhMySSH Build Script
# Builds for all supported platforms and architectures

set -e

# Configuration
APP_NAME="OhMySSH"
VERSION=${1:-"dev"}
BUILD_DIR="build"
CMD_DIR="cmd/ohmyssh"

# Create build directory
mkdir -p "$BUILD_DIR"

# Function to build for a specific platform
build_platform() {
    local os=$1
    local arch=$2
    local ext=$3
    
    echo "Building ${APP_NAME} for ${os}/${arch}..."
    
    local output_name="${APP_NAME}-${os}-${arch}${ext}"
    
    GOOS="$os" GOARCH="$arch" go build \
        -o "$BUILD_DIR/$output_name" \
        "$CMD_DIR/main.go"
    
    echo "✓ Built: $BUILD_DIR/$output_name"
}

echo "Building $APP_NAME v$VERSION for all platforms..."
echo ""

# Build for all platforms and architectures
build_platform "linux" "amd64" ""
build_platform "linux" "arm64" ""
build_platform "windows" "amd64" ".exe"
build_platform "windows" "arm64" ".exe"
build_platform "darwin" "amd64" ""
build_platform "darwin" "arm64" ""

echo ""
echo "All builds completed successfully!"
echo ""
echo "Built binaries in $BUILD_DIR/:"
ls -la "$BUILD_DIR/$APP_NAME"*
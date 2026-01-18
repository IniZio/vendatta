#!/bin/bash

set -e

REPO="IniZio/nexus"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $OS in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            echo "Unsupported OS: $OS"
            exit 1
            ;;
    esac

    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac

    BINARY_NAME="nexus-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="${BINARY_NAME}.exe"
    fi

    echo "$BINARY_NAME"
}

# Get latest release from GitHub API
get_latest_release() {
    curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Download and install
install_binary() {
    BINARY_NAME=$1
    TAG=$2

    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$TAG/$BINARY_NAME"

    echo "Downloading $BINARY_NAME from $DOWNLOAD_URL"

    if command -v curl >/dev/null 2>&1; then
        curl -L -o "/tmp/nexus" "$DOWNLOAD_URL"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "/tmp/nexus" "$DOWNLOAD_URL"
    else
        echo "Neither curl nor wget found. Please install one of them."
        exit 1
    fi

    chmod +x "/tmp/nexus"

    mkdir -p "$INSTALL_DIR"
    mv "/tmp/nexus" "$INSTALL_DIR/nexus"

    echo "Vendetta $TAG installed successfully to $INSTALL_DIR/nexus"
    echo "Run 'nexus --help' to get started"
}

# Build and install from source with version info
install_from_source() {
    TAG=$1
    BINARY_NAME=$(detect_platform)

    echo "Building nexus from source..."

    # Get version and date
    VERSION=${TAG#v}
    BUILDDATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    LDFLAGS="-X main.version=$VERSION -X main.buildDate=$BUILDDATE"

    # Build
    GOOS=$(uname -s | tr '[:upper:]' '[:lower:]')
    GOARCH=$(uname -m)
    case $GOARCH in
        x86_64) GOARCH="amd64" ;;
        aarch64|arm64) GOARCH="arm64" ;;
    esac

    TEMP_DIR=$(mktemp -d)
    cd "$TEMP_DIR"

    # Clone if repo not present, otherwise use existing
    if [ -d "/tmp/nexus-repo" ]; then
        cp -r /tmp/nexus-repo/* .
    else
        git clone --depth 1 --branch "$TAG" "https://github.com/$REPO.git" .
    fi

    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "$LDFLAGS" -o "/tmp/nexus" ./cmd/nexus

    mkdir -p "$INSTALL_DIR"
    mv "/tmp/nexus" "$INSTALL_DIR/nexus"
    chmod +x "$INSTALL_DIR/nexus"

    cd - > /dev/null
    rm -rf "$TEMP_DIR"

    echo "Vendetta $TAG installed successfully to $INSTALL_DIR/nexus"
    echo "Run 'nexus --help' to get started"
}

main() {
    echo "Installing nexus..."

    TAG=$(get_latest_release)
    BINARY_NAME=$(detect_platform)

    if [ -z "$TAG" ]; then
        echo "Failed to get latest release, building from source..."
        install_from_source "main"
        exit 0
    fi

    echo "Latest release: $TAG"
    echo "Platform: $BINARY_NAME"

    # Try to download pre-built binary first
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$TAG/$BINARY_NAME"
    if curl -sL --fail -o /dev/null "$DOWNLOAD_URL" 2>/dev/null; then
        install_binary "$BINARY_NAME" "$TAG"
    else
        echo "Pre-built binary not available for $BINARY_NAME, building from source..."
        install_from_source "$TAG"
    fi
}

main "$@"

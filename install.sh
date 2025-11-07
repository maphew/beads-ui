#!/bin/bash

# beady installer script
# Installs the latest beady binary for your platform

set -e

REPO="maphew/beady"
BINARY_NAME="beady"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $OS in
    linux)
        OS="Linux"
        ;;
    darwin)
        OS="Darwin"
        ;;
    mingw*|msys*|cygwin*)
        OS="Windows"
        BINARY_NAME="beady.exe"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

case $ARCH in
    x86_64|amd64)
        ARCH="x86_64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    i386|i686)
        ARCH="i386"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Get latest release info
echo "Fetching latest release information..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest")
VERSION=$(echo "$LATEST_RELEASE" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo "Failed to get latest version"
    exit 1
fi

echo "Latest version: $VERSION"

# Find the appropriate asset
ASSET_NAME="${BINARY_NAME}_${VERSION#v}_${OS}_${ARCH}"
if [ "$OS" = "Windows" ]; then
    DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | grep "browser_download_url.*${ASSET_NAME}\.zip" | cut -d '"' -f 4)
    EXT="zip"
else
    DOWNLOAD_URL=$(echo "$LATEST_RELEASE" | grep "browser_download_url.*${ASSET_NAME}\.tar\.gz" | cut -d '"' -f 4)
    EXT="tar.gz"
fi

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Could not find download URL for $ASSET_NAME"
    exit 1
fi

echo "Downloading $DOWNLOAD_URL..."

# Download and extract
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

if [ "$EXT" = "zip" ]; then
    curl -L -o "beady.$EXT" "$DOWNLOAD_URL"
    unzip "beady.$EXT"
else
    curl -L -o "beady.$EXT" "$DOWNLOAD_URL"
    tar -xzf "beady.$EXT"
fi

# Install to ~/bin or /usr/local/bin
INSTALL_DIR="${INSTALL_DIR:-$HOME/bin}"
if [ ! -d "$INSTALL_DIR" ]; then
    INSTALL_DIR="/usr/local/bin"
    if [ ! -w "$INSTALL_DIR" ]; then
        echo "Installing to $HOME/bin (you may need to add it to PATH)"
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
    fi
fi

echo "Installing to $INSTALL_DIR/$BINARY_NAME"
mv "$BINARY_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Clean up
cd /
rm -rf "$TEMP_DIR"

echo "Installation complete! Run '$BINARY_NAME --help' to get started."
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "Note: You may need to add $INSTALL_DIR to your PATH"
fi

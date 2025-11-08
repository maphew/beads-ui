#!/bin/bash

# beady installer script
# Installs the latest beady binary for your platform
# Usage: ./install.sh [uninstall [--yes]]

set -e

REPO="maphew/beady"
BINARY_NAME="beady"

# Uninstall function
uninstall() {
    local unattended=${1:-}

    # Search for beady in PATH (deduplicated)
    local found_paths=()
    local seen_paths=()
    IFS=':' read -ra path_dirs <<< "$PATH"
    for dir in "${path_dirs[@]}"; do
        # Skip empty directories
        [ -z "$dir" ] && continue
        local full_path="$dir/$BINARY_NAME"
        # Check if we've already processed this exact path
        if [[ ! " ${seen_paths[@]} " =~ " ${full_path} " ]]; then
            seen_paths+=("$full_path")
            if [ -f "$full_path" ] || [ -f "$dir/$BINARY_NAME.exe" ]; then
                found_paths+=("$full_path")
            fi
        fi
    done

    # Also check common install locations
    local common_locations=(
        "$HOME/.local/bin/$BINARY_NAME"
        "$HOME/bin/$BINARY_NAME"
        "/usr/local/bin/$BINARY_NAME"
    )
    for loc in "${common_locations[@]}"; do
        if [ -f "$loc" ] && [[ ! " ${found_paths[@]} " =~ " ${loc} " ]]; then
            found_paths+=("$loc")
        fi
    done

    if [ ${#found_paths[@]} -eq 0 ]; then
        echo "No beady installation found in PATH or common locations"
        exit 1
    fi

    if [ ${#found_paths[@]} -eq 1 ]; then
        echo "Found beady at: ${found_paths[0]}"
    else
        echo "Found multiple beady installations:"
        for i in "${!found_paths[@]}"; do
            echo "  $((i + 1)). ${found_paths[$i]}"
        done
    fi

    local target_path="${found_paths[0]}"

    # If multiple found and not unattended, let user choose
    if [ ${#found_paths[@]} -gt 1 ] && [ "$unattended" != "--yes" ]; then
        read -p "Remove all found installations? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Aborted"
            exit 0
        fi
    elif [ "$unattended" != "--yes" ]; then
        read -p "Remove $target_path? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Aborted"
            exit 0
        fi
    fi

    # Remove all found installations
    local removed_count=0
    for path in "${found_paths[@]}"; do
        if [ -f "$path" ]; then
            rm -f "$path"
            echo "Removed: $path"
            ((removed_count++))
        fi
    done

    if [ $removed_count -gt 0 ]; then
        echo "Uninstallation complete! Removed $removed_count beady installation(s)."
    else
        echo "No files were removed"
        exit 1
    fi
}

# Handle uninstall command
if [ "$1" = "uninstall" ]; then
    uninstall "$2"
    exit 0
elif [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "beady installer script"
    echo ""
    echo "Usage:"
    echo "  ./install.sh                  Install the latest beady release"
    echo "  ./install.sh uninstall        Uninstall beady (interactive)"
    echo "  ./install.sh uninstall --yes  Uninstall beady (unattended)"
    echo "  ./install.sh --help           Show this help message"
    exit 0
elif [ -n "$1" ]; then
    echo "Unknown option: $1"
    echo "Use ./install.sh --help for usage information"
    exit 1
fi

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

# Install to ~/.local/bin (XDG standard) or /usr/local/bin
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
if [ ! -d "$INSTALL_DIR" ]; then
    INSTALL_DIR="/usr/local/bin"
    if [ ! -w "$INSTALL_DIR" ]; then
        echo "Installing to $HOME/.local/bin (you may need to add it to PATH)"
        INSTALL_DIR="$HOME/.local/bin"
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

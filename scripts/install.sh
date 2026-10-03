#!/bin/bash
# Universal installer for wsl-dev-setup
# Usage: curl -fsSL https://github.com/jaircarrillo-lsx/wsl-dev-setup/releases/latest/download/install.sh | bash

set -euo pipefail

REPO="jaircarrillo-lsx/wsl-dev-setup"
BINARY="wsl-dev-setup"
INSTALL_DIR="${HOME}/.local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_error() { echo -e "${RED}[ERROR]${NC} $*"; }
log_success() { echo -e "${GREEN}[OK]${NC} $*"; }

detect_os_arch() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $OS in
        linux) OS="linux" ;;
        darwin) OS="darwin" ;;
        *) log_error "Unsupported OS: $OS"; exit 1 ;;
    esac

    case $ARCH in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) log_error "Unsupported architecture: $ARCH"; exit 1 ;;
    esac

    log_info "Detected: $OS/$ARCH"
}

download_latest() {
    log_info "Fetching latest release..."
    
    LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
    ASSET_NAME="${BINARY}_${OS}_${ARCH}"
    
    if [[ "$OS" == "windows" ]]; then
        ASSET_NAME="${ASSET_NAME}.zip"
    else
        ASSET_NAME="${ASSET_NAME}.tar.gz"
    fi

    # Build curl args with optional auth (check env var and command line arg)
    local github_token="${GITHUB_TOKEN:-${GITHUB_TOKEN_ARG:-}}"
    
    CURL_ARGS=(-fsSL -H "Accept: application/vnd.github+json")
    if [[ -n "$github_token" ]]; then
        CURL_ARGS+=(-H "Authorization: Bearer $github_token")
    fi
    CURL_ARGS+=("$LATEST_URL")

    # Retry with backoff
    local max_retries=3
    local retry=0
    local DOWNLOAD_URL=""

    while [[ $retry -lt $max_retries && -z "$DOWNLOAD_URL" ]]; do
        DOWNLOAD_URL=$(curl "${CURL_ARGS[@]}" 2>/dev/null | grep "browser_download_url" | grep "$ASSET_NAME" | head -1 | cut -d '"' -f 4)
        
        if [[ -z "$DOWNLOAD_URL" ]]; then
            retry=$((retry + 1))
            if [[ $retry -lt $max_retries ]]; then
                log_warn "Rate limited or asset not found, retrying in $((retry * 5))s... (attempt $retry/$max_retries)"
                sleep $((retry * 5))
            fi
        fi
    done

    if [[ -z "$DOWNLOAD_URL" ]]; then
        log_error "Could not find release asset for $OS/$ARCH after $max_retries attempts"
        log_info "Try setting GITHUB_TOKEN environment variable or use --token flag"
        exit 1
    fi

    log_info "Downloading from $DOWNLOAD_URL"
    
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    # Download with auth if available
    DOWNLOAD_ARGS=(-fsSL)
    if [[ -n "$github_token" ]]; then
        DOWNLOAD_ARGS+=(-H "Authorization: Bearer $github_token")
    fi
    DOWNLOAD_ARGS+=(-o "$TMP_DIR/$ASSET_NAME" "$DOWNLOAD_URL")

    curl "${DOWNLOAD_ARGS[@]}"

    if [[ "$ASSET_NAME" == *.zip ]]; then
        unzip -q "$TMP_DIR/$ASSET_NAME" -d "$TMP_DIR"
    else
        tar -xzf "$TMP_DIR/$ASSET_NAME" -C "$TMP_DIR"
    fi

    # The binary inside archive has platform suffix, find and rename it
    FOUND=$(find "$TMP_DIR" -name "wsl-dev-setup*" -type f ! -name "*.tar.gz" ! -name "*.zip" ! -name "install.sh" | head -1)
    if [[ -n "$FOUND" ]]; then
        cp "$FOUND" "$TMP_DIR/$BINARY"
    else
        log_error "Binary not found in archive"
        exit 1
    fi

    chmod +x "$TMP_DIR/$BINARY"
    INSTALL_PATH="$TMP_DIR/$BINARY"
}

install_binary() {
    mkdir -p "$INSTALL_DIR"
    
    if [[ -f "$INSTALL_DIR/$BINARY" ]]; then
        log_warn "Existing installation found at $INSTALL_DIR/$BINARY"
        log_info "Backing up..."
        mv "$INSTALL_DIR/$BINARY" "$INSTALL_DIR/$BINARY.bak.$(date +%s)"
    fi

    cp "$INSTALL_PATH" "$INSTALL_DIR/$BINARY"
    log_success "Installed to $INSTALL_DIR/$BINARY"

    # Add to PATH if needed
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        log_warn "$INSTALL_DIR is not in your PATH"
        log_info "Add this to your shell config (.bashrc, .zshrc, etc.):"
        echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
}

verify_install() {
    if command -v "$BINARY" >/dev/null 2>&1; then
        VERSION=$("$BINARY" version 2>/dev/null | head -1 || echo "unknown")
        log_success "Installation verified: $VERSION"
    else
        log_warn "Binary installed but not in PATH. Restart your shell or run:"
        echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
}

main() {
    # Parse arguments
    local GITHUB_TOKEN_ARG=""
    while [[ $# -gt 0 ]]; do
        case $1 in
            --token|-t)
                GITHUB_TOKEN_ARG="$2"
                shift 2
                ;;
            --help|-h)
                echo "Usage: $0 [--token TOKEN]"
                echo "  --token, -t  GitHub token for API authentication (or set GITHUB_TOKEN env var)"
                exit 0
                ;;
            *)
                log_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    export GITHUB_TOKEN_ARG

    echo "=================================="
    echo "  wsl-dev-setup installer"
    echo "=================================="
    echo

    detect_os_arch
    download_latest
    install_binary
    verify_install

    echo
    log_success "Installation complete!"
    echo
    echo "Next steps:"
    echo "  1. Restart your terminal or run: exec zsh"
    echo "  2. Run: wsl-dev-setup install"
    echo "  3. Run: p10k configure  (to customize prompt)"
}

main "$@"
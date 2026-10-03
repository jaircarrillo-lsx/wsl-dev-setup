# wsl-dev-setup

> One-command development environment for WSL, Termux, Docker, macOS, Windows, and Linux.

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![Release](https://img.shields.io/github/v/release/lunar-debian/wsl-dev-setup)](https://github.com/lunar-debian/wsl-dev-setup/releases)
[![License](https://img.shields.io/github/license/lunar-debian/wsl-dev-setup)](LICENSE)
[![Build](https://github.com/lunar-debian/wsl-dev-setup/actions/workflows/build.yml/badge.svg)](https://github.com/lunar-debian/wsl-dev-setup/actions/workflows/build.yml)

Installs a complete, modern development environment with:

- **Shell**: Oh My Zsh + Powerlevel10k (lean/classic/rainbow themes)
- **Plugins**: zsh-autosuggestions, zsh-syntax-highlighting, zsh-completions, fzf-tab, zsh-vi-mode
- **Modern CLI Tools**: eza, bat, fd, ripgrep, zoxide, btop, delta, lazygit, fzf
- **Editor Config**: Git config with delta, useful aliases
- **Opencode**: Pre-configured theme, keybindings, shell integration
- **Fonts**: MesloLGS NF + JetBrainsMono NF (Nerd Fonts)

## Quick Start

### Universal Installer (recommended)

```bash
curl -fsSL https://github.com/lunar-debian/wsl-dev-setup/releases/latest/download/install.sh | bash
```

Then restart your terminal and run:

```bash
wsl-dev-setup install
```

### Package Managers

```bash
# Homebrew (macOS/Linux)
brew install lunar-debian/tap/wsl-dev-setup

# Go install
go install github.com/lunar-debian/wsl-dev-setup@latest

# Docker
docker run --rm -it -v $HOME:/root ghcr.io/lunar-debian/wsl-dev-setup install
```

### Direct Binary Download

Download from [GitHub Releases](https://github.com/lunar-debian/wsl-dev-setup/releases) for your platform:
- Linux (amd64/arm64)
- macOS (Intel/Apple Silicon)
- Windows (amd64)
- Android/Termux (arm64)

## Usage

```bash
# Install with default settings (lean theme)
wsl-dev-setup install

# Customize theme
wsl-dev-setup install --style=classic   # Full featured
wsl-dev-setup install --style=rainbow   # Colorful
wsl-dev-setup install --style=lean      # Minimal (default)

# Skip components
wsl-dev-setup install --skip-fonts
wsl-dev-setup install --skip-opencode
wsl-dev-setup install --skip-chsh
wsl-dev-setup install --skip-system-base
wsl-dev-setup install --skip-omz
wsl-dev-setup install --skip-cli-tools

# Dry run (see what would be done)
wsl-dev-setup install --dry-run

# Force overwrite (with backup)
wsl-dev-setup install --force

# Check environment
wsl-dev-setup doctor

# Uninstall (restores backups)
wsl-dev-setup uninstall

# Version info
wsl-dev-setup version
```

## Platform Support

| Platform | Status | Notes |
|----------|--------|-------|
| WSL (Debian/Ubuntu) | ✅ Full | Primary target |
| Termux (Android) | ✅ Full | ARM64 binary |
| Docker/Containers | ✅ Full | `--skip-chsh --skip-fonts` |
| macOS (Intel/ARM) | ✅ Full | Homebrew supported |
| Windows (native) | ✅ Full | `.exe` binary |
| Linux (generic) | ✅ Full | Any distro with apt/pkg/brew/pacman |

## What Gets Installed

### System Packages (requires sudo)
```
zsh, git, curl, wget, fzf, bat, eza, fd-find, ripgrep, zoxide,
btop, git-delta, lazygit, unzip, fontconfig, ca-certificates
```

### Oh My Zsh + Powerlevel10k
- Unattended installation
- Lean/Classic/Rainbow theme presets
- Instant prompt enabled

### Zsh Plugins (cloned to `$ZSH_CUSTOM/plugins`)
- zsh-autosuggestions
- zsh-syntax-highlighting
- zsh-completions
- fzf-tab
- zsh-vi-mode

### Modern CLI Tools (from GitHub Releases)
- **eza** - Modern ls replacement
- **bat** - Cat with syntax highlighting
- **fd** - Fast find alternative
- **ripgrep (rg)** - Fast grep
- **zoxide** - Smart cd
- **btop** - Better top
- **delta** - Better git diff
- **lazygit** - TUI for git
- **fzf** - Fuzzy finder

### Configuration Files
- `~/.zshrc` - Optimized with aliases, keybindings, functions
- `~/.p10k.zsh` - Powerlevel10k theme config
- `~/.gitconfig` - Git aliases, delta pager
- `~/.config/opencode/config.json` - Opencode theme & keybindings

## Docker / Dev Containers

### As base image
```dockerfile
FROM ghcr.io/lunar-debian/wsl-dev-setup:latest
# Your Dockerfile continues...
```

### In devcontainer.json
```json
{
  "postCreateCommand": "wsl-dev-setup install --skip-fonts --skip-chsh --style=lean"
}
```

### Interactive container
```bash
docker run --rm -it -v $HOME:/root ghcr.io/lunar-debian/wsl-dev-setup install
```

## Termux (Android)

```bash
# Install via binary
curl -fsSL https://github.com/lunar-debian/wsl-dev-setup/releases/latest/download/wsl-dev-setup_linux_arm64.tar.gz | tar -xz -C $PREFIX/bin

# Or use the universal installer
curl -fsSL https://github.com/lunar-debian/wsl-dev-setup/releases/latest/download/install.sh | bash

# Then run
wsl-dev-setup install --skip-fonts --skip-chsh
```

## Opencode Integration

The installer configures opencode with:
- Catppuccin Mocha theme
- Custom keybindings (leader: Ctrl+b)
- Shell integration for seamless terminal switching
- JetBrainsMono Nerd Font

## Uninstall

```bash
wsl-dev-setup uninstall
```

This restores all backed-up files (`.zshrc`, `.p10k.zsh`, `.gitconfig`, opencode config) and removes installed components.

## Development

```bash
# Build
go build -o wsl-dev-setup .

# Test
go test ./...

# Run locally
./wsl-dev-setup install --dry-run
```

## Architecture

```
wsl-dev-setup/
├── cmd/                    # CLI commands (cobra)
├── internal/
│   ├── config/             # Configuration parsing
│   ├── installer/          # Installation orchestrator
│   │   └── steps/          # Individual installation steps
│   └── utils/              # Helpers (detect, exec, fs, backup, logger)
├── scripts/
│   └── install.sh          # Universal installer
├── .github/workflows/      # CI/CD
├── Dockerfile              # Multi-stage build
├── .goreleaser.yaml        # Release configuration
└── main.go                 # Entry point
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./...`
5. Submit a PR

---

**Made with ❤️ for developers who want a beautiful, productive terminal everywhere.**
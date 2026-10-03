package installer

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
)

type CLIToolsStep struct {
	arch   string
	osName string
}

func NewCLIToolsStep() *CLIToolsStep {
	arch := runtime.GOARCH
	osName := "linux"
	if runtime.GOOS == "darwin" {
		osName = "darwin"
	} else if runtime.GOOS == "windows" {
		osName = "windows"
	}
	return &CLIToolsStep{
		arch:   arch,
		osName: osName,
	}
}

func (s *CLIToolsStep) Name() string {
	return "cli_tools"
}

type cliTool struct {
	name         string
	repo         string
	binaryName   string
	versionArg   string
	assetPattern string
}

func (s *CLIToolsStep) Check(ctx *Context) (bool, error) {
	tools := s.getTools()
	binDir := filepath.Join(ctx.Env.HomeDir, ".local", "bin")

	for _, tool := range tools {
		binary := tool.binaryName
		if binary == "" {
			binary = tool.name
		}
		path := filepath.Join(binDir, binary)
		if !utils.FileExists(path) && !utils.CommandExists(binary) {
			return false, nil
		}
	}
	return true, nil
}

func (s *CLIToolsStep) Run(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would install CLI tools from GitHub releases")
		return nil
	}

	if ctx.Config.SkipCLITools {
		ctx.Logger.Info("Skipping CLI tools installation")
		return nil
	}

	tools := s.getTools()
	binDir := filepath.Join(ctx.Env.HomeDir, ".local", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	for _, tool := range tools {
		if err := s.installTool(ctx, tool, binDir); err != nil {
			ctx.Logger.Warn("Failed to install tool", "tool", tool.name, "error", err)
		}
	}

	return nil
}

func (s *CLIToolsStep) Rollback(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would remove CLI tools")
		return nil
	}

	if ctx.Config.SkipCLITools {
		ctx.Logger.Info("Skipping CLI tools removal")
		return nil
	}

	tools := s.getTools()
	binDir := filepath.Join(ctx.Env.HomeDir, ".local", "bin")

	for _, tool := range tools {
		binary := tool.binaryName
		if binary == "" {
			binary = tool.name
		}
		path := filepath.Join(binDir, binary)
		if utils.FileExists(path) {
			ctx.Logger.Info("Removing tool", "tool", tool.name)
			os.Remove(path)
		}
	}
	return nil
}

func (s *CLIToolsStep) getTools() []cliTool {
	return []cliTool{
		{
			name:         "eza",
			repo:         "eza-community/eza",
			binaryName:   "eza",
			assetPattern: "eza_{version}_{os}_{arch}.tar.gz",
		},
		{
			name:         "bat",
			repo:         "sharkdp/bat",
			binaryName:   "bat",
			assetPattern: "bat-{version}-{os}-{arch}.tar.gz",
		},
		{
			name:         "fd",
			repo:         "sharkdp/fd",
			binaryName:   "fd",
			assetPattern: "fd-{version}-{os}-{arch}.tar.gz",
		},
		{
			name:         "ripgrep",
			repo:         "BurntSushi/ripgrep",
			binaryName:   "rg",
			assetPattern: "ripgrep-{version}-{os}-{arch}.tar.gz",
		},
		{
			name:         "zoxide",
			repo:         "ajeetdsouza/zoxide",
			binaryName:   "zoxide",
			assetPattern: "zoxide-{version}-{os}-{arch}.tar.gz",
		},
		{
			name:         "btop",
			repo:         "aristocratos/btop",
			binaryName:   "btop",
			assetPattern: "btop-{version}-{os}-{arch}.tar.gz",
		},
		{
			name:         "delta",
			repo:         "dandavison/delta",
			binaryName:   "delta",
			assetPattern: "delta-{version}-{os}-{arch}.tar.gz",
		},
		{
			name:         "lazygit",
			repo:         "jesseduffield/lazygit",
			binaryName:   "lazygit",
			assetPattern: "lazygit_{version}_{os}_{arch}.tar.gz",
		},
	}
}

func (s *CLIToolsStep) installTool(ctx *Context, tool cliTool, binDir string) error {
	binary := tool.binaryName
	if binary == "" {
		binary = tool.name
	}

	targetPath := filepath.Join(binDir, binary)
	if utils.FileExists(targetPath) {
		ctx.Logger.Debug("Tool already installed", "tool", tool.name)
		return nil
	}

	ctx.Logger.Info("Installing tool", "tool", tool.name)

	version, err := getLatestRelease(ctx, tool.repo)
	if err != nil {
		return err
	}

	assetName := strings.ReplaceAll(tool.assetPattern, "{version}", version)
	assetName = strings.ReplaceAll(assetName, "{os}", s.osName)
	assetName = strings.ReplaceAll(assetName, "{arch}", s.arch)
	if strings.Contains(assetName, "windows") && !strings.HasSuffix(assetName, ".zip") {
		assetName += ".zip"
	} else if !strings.HasSuffix(assetName, ".tar.gz") && !strings.HasSuffix(assetName, ".zip") {
		assetName += ".tar.gz"
	}

	downloadURL := "https://github.com/" + tool.repo + "/releases/download/" + version + "/" + assetName
	tmpDir := filepath.Join(os.TempDir(), "wsl-dev-setup-tools", tool.name)
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, assetName)
	if _, err := utils.RunCommand(context.Background(), ctx.Logger, "curl", "-fL", "-o", archivePath, downloadURL); err != nil {
		return err
	}

	var extractedPath string
	if strings.HasSuffix(archivePath, ".zip") {
		if _, err := utils.RunCommand(context.Background(), ctx.Logger, "unzip", "-q", "-o", archivePath, "-d", tmpDir); err != nil {
			return err
		}
		extractedPath = findBinary(tmpDir, binary)
	} else {
		if _, err := utils.RunCommand(context.Background(), ctx.Logger, "tar", "-xzf", archivePath, "-C", tmpDir); err != nil {
			return err
		}
		extractedPath = findBinary(tmpDir, binary)
	}

	if extractedPath == "" {
		ctx.Logger.Warn("Binary not found in archive", "tool", tool.name)
		return nil
	}

	if err := os.Rename(extractedPath, targetPath); err != nil {
		return err
	}

	if err := os.Chmod(targetPath, 0755); err != nil {
		return err
	}

	ctx.Logger.Info("Tool installed", "tool", tool.name, "path", targetPath)
	return nil
}

func getLatestRelease(ctx *Context, repo string) (string, error) {
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	result, err := utils.RunCommand(context.Background(), ctx.Logger, "curl", "-fsSL", url)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(result.Stdout, "\n") {
		if strings.Contains(line, "\"tag_name\"") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				return strings.TrimPrefix(parts[3], "v"), nil
			}
		}
	}
	return "", nil
}

func findBinary(root, binary string) string {
	var found string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Name() == binary && info.Mode()&0111 != 0 {
			found = path
		}
		return nil
	})
	return found
}
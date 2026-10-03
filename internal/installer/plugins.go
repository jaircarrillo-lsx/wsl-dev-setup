package installer

import (
	"context"
	"os"
	"path/filepath"

	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
)

type PluginsStep struct{}

func NewPluginsStep() *PluginsStep {
	return &PluginsStep{}
}

func (s *PluginsStep) Name() string {
	return "plugins"
}

func (s *PluginsStep) Check(ctx *Context) (bool, error) {
	zshCustom := os.Getenv("ZSH_CUSTOM")
	if zshCustom == "" {
		zshCustom = filepath.Join(ctx.Env.HomeDir, ".oh-my-zsh", "custom")
	}

	plugins := []string{
		"zsh-autosuggestions",
		"zsh-syntax-highlighting",
		"zsh-completions",
		"fzf-tab",
		"zsh-vi-mode",
	}

	pluginsDir := filepath.Join(zshCustom, "plugins")
	for _, p := range plugins {
		if !utils.DirExists(filepath.Join(pluginsDir, p)) {
			return false, nil
		}
	}
	return true, nil
}

func (s *PluginsStep) Run(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would install Zsh plugins")
		return nil
	}

	zshCustom := os.Getenv("ZSH_CUSTOM")
	if zshCustom == "" {
		zshCustom = filepath.Join(ctx.Env.HomeDir, ".oh-my-zsh", "custom")
	}

	pluginsDir := filepath.Join(zshCustom, "plugins")
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return err
	}

	plugins := map[string]string{
		"zsh-autosuggestions":  "https://github.com/zsh-users/zsh-autosuggestions",
		"zsh-syntax-highlighting": "https://github.com/zsh-users/zsh-syntax-highlighting",
		"zsh-completions":      "https://github.com/zsh-users/zsh-completions",
		"fzf-tab":              "https://github.com/Aloxaf/fzf-tab",
		"zsh-vi-mode":          "https://github.com/jeffreytse/zsh-vi-mode",
	}

	for name, url := range plugins {
		target := filepath.Join(pluginsDir, name)
		if utils.DirExists(target) {
			ctx.Logger.Info("Plugin already exists", "plugin", name)
			continue
		}

		ctx.Logger.Info("Installing plugin", "plugin", name)
		_, err := utils.RunCommand(context.Background(), ctx.Logger, "git", "clone", "--depth=1", url, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PluginsStep) Rollback(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would remove Zsh plugins")
		return nil
	}

	zshCustom := os.Getenv("ZSH_CUSTOM")
	if zshCustom == "" {
		zshCustom = filepath.Join(ctx.Env.HomeDir, ".oh-my-zsh", "custom")
	}

	pluginsDir := filepath.Join(zshCustom, "plugins")
	plugins := []string{
		"zsh-autosuggestions",
		"zsh-syntax-highlighting",
		"zsh-completions",
		"fzf-tab",
		"zsh-vi-mode",
	}

	for _, name := range plugins {
		target := filepath.Join(pluginsDir, name)
		if utils.DirExists(target) {
			ctx.Logger.Info("Removing plugin", "plugin", name)
			os.RemoveAll(target)
		}
	}
	return nil
}
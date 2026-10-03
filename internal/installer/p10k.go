package installer

import (
	"context"
	"os"
	"path/filepath"

	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
)

type P10kStep struct{}

func NewP10kStep() *P10kStep {
	return &P10kStep{}
}

func (s *P10kStep) Name() string {
	return "p10k"
}

func (s *P10kStep) Check(ctx *Context) (bool, error) {
	zshCustom := os.Getenv("ZSH_CUSTOM")
	if zshCustom == "" {
		zshCustom = filepath.Join(ctx.Env.HomeDir, ".oh-my-zsh", "custom")
	}
	p10kPath := filepath.Join(zshCustom, "themes", "powerlevel10k")
	return utils.DirExists(p10kPath), nil
}

func (s *P10kStep) Run(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would install Powerlevel10k theme")
		return nil
	}

	zshCustom := os.Getenv("ZSH_CUSTOM")
	if zshCustom == "" {
		zshCustom = filepath.Join(ctx.Env.HomeDir, ".oh-my-zsh", "custom")
	}

	p10kPath := filepath.Join(zshCustom, "themes", "powerlevel10k")
	if utils.DirExists(p10kPath) {
		ctx.Logger.Info("Powerlevel10k already installed")
		return nil
	}

	ctx.Logger.Info("Installing Powerlevel10k theme...")

	if err := os.MkdirAll(filepath.Dir(p10kPath), 0755); err != nil {
		return err
	}

	_, err := utils.RunCommand(context.Background(), ctx.Logger, "git", "clone", "--depth=1", "https://github.com/romkatv/powerlevel10k.git", p10kPath)
	return err
}

func (s *P10kStep) Rollback(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would remove Powerlevel10k")
		return nil
	}

	zshCustom := os.Getenv("ZSH_CUSTOM")
	if zshCustom == "" {
		zshCustom = filepath.Join(ctx.Env.HomeDir, ".oh-my-zsh", "custom")
	}

	p10kPath := filepath.Join(zshCustom, "themes", "powerlevel10k")
	if utils.DirExists(p10kPath) {
		ctx.Logger.Info("Removing Powerlevel10k")
		return os.RemoveAll(p10kPath)
	}
	return nil
}
package installer

import (
	"context"
	"os"
	"path/filepath"

	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
)

type OMZStep struct{}

func NewOMZStep() *OMZStep {
	return &OMZStep{}
}

func (s *OMZStep) Name() string {
	return "omz"
}

func (s *OMZStep) Check(ctx *Context) (bool, error) {
	omzPath := filepath.Join(ctx.Env.HomeDir, ".oh-my-zsh")
	return utils.DirExists(omzPath) && utils.FileExists(filepath.Join(omzPath, "oh-my-zsh.sh")), nil
}

func (s *OMZStep) Run(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would install Oh My Zsh")
		return nil
	}

	if ctx.Config.SkipOMZ {
		ctx.Logger.Info("Skipping Oh My Zsh installation")
		return nil
	}

	omzPath := filepath.Join(ctx.Env.HomeDir, ".oh-my-zsh")
	if utils.DirExists(omzPath) {
		ctx.Logger.Info("Oh My Zsh already installed")
		return nil
	}

	ctx.Logger.Info("Installing Oh My Zsh...")

	// Download and execute the install script
	scriptURL := "https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh"
	tmpScript := filepath.Join(os.TempDir(), "omz_install.sh")

	if _, err := utils.RunCommand(context.Background(), ctx.Logger, "curl", "-fsSL", "-o", tmpScript, scriptURL); err != nil {
		return err
	}
	defer os.Remove(tmpScript)

	_, err := utils.RunCommand(context.Background(), ctx.Logger, "sh", tmpScript)
	return err
}

func (s *OMZStep) Rollback(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would remove Oh My Zsh")
		return nil
	}

	omzPath := filepath.Join(ctx.Env.HomeDir, ".oh-my-zsh")
	if utils.DirExists(omzPath) {
		ctx.Logger.Info("Removing Oh My Zsh")
		return os.RemoveAll(omzPath)
	}
	return nil
}
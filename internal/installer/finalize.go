package installer

import (
	"context"
	"os"
	"path/filepath"

	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
)

type FinalizeStep struct{}

func NewFinalizeStep() *FinalizeStep {
	return &FinalizeStep{}
}

func (s *FinalizeStep) Name() string {
	return "finalize"
}

func (s *FinalizeStep) Check(ctx *Context) (bool, error) {
	if ctx.Config.SkipChsh {
		return true, nil
	}

	currentShell := os.Getenv("SHELL")
	zshPath, _ := utils.LookPath("zsh")
	return currentShell == zshPath, nil
}

func (s *FinalizeStep) Run(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would finalize installation (chsh, font cache, verification)")
		return nil
	}

	ctx.Logger.Info("Finalizing installation...")

	if !ctx.Config.SkipChsh {
		if err := s.changeShell(ctx); err != nil {
			ctx.Logger.Warn("Failed to change shell", "error", err)
		}
	}

	if !ctx.Config.SkipFonts && ctx.Env.HasDisplay {
		ctx.Logger.Info("Updating font cache...")
		utils.RunCommand(context.Background(), ctx.Logger, "fc-cache", "-f")
	}

	if err := s.verifyInstallation(ctx); err != nil {
		ctx.Logger.Warn("Verification had issues", "error", err)
	}

	return nil
}

func (s *FinalizeStep) Rollback(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would revert shell change")
		return nil
	}

	if !ctx.Config.SkipChsh {
		ctx.Logger.Info("Reverting shell to bash")
		utils.RunCommand(context.Background(), ctx.Logger, "chsh", "-s", "$(which bash)")
	}
	return nil
}

func (s *FinalizeStep) changeShell(ctx *Context) error {
	zshPath, err := utils.LookPath("zsh")
	if err != nil {
		return err
	}

	currentShell := os.Getenv("SHELL")
	if currentShell == zshPath {
		ctx.Logger.Info("Shell already set to zsh")
		return nil
	}

	ctx.Logger.Info("Changing default shell to zsh", "shell", zshPath)
	_, err = utils.RunCommand(context.Background(), ctx.Logger, "chsh", "-s", zshPath)
	return err
}

func (s *FinalizeStep) verifyInstallation(ctx *Context) error {
	ctx.Logger.Info("Verifying installation...")

	checks := map[string]string{
		"zsh":       "zsh --version",
		"git":       "git --version",
		"oh-my-zsh": "test -d ~/.oh-my-zsh",
		"p10k":      "test -d ~/.oh-my-zsh/custom/themes/powerlevel10k",
		"eza":       "eza --version",
		"bat":       "bat --version",
		"fd":        "fd --version",
		"rg":        "rg --version",
		"zoxide":    "zoxide --version",
		"btop":      "btop --version",
		"delta":     "delta --version",
		"lazygit":   "lazygit --version",
		"fzf":       "fzf --version",
	}

	allPassed := true
	for name, cmd := range checks {
		if err := s.runCheck(ctx, name, cmd); err != nil {
			ctx.Logger.Warn("Check failed", "check", name, "error", err)
			allPassed = false
		} else {
			ctx.Logger.Info("Check passed", "check", name)
		}
	}

	if !allPassed {
		ctx.Logger.Warn("Some checks failed - installation may be incomplete")
	} else {
		ctx.Logger.Info("All verification checks passed!")
	}

	return nil
}

func (s *FinalizeStep) runCheck(ctx *Context, name, cmdStr string) error {
	if cmdStr == "test -d ~/.oh-my-zsh" {
		home := ctx.Env.HomeDir
		if _, err := os.Stat(filepath.Join(home, ".oh-my-zsh")); err != nil {
			return err
		}
		return nil
	}
	if cmdStr == "test -d ~/.oh-my-zsh/custom/themes/powerlevel10k" {
		home := ctx.Env.HomeDir
		if _, err := os.Stat(filepath.Join(home, ".oh-my-zsh", "custom", "themes", "powerlevel10k")); err != nil {
			return err
		}
		return nil
	}

	_, err := utils.RunCommand(context.Background(), ctx.Logger, "sh", "-c", cmdStr)
	return err
}
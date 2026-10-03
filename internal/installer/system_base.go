package installer

import (
	"context"

	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
)

type SystemBaseStep struct{}

func NewSystemBaseStep() *SystemBaseStep {
	return &SystemBaseStep{}
}

func (s *SystemBaseStep) Name() string {
	return "system_base"
}

func (s *SystemBaseStep) Check(ctx *Context) (bool, error) {
	required := []string{"zsh", "git", "curl", "wget", "fzf"}
	for _, cmd := range required {
		if !utils.CommandExists(cmd) {
			return false, nil
		}
	}
	return true, nil
}

func (s *SystemBaseStep) Run(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would install system packages")
		return nil
	}

	if ctx.Config.SkipSystemBase {
		ctx.Logger.Info("Skipping system package installation")
		return nil
	}

	packages := []string{
		"zsh", "git", "curl", "wget", "fzf",
		"bat", "eza", "fd-find", "ripgrep", "zoxide",
		"btop", "git-delta", "lazygit", "unzip",
		"fontconfig", "ca-certificates", "gnupg",
	}

	ctx.Logger.Info("Installing system packages", "packages", packages)

	switch ctx.Env.PackageMgr {
	case "apt":
		if err := runAptInstall(ctx, packages); err != nil {
			return err
		}
	case "pkg":
		if err := runPkgInstall(ctx, packages); err != nil {
			return err
		}
	case "brew":
		if err := runBrewInstall(ctx, packages); err != nil {
			return err
		}
	case "pacman":
		if err := runPacmanInstall(ctx, packages); err != nil {
			return err
		}
	case "dnf":
		if err := runDnfInstall(ctx, packages); err != nil {
			return err
		}
	case "zypper":
		if err := runZypperInstall(ctx, packages); err != nil {
			return err
		}
	default:
		ctx.Logger.Warn("Unknown package manager, skipping system packages", "mgr", ctx.Env.PackageMgr)
	}

	return nil
}

func (s *SystemBaseStep) Rollback(ctx *Context) error {
	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would remove system packages")
		return nil
	}

	if ctx.Config.SkipSystemBase {
		ctx.Logger.Info("Skipping system package removal")
		return nil
	}

	packages := []string{
		"zsh", "git", "curl", "wget", "fzf",
		"bat", "eza", "fd-find", "ripgrep", "zoxide",
		"btop", "git-delta", "lazygit", "unzip",
		"fontconfig", "ca-certificates", "gnupg",
	}

	var removeCmd []string
	switch ctx.Env.PackageMgr {
	case "apt":
		removeCmd = []string{"apt", "remove", "--auto-remove", "-y"}
	case "pkg":
		removeCmd = []string{"pkg", "uninstall", "-y"}
	case "brew":
		removeCmd = []string{"brew", "uninstall"}
	case "pacman":
		removeCmd = []string{"pacman", "-Rns", "--noconfirm"}
	case "dnf":
		removeCmd = []string{"dnf", "remove", "-y"}
	case "zypper":
		removeCmd = []string{"zypper", "remove", "-y"}
	default:
		return nil
	}

	args := append(removeCmd, packages...)
	ctx.Logger.Info("Removing system packages", "packages", packages)
	_, err := utils.RunCommand(context.Background(), ctx.Logger, args[0], args[1:]...)
	return err
}

func runAptInstall(ctx *Context, packages []string) error {
	if _, err := utils.RunCommand(context.Background(), ctx.Logger, "apt", "update"); err != nil {
		return err
	}
	args := append([]string{"install", "-y", "--no-install-recommends"}, packages...)
	_, err := utils.RunCommand(context.Background(), ctx.Logger, "apt", args...)
	return err
}

func runPkgInstall(ctx *Context, packages []string) error {
	args := append([]string{"install", "-y"}, packages...)
	_, err := utils.RunCommand(context.Background(), ctx.Logger, "pkg", args...)
	return err
}

func runBrewInstall(ctx *Context, packages []string) error {
	args := append([]string{"install"}, packages...)
	_, err := utils.RunCommand(context.Background(), ctx.Logger, "brew", args...)
	return err
}

func runPacmanInstall(ctx *Context, packages []string) error {
	args := append([]string{"-S", "--noconfirm", "--needed"}, packages...)
	_, err := utils.RunCommand(context.Background(), ctx.Logger, "pacman", args...)
	return err
}

func runDnfInstall(ctx *Context, packages []string) error {
	args := append([]string{"install", "-y"}, packages...)
	_, err := utils.RunCommand(context.Background(), ctx.Logger, "dnf", args...)
	return err
}

func runZypperInstall(ctx *Context, packages []string) error {
	args := append([]string{"install", "-y"}, packages...)
	_, err := utils.RunCommand(context.Background(), ctx.Logger, "zypper", args...)
	return err
}
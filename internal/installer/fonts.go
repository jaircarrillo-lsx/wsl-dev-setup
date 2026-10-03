package installer

import (
	"context"
	"os"
	"path/filepath"

	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
)

type FontsStep struct{}

func NewFontsStep() *FontsStep {
	return &FontsStep{}
}

func (s *FontsStep) Name() string {
	return "fonts"
}

func (s *FontsStep) Check(ctx *Context) (bool, error) {
	if ctx.Config.SkipFonts {
		return true, nil
	}

	fontDir := filepath.Join(ctx.Env.HomeDir, ".local", "share", "fonts")
	mesloDir := filepath.Join(fontDir, "MesloLGS NF")
	jetbrainsDir := filepath.Join(fontDir, "JetBrainsMono NF")

	return utils.DirExists(mesloDir) && utils.DirExists(jetbrainsDir), nil
}

func (s *FontsStep) Run(ctx *Context) error {
	if ctx.Config.SkipFonts {
		ctx.Logger.Info("Skipping fonts installation")
		return nil
	}

	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would install Nerd Fonts (MesloLGS NF, JetBrainsMono NF)")
		return nil
	}

	fontDir := filepath.Join(ctx.Env.HomeDir, ".local", "share", "fonts")
	if err := os.MkdirAll(fontDir, 0755); err != nil {
		return err
	}

	ctx.Logger.Info("Downloading Nerd Fonts...")

	if err := downloadFont(ctx, "MesloLGS NF", "https://github.com/ryanoasis/nerd-fonts/releases/latest/download/Meslo.zip"); err != nil {
		return err
	}

	if err := downloadFont(ctx, "JetBrainsMono NF", "https://github.com/ryanoasis/nerd-fonts/releases/latest/download/JetBrainsMono.zip"); err != nil {
		return err
	}

	ctx.Logger.Info("Updating font cache...")
	_, err := utils.RunCommand(context.Background(), ctx.Logger, "fc-cache", "-f", fontDir)
	return err
}

func (s *FontsStep) Rollback(ctx *Context) error {
	if ctx.Config.SkipFonts || ctx.DryRun {
		return nil
	}

	fontDir := filepath.Join(ctx.Env.HomeDir, ".local", "share", "fonts")
	dirs := []string{
		filepath.Join(fontDir, "MesloLGS NF"),
		filepath.Join(fontDir, "JetBrainsMono NF"),
	}

	for _, d := range dirs {
		if utils.DirExists(d) {
			ctx.Logger.Info("Removing font directory", "dir", d)
			os.RemoveAll(d)
		}
	}

	utils.RunCommand(context.Background(), ctx.Logger, "fc-cache", "-f", fontDir)
	return nil
}

func downloadFont(ctx *Context, name, url string) error {
	fontDir := filepath.Join(ctx.Env.HomeDir, ".local", "share", "fonts")
	targetDir := filepath.Join(fontDir, name)

	if utils.DirExists(targetDir) {
		ctx.Logger.Info("Font already exists", "font", name)
		return nil
	}

	ctx.Logger.Info("Downloading font", "font", name)
	tmpDir := filepath.Join(os.TempDir(), "wsl-dev-setup-fonts", name)
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	zipPath := filepath.Join(tmpDir, name+".zip")
	if _, err := utils.RunCommand(context.Background(), ctx.Logger, "curl", "-fL", "-o", zipPath, url); err != nil {
		return err
	}

	if _, err := utils.RunCommand(context.Background(), ctx.Logger, "unzip", "-q", "-o", zipPath, "-d", targetDir); err != nil {
		return err
	}

	ctx.Logger.Info("Font installed", "font", name, "path", targetDir)
	return nil
}
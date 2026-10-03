package installer

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
)

type OpencodeStep struct{}

func NewOpencodeStep() *OpencodeStep {
	return &OpencodeStep{}
}

func (s *OpencodeStep) Name() string {
	return "opencode"
}

func (s *OpencodeStep) Check(ctx *Context) (bool, error) {
	if ctx.Config.SkipOpencode {
		return true, nil
	}

	configPath := filepath.Join(ctx.Env.ConfigDir, "opencode", "config.json")
	return utils.FileExists(configPath), nil
}

func (s *OpencodeStep) Run(ctx *Context) error {
	if ctx.Config.SkipOpencode {
		ctx.Logger.Info("Skipping opencode configuration")
		return nil
	}

	if ctx.DryRun {
		ctx.Logger.Info("[DRY RUN] Would write opencode configuration")
		return nil
	}

	ctx.Logger.Info("Writing opencode configuration...")

	configDir := filepath.Join(ctx.Env.ConfigDir, "opencode")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.json")
	ctx.Backup.Backup(configPath)

	config := OpencodeConfig{
		Theme: "catppuccin-mocha",
		Keybindings: Keybindings{
			Leader: "ctrl+b",
			Commands: map[string]string{
				"new-session":    "n",
				"close-session":  "q",
				"next-pane":      "Tab",
				"prev-pane":      "Shift+Tab",
				"split-horizontal": "h",
				"split-vertical":   "v",
			},
		},
		ShellIntegration: true,
		Editor: EditorConfig{
			Command: "nvim",
			Args:    []string{},
		},
		Terminal: TerminalConfig{
			FontFamily: "JetBrainsMono Nerd Font",
			FontSize:   14,
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return utils.WriteFileAtomic(configPath, string(data))
}

func (s *OpencodeStep) Rollback(ctx *Context) error {
	if ctx.Config.SkipOpencode || ctx.DryRun {
		return nil
	}

	configPath := filepath.Join(ctx.Env.ConfigDir, "opencode", "config.json")
	ctx.Backup.Backup(configPath)
	ctx.Backup.RestoreAll()
	return nil
}

type OpencodeConfig struct {
	Theme            string            `json:"theme"`
	Keybindings      Keybindings       `json:"keybindings"`
	ShellIntegration bool              `json:"shellIntegration"`
	Editor           EditorConfig      `json:"editor"`
	Terminal         TerminalConfig    `json:"terminal"`
	Plugins          []string          `json:"plugins,omitempty"`
	Settings         map[string]any    `json:"settings,omitempty"`
}

type Keybindings struct {
	Leader   string            `json:"leader"`
	Commands map[string]string `json:"commands"`
}

type EditorConfig struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type TerminalConfig struct {
	FontFamily string `json:"fontFamily"`
	FontSize   int    `json:"fontSize"`
}
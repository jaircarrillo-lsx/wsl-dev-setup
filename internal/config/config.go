package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type InstallConfig struct {
	Style           string
	SkipFonts       bool
	SkipOpencode    bool
	SkipChsh        bool
	SkipSystemBase  bool
	SkipOMZ         bool
	SkipCLITools    bool
	Force           bool
	DryRun          bool
	NonInteractive  bool
	NoColor         bool
	HomeDir         string
	ConfigDir       string
}

type UninstallConfig struct {
	Force    bool
	NoColor  bool
	HomeDir  string
	ConfigDir string
}

func LoadInstallConfig(cmd *cobra.Command) InstallConfig {
	viper.BindPFlags(cmd.Flags())
	return InstallConfig{
		Style:           viper.GetString("style"),
		SkipFonts:       viper.GetBool("skip-fonts"),
		SkipOpencode:    viper.GetBool("skip-opencode"),
		SkipChsh:        viper.GetBool("skip-chsh"),
		SkipSystemBase:  viper.GetBool("skip-system-base"),
		SkipOMZ:         viper.GetBool("skip-omz"),
		SkipCLITools:    viper.GetBool("skip-cli-tools"),
		Force:           viper.GetBool("force"),
		DryRun:          viper.GetBool("dry-run"),
		NonInteractive:  viper.GetBool("non-interactive"),
		NoColor:         viper.GetBool("no-color"),
		HomeDir:         viper.GetString("home"),
		ConfigDir:       viper.GetString("config-dir"),
	}
}

func LoadUninstallConfig(cmd *cobra.Command) UninstallConfig {
	viper.BindPFlags(cmd.Flags())
	return UninstallConfig{
		Force:     viper.GetBool("force"),
		NoColor:   viper.GetBool("no-color"),
		HomeDir:   viper.GetString("home"),
		ConfigDir: viper.GetString("config-dir"),
	}
}
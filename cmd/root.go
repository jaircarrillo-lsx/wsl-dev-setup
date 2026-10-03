package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Version string
	Commit  string
	Date    string

	rootCmd = &cobra.Command{
		Use:     "wsl-dev-setup",
		Short:   "One-command dev environment for WSL/Termux/Docker: Oh My Zsh + Powerlevel10k + modern CLI tools + opencode",
		Long:    `wsl-dev-setup installs a complete development environment with Oh My Zsh, Powerlevel10k, modern CLI tools (eza, bat, fd, ripgrep, zoxide, btop, delta, lazygit, fzf), and opencode configuration.`,
		Version: Version,
	}

	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install the complete development environment",
		RunE:  runInstall,
	}

	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall and restore backups",
		RunE:  runUninstall,
	}

	doctorCmd = &cobra.Command{
		Use:   "doctor",
		Short: "Check environment and installation status",
		RunE:  runDoctor,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run:   runVersion,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(installCmd, uninstallCmd, doctorCmd, versionCmd)

	installCmd.Flags().String("style", "lean", "Powerlevel10k style: lean, classic, rainbow")
	installCmd.Flags().Bool("skip-fonts", false, "Skip Nerd Fonts installation")
	installCmd.Flags().Bool("skip-opencode", false, "Skip opencode configuration")
	installCmd.Flags().Bool("skip-chsh", false, "Skip changing default shell")
	installCmd.Flags().Bool("skip-system-base", false, "Skip system package installation (apt/pkg/brew)")
	installCmd.Flags().Bool("skip-omz", false, "Skip Oh My Zsh installation")
	installCmd.Flags().Bool("skip-cli-tools", false, "Skip CLI tools installation")
	installCmd.Flags().Bool("force", false, "Overwrite existing configs (with backup)")
	installCmd.Flags().Bool("dry-run", false, "Show what would be done without executing")
	installCmd.Flags().Bool("non-interactive", false, "Run without prompts (for CI)")
	installCmd.Flags().Bool("no-color", false, "Disable colored output")
	installCmd.Flags().String("home", "", "Override $HOME directory")
	installCmd.Flags().String("config-dir", "", "Override $XDG_CONFIG_HOME directory")

	uninstallCmd.Flags().Bool("force", false, "Force uninstall without confirmation")
}
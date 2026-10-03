package cmd

import (
	"fmt"

	"github.com/lunar-debian/wsl-dev-setup/internal/config"
	"github.com/lunar-debian/wsl-dev-setup/internal/installer"
	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
	"github.com/spf13/cobra"
)

func runInstall(cmd *cobra.Command, args []string) error {
	cfg := config.LoadInstallConfig(cmd)

	logger := utils.NewLogger(cfg.NoColor)

	if cfg.DryRun {
		logger.Info("DRY RUN mode - no changes will be made")
	}

	detect := utils.DetectEnvironment(cfg.HomeDir, cfg.ConfigDir)
	logger.Debug("Environment detected", "os", detect.OS, "arch", detect.Arch, "distro", detect.Distro, "wsl", detect.IsWSL, "termux", detect.IsTermux, "container", detect.IsContainer)

	inst := installer.New(logger, detect, cfg)
	return inst.Run()
}

func runUninstall(cmd *cobra.Command, args []string) error {
	cfg := config.LoadUninstallConfig(cmd)

	logger := utils.NewLogger(cfg.NoColor)

	detect := utils.DetectEnvironment("", "")

	if !cfg.Force {
		fmt.Print("This will restore backups and remove installed components. Continue? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			logger.Info("Uninstall cancelled")
			return nil
		}
	}

	inst := installer.New(logger, detect, config.InstallConfig{Force: cfg.Force})
	return inst.Uninstall()
}
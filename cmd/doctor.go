package cmd

import (
	"fmt"

	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
	"github.com/spf13/cobra"
)

func runDoctor(cmd *cobra.Command, args []string) error {
	logger := utils.NewLogger(false)

	detect := utils.DetectEnvironment("", "")

	logger.Info("Environment Detection")
	logger.Info("=====================")
	fmt.Printf("OS:           %s\n", detect.OS)
	fmt.Printf("Arch:         %s\n", detect.Arch)
	fmt.Printf("Distro:       %s\n", detect.Distro)
	fmt.Printf("WSL:          %v\n", detect.IsWSL)
	fmt.Printf("Termux:       %v\n", detect.IsTermux)
	fmt.Printf("Container:    %v\n", detect.IsContainer)
	fmt.Printf("Shell:        %s\n", detect.Shell)
	fmt.Printf("User:         %s\n", detect.User)
	fmt.Printf("Home:         %s\n", detect.HomeDir)
	fmt.Printf("Config Dir:   %s\n", detect.ConfigDir)
	fmt.Printf("Has apt:      %v\n", detect.HasApt)
	fmt.Printf("Has pkg:      %v\n", detect.HasPkg)
	fmt.Printf("Has brew:     %v\n", detect.HasBrew)
	fmt.Printf("Has systemd:  %v\n", detect.HasSystemd)
	fmt.Printf("Has display:  %v\n", detect.HasDisplay)

	logger.Info("\nTool Availability")
	logger.Info("==================")
	tools := []string{"zsh", "git", "curl", "wget", "fzf", "bat", "eza", "fd", "rg", "zoxide", "btop", "delta", "lazygit", "fc-cache"}
	for _, t := range tools {
		path, _ := utils.LookPath(t)
		status := "✓"
		if path == "" {
			status = "✗"
		}
		fmt.Printf("  %s %-12s %s\n", status, t, path)
	}

	return nil
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("wsl-dev-setup version %s\n", Version)
	fmt.Printf("commit: %s\n", Commit)
	fmt.Printf("date:   %s\n", Date)
}
package main

import (
	"fmt"
	"os"

	"github.com/lunar-debian/wsl-dev-setup/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
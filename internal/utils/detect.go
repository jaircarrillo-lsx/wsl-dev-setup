package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Environment struct {
	OS           string
	Arch         string
	Distro       string
	IsWSL        bool
	IsTermux     bool
	IsContainer  bool
	HasApt       bool
	HasPkg       bool
	HasBrew      bool
	HasChoco     bool
	HasPacman    bool
	HasDnf       bool
	HasZypper    bool
	Shell        string
	User         string
	HomeDir      string
	ConfigDir    string
	HasSystemd   bool
	HasDisplay   bool
	PackageMgr   string
}

func DetectEnvironment(homeOverride, configOverride string) *Environment {
	env := &Environment{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	env.User = os.Getenv("USER")
	if env.User == "" {
		env.User = os.Getenv("USERNAME")
	}

	env.HomeDir = homeOverride
	if env.HomeDir == "" {
		env.HomeDir = os.Getenv("HOME")
	}
	if env.HomeDir == "" {
		env.HomeDir = os.Getenv("USERPROFILE")
	}

	env.ConfigDir = configOverride
	if env.ConfigDir == "" {
		env.ConfigDir = os.Getenv("XDG_CONFIG_HOME")
	}
	if env.ConfigDir == "" && env.HomeDir != "" {
		env.ConfigDir = filepath.Join(env.HomeDir, ".config")
	}

	env.Shell = filepath.Base(os.Getenv("SHELL"))
	if env.Shell == "" {
		env.Shell = "bash"
	}

	env.IsWSL = isWSL()
	env.IsTermux = isTermux()
	env.IsContainer = isContainer()
	env.HasSystemd = hasSystemd()
	env.HasDisplay = hasDisplay()

	env.Distro = detectDistro()
	env.HasApt = hasCommand("apt")
	env.HasPkg = hasCommand("pkg")
	env.HasBrew = hasCommand("brew")
	env.HasChoco = hasCommand("choco")
	env.HasPacman = hasCommand("pacman")
	env.HasDnf = hasCommand("dnf")
	env.HasZypper = hasCommand("zypper")

	switch {
	case env.HasApt:
		env.PackageMgr = "apt"
	case env.HasPkg:
		env.PackageMgr = "pkg"
	case env.HasBrew:
		env.PackageMgr = "brew"
	case env.HasPacman:
		env.PackageMgr = "pacman"
	case env.HasDnf:
		env.PackageMgr = "dnf"
	case env.HasZypper:
		env.PackageMgr = "zypper"
	}

	return env
}

func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	content := strings.ToLower(string(data))
	return strings.Contains(content, "microsoft") || strings.Contains(content, "wsl")
}

func isTermux() bool {
	return os.Getenv("TERMUX_VERSION") != "" ||
		strings.HasPrefix(os.Getenv("PREFIX"), "/data/data/com.termux") ||
		strings.Contains(os.Getenv("HOME"), "com.termux")
}

func isContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	if _, err := os.Stat("/run/.containerenv"); err == nil {
		return true
	}
	if v := os.Getenv("container"); v != "" {
		return true
	}
	return false
}

func hasSystemd() bool {
	return os.Getenv("SYSTEMD_EXEC_PID") != "" ||
		hasCommand("systemctl")
}

func hasDisplay() bool {
	return os.Getenv("DISPLAY") != "" ||
		os.Getenv("WAYLAND_DISPLAY") != "" ||
		os.Getenv("TERM_PROGRAM") != ""
}

func detectDistro() string {
	if isTermux() {
		return "termux"
	}
	if _, err := os.Stat("/etc/os-release"); err == nil {
		data, _ := os.ReadFile("/etc/os-release")
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "ID=") {
				return strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
			}
		}
	}
	return "unknown"
}

func hasCommand(name string) bool {
	_, err := LookPath(name)
	return err == nil
}

func LookPath(name string) (string, error) {
	return executablePath(name)
}

func executablePath(name string) (string, error) {
	if filepath.IsAbs(name) {
		if info, err := os.Stat(name); err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
			return name, nil
		}
		return "", os.ErrNotExist
	}

	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		path := filepath.Join(dir, name)
		if info, err := os.Stat(path); err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
			return path, nil
		}
	}
	return "", os.ErrNotExist
}
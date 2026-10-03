package utils

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type ExecResult struct {
	Stdout string
	Stderr string
	ExitCode int
}

func RunCommand(ctx context.Context, logger *log.Logger, name string, args ...string) (ExecResult, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger.Debug("Running command", "cmd", name, "args", args)
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	result := ExecResult{
		Stdout:   strings.TrimSpace(stdout.String()),
		Stderr:   strings.TrimSpace(stderr.String()),
		ExitCode: cmd.ProcessState.ExitCode(),
	}

	if err != nil {
		logger.Debug("Command failed", "cmd", name, "duration", duration, "exit_code", result.ExitCode, "stderr", result.Stderr)
		return result, err
	}

	logger.Debug("Command succeeded", "cmd", name, "duration", duration)
	return result, nil
}

func RunCommandWithInput(ctx context.Context, logger *log.Logger, input string, name string, args ...string) (ExecResult, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	logger.Debug("Running command with input", "cmd", name, "args", args)
	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	result := ExecResult{
		Stdout:   strings.TrimSpace(stdout.String()),
		Stderr:   strings.TrimSpace(stderr.String()),
		ExitCode: cmd.ProcessState.ExitCode(),
	}

	if err != nil {
		logger.Debug("Command failed", "cmd", name, "duration", duration, "exit_code", result.ExitCode, "stderr", result.Stderr)
		return result, err
	}

	logger.Debug("Command succeeded", "cmd", name, "duration", duration)
	return result, nil
}

func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	return string(data), err
}

func WriteFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func WriteFileAtomic(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpPath := path + ".tmp." + time.Now().Format("20060102150405")
	if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func BackupFile(logger *log.Logger, path string) (string, error) {
	if !FileExists(path) {
		return "", nil
	}

	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(os.Getenv("HOME"), ".wsl-dev-setup", "backups", timestamp)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return "", err
	}

	relPath, _ := filepath.Rel(os.Getenv("HOME"), path)
	backupPath := filepath.Join(backupDir, relPath)

	if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return "", err
	}

	logger.Debug("Backed up file", "from", path, "to", backupPath)
	return backupPath, nil
}

func RestoreBackup(logger *log.Logger, backupPath, originalPath string) error {
	if !FileExists(backupPath) {
		return nil
	}

	data, err := os.ReadFile(backupPath)
	if err != nil {
		return err
	}

	if err := WriteFileAtomic(originalPath, string(data)); err != nil {
		return err
	}

	logger.Debug("Restored backup", "from", backupPath, "to", originalPath)
	return nil
}

func FileChecksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

func RenderTemplate(logger *log.Logger, tmplContent string, data interface{}) (string, error) {
	tmpl, err := template.New("config").Parse(tmplContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
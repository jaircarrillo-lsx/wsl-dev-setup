package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
)

type BackupManager struct {
	logger    *log.Logger
	backupDir string
	timestamp string
	files     []string
}

func NewBackupManager(logger *log.Logger) *BackupManager {
	timestamp := time.Now().Format("20060102_150405")
	backupDir := filepath.Join(os.Getenv("HOME"), ".wsl-dev-setup", "backups", timestamp)
	return &BackupManager{
		logger:    logger,
		backupDir: backupDir,
		timestamp: timestamp,
		files:     []string{},
	}
}

func (b *BackupManager) Backup(path string) error {
	if !FileExists(path) {
		return nil
	}

	relPath, _ := filepath.Rel(os.Getenv("HOME"), path)
	backupPath := filepath.Join(b.backupDir, relPath)

	if err := EnsureDir(filepath.Dir(backupPath)); err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return err
	}

	b.files = append(b.files, path)
	b.logger.Debug("Backed up", "file", path, "to", backupPath)
	return nil
}

func (b *BackupManager) RestoreAll() error {
	for _, originalPath := range b.files {
		relPath, _ := filepath.Rel(os.Getenv("HOME"), originalPath)
		backupPath := filepath.Join(b.backupDir, relPath)

		if !FileExists(backupPath) {
			continue
		}

		data, err := os.ReadFile(backupPath)
		if err != nil {
			b.logger.Warn("Failed to read backup", "file", backupPath, "error", err)
			continue
		}

		if err := WriteFileAtomic(originalPath, string(data)); err != nil {
			b.logger.Warn("Failed to restore backup", "file", originalPath, "error", err)
			continue
		}

		b.logger.Info("Restored", "file", originalPath)
	}
	return nil
}

func (b *BackupManager) ListBackups() ([]string, error) {
	baseDir := filepath.Join(os.Getenv("HOME"), ".wsl-dev-setup", "backups")
	if !DirExists(baseDir) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	var backups []string
	for _, e := range entries {
		if e.IsDir() {
			backups = append(backups, e.Name())
		}
	}
	return backups, nil
}

func (b *BackupManager) CleanOldBackups(keep int) error {
	backups, err := b.ListBackups()
	if err != nil {
		return err
	}

	if len(backups) <= keep {
		return nil
	}

	for _, backup := range backups[:len(backups)-keep] {
		path := filepath.Join(filepath.Dir(b.backupDir), backup)
		if err := os.RemoveAll(path); err != nil {
			b.logger.Warn("Failed to remove old backup", "backup", backup, "error", err)
		} else {
			b.logger.Info("Removed old backup", "backup", backup)
		}
	}
	return nil
}
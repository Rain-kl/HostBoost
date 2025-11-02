package hostsync

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// Default backup directory name
	defaultBackupDir = ".hostsync_backup"
	// Maximum number of backups to keep
	maxBackupFiles = 10
)

// createBackup creates a backup of the system hosts file
func (s *Syncer) createBackup() error {
	// Create backup directory if it doesn't exist
	backupDir := s.getBackupDir()
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("%w: %v", ErrBackupFailed, err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupFilename := fmt.Sprintf("hosts_backup_%s", timestamp)
	backupPath := filepath.Join(backupDir, backupFilename)

	// Read current system hosts file
	hostsContent, err := os.ReadFile(s.systemHostsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If hosts file doesn't exist, no need to backup
			return nil
		}
		return fmt.Errorf("%w: failed to read hosts file: %v", ErrBackupFailed, err)
	}

	// Write backup file
	if err := os.WriteFile(backupPath, hostsContent, 0644); err != nil {
		return fmt.Errorf("%w: failed to write backup file: %v", ErrBackupFailed, err)
	}

	// Clean up old backups
	if err := s.cleanupOldBackups(); err != nil {
		// Log warning but don't fail the backup operation
		fmt.Printf("Warning: failed to cleanup old backups: %v\n", err)
	}

	return nil
}

// getBackupDir returns the backup directory path
func (s *Syncer) getBackupDir() string {
	// Get the directory of hosts.json
	hostsDir := filepath.Dir(s.hostsJSONPath)
	return filepath.Join(hostsDir, defaultBackupDir)
}

// cleanupOldBackups removes old backup files, keeping only the most recent ones
func (s *Syncer) cleanupOldBackups() error {
	backupDir := s.getBackupDir()

	// Read all files in backup directory
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	// Filter backup files
	var backupFiles []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == "" {
			backupFiles = append(backupFiles, entry)
		}
	}

	// If we have more than max backups, delete oldest ones
	if len(backupFiles) > maxBackupFiles {
		// Get file info with modification times
		type fileWithTime struct {
			name    string
			modTime time.Time
		}
		var filesWithTime []fileWithTime

		for _, file := range backupFiles {
			info, err := file.Info()
			if err != nil {
				continue
			}
			filesWithTime = append(filesWithTime, fileWithTime{
				name:    file.Name(),
				modTime: info.ModTime(),
			})
		}

		// Sort by modification time (newest first)
		for i := 0; i < len(filesWithTime)-1; i++ {
			for j := i + 1; j < len(filesWithTime); j++ {
				if filesWithTime[i].modTime.Before(filesWithTime[j].modTime) {
					filesWithTime[i], filesWithTime[j] = filesWithTime[j], filesWithTime[i]
				}
			}
		}

		// Delete oldest files beyond maxBackupFiles
		for i := maxBackupFiles; i < len(filesWithTime); i++ {
			filePath := filepath.Join(backupDir, filesWithTime[i].name)
			if err := os.Remove(filePath); err != nil {
				fmt.Printf("Warning: failed to remove old backup %s: %v\n", filePath, err)
			}
		}
	}

	return nil
}

// ListBackups returns a list of available backup files
func (s *Syncer) ListBackups() ([]string, error) {
	backupDir := s.getBackupDir()

	// Check if backup directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read all files in backup directory
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, err
	}

	// Filter and collect backup file names
	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == "" {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}

// RestoreFromBackup restores the system hosts file from a backup
func (s *Syncer) RestoreFromBackup(backupName string) error {
	backupDir := s.getBackupDir()
	backupPath := filepath.Join(backupDir, backupName)

	// Check if backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file %s not found", backupName)
	}

	// Read backup content
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	// Write to system hosts file
	if err := os.WriteFile(s.systemHostsPath, backupContent, 0644); err != nil {
		if os.IsPermission(err) {
			return ErrPermissionDenied
		}
		return fmt.Errorf("failed to restore hosts file: %w", err)
	}

	return nil
}

// RestoreLatestBackup restores the system hosts file from the most recent backup
func (s *Syncer) RestoreLatestBackup() error {
	backups, err := s.ListBackups()
	if err != nil {
		return err
	}

	if len(backups) == 0 {
		return fmt.Errorf("no backup files found")
	}

	// Get file info to find the latest backup
	backupDir := s.getBackupDir()
	var latestBackup string
	var latestTime time.Time

	for _, backup := range backups {
		info, err := os.Stat(filepath.Join(backupDir, backup))
		if err != nil {
			continue
		}
		if latestBackup == "" || info.ModTime().After(latestTime) {
			latestBackup = backup
			latestTime = info.ModTime()
		}
	}

	if latestBackup == "" {
		return fmt.Errorf("no valid backup files found")
	}

	return s.RestoreFromBackup(latestBackup)
}

// DeleteAllBackups removes all backup files
func (s *Syncer) DeleteAllBackups() error {
	backupDir := s.getBackupDir()

	// Check if backup directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return nil // Nothing to delete
	}

	// Remove the entire backup directory
	if err := os.RemoveAll(backupDir); err != nil {
		return fmt.Errorf("failed to delete backup directory: %w", err)
	}

	return nil
}

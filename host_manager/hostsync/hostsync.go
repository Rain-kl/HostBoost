package hostsync

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
)

var (
	// ErrInvalidHostsFile indicates the hosts.json file is invalid or cannot be read
	ErrInvalidHostsFile = errors.New("invalid hosts file")
	// ErrPermissionDenied indicates insufficient permissions to modify system hosts file
	ErrPermissionDenied = errors.New("permission denied to modify system hosts file")
	// ErrBackupFailed indicates backup operation failed
	ErrBackupFailed = errors.New("backup operation failed")
)

// HostEntry represents a single host entry
type HostEntry struct {
	Domain string `json:"domain"`
	IP     string `json:"ip"`
	Type   string `json:"type"`
}

// Syncer handles synchronization between hosts.json and system hosts file
type Syncer struct {
	hostsJSONPath     string
	systemHostsPath   string
	backupEnabled     bool
	autoFlushDNSCache bool
}

// NewSyncer creates a new Syncer instance
func NewSyncer(hostsJSONPath string) *Syncer {
	return &Syncer{
		hostsJSONPath:     hostsJSONPath,
		systemHostsPath:   getSystemHostsPath(),
		backupEnabled:     true,
		autoFlushDNSCache: true, // Enable DNS cache flush by default
	}
}

// SetBackupEnabled enables or disables backup before sync
func (s *Syncer) SetBackupEnabled(enabled bool) {
	s.backupEnabled = enabled
}

// SetAutoFlushDNSCache enables or disables automatic DNS cache flush after sync
func (s *Syncer) SetAutoFlushDNSCache(enabled bool) {
	s.autoFlushDNSCache = enabled
}

// Sync reads hosts.json and synchronizes entries to system hosts file
func (s *Syncer) Sync() error {
	// Read hosts.json
	entries, err := s.readHostsJSON()
	if err != nil {
		return fmt.Errorf("failed to read hosts.json: %w", err)
	}

	// Create backup if enabled
	if s.backupEnabled {
		if err := s.createBackup(); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
	}

	// Read current system hosts file (only get non-managed lines)
	_, otherLines, err := s.parseSystemHosts()
	if err != nil {
		return fmt.Errorf("failed to parse system hosts file: %w", err)
	}

	// Write back to system hosts file (completely replace managed section with entries from hosts.json)
	if err := s.writeSystemHosts(entries, otherLines); err != nil {
		return fmt.Errorf("failed to write system hosts file: %w", err)
	}

	// Flush DNS cache after successful sync if enabled
	if s.autoFlushDNSCache {
		if err := s.FlushDNSCache(); err != nil {
			// Log warning but don't fail the sync operation
			fmt.Printf("Warning: failed to flush DNS cache: %v\n", err)
		}
	}

	return nil
}

// SyncFromJSON reads hosts.json and returns the entries without modifying system hosts
func (s *Syncer) SyncFromJSON() ([]HostEntry, error) {
	return s.readHostsJSON()
}

// readHostsJSON reads and parses the hosts.json file
func (s *Syncer) readHostsJSON() ([]HostEntry, error) {
	data, err := os.ReadFile(s.hostsJSONPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []HostEntry{}, nil // Return empty if file doesn't exist
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidHostsFile, err)
	}

	var entries []HostEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON format: %v", ErrInvalidHostsFile, err)
	}

	return entries, nil
}

// getSystemHostsPath returns the system hosts file path based on OS
func getSystemHostsPath() string {
	switch runtime.GOOS {
	case "windows":
		return "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "darwin", "linux":
		return "/etc/hosts"
	default:
		return "/etc/hosts"
	}
}

// ValidatePermissions checks if the current process has permissions to modify system hosts
func (s *Syncer) ValidatePermissions() error {
	// Try to open the file with write permissions
	file, err := os.OpenFile(s.systemHostsPath, os.O_RDWR, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return ErrPermissionDenied
		}
		return err
	}
	defer file.Close()

	return nil
}

// GetSystemHostsPath returns the current system hosts file path
func (s *Syncer) GetSystemHostsPath() string {
	return s.systemHostsPath
}

// GetHostsJSONPath returns the current hosts.json file path
func (s *Syncer) GetHostsJSONPath() string {
	return s.hostsJSONPath
}

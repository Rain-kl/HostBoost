package hostsync

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	// Marker comments to identify managed section
	managedSectionStart = "# === HostBoost Managed Section Start ==="
	managedSectionEnd   = "# === HostBoost Managed Section End ==="
)

// parseSystemHosts reads and parses the system hosts file
// Returns managed entries, other lines (comments, empty lines, non-managed entries), and error
func (s *Syncer) parseSystemHosts() ([]HostEntry, []string, error) {
	file, err := os.Open(s.systemHostsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If hosts file doesn't exist, return empty
			return []HostEntry{}, []string{}, nil
		}
		return nil, nil, err
	}
	defer file.Close()

	var entries []HostEntry
	var otherLines []string
	inManagedSection := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Check for managed section markers
		if trimmedLine == managedSectionStart {
			inManagedSection = true
			continue
		}
		if trimmedLine == managedSectionEnd {
			inManagedSection = false
			continue
		}

		// If in managed section, parse as host entry
		if inManagedSection {
			entry := parseHostLine(trimmedLine)
			if entry != nil {
				entries = append(entries, *entry)
			}
			continue
		}

		// Keep non-managed lines as-is
		otherLines = append(otherLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return entries, otherLines, nil
}

// parseHostLine parses a single host line into a HostEntry
// Returns nil if the line is invalid or a comment
func parseHostLine(line string) *HostEntry {
	line = strings.TrimSpace(line)

	// Skip empty lines and comments
	if line == "" || strings.HasPrefix(line, "#") {
		return nil
	}

	// Split by whitespace
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return nil
	}

	ip := fields[0]
	domain := fields[1]

	// Extract type from comment if present
	hostType := ""
	if len(fields) > 2 && strings.HasPrefix(fields[2], "#") {
		// Format: "127.0.0.1 example.com # type:cloudflare"
		for _, field := range fields[2:] {
			if strings.HasPrefix(field, "type:") {
				hostType = strings.TrimPrefix(field, "type:")
				break
			}
		}
	}

	return &HostEntry{
		Domain: domain,
		IP:     ip,
		Type:   hostType,
	}
}

// writeSystemHosts writes entries and other lines back to system hosts file
func (s *Syncer) writeSystemHosts(entries []HostEntry, otherLines []string) error {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "hosts-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath) // Clean up temp file

	writer := bufio.NewWriter(tempFile)

	// Write other lines (non-managed content)
	for _, line := range otherLines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			tempFile.Close()
			return fmt.Errorf("failed to write line: %w", err)
		}
	}

	// Write managed section
	if len(entries) > 0 {
		// Add a newline before managed section if otherLines doesn't end with empty line
		if len(otherLines) > 0 && strings.TrimSpace(otherLines[len(otherLines)-1]) != "" {
			writer.WriteString("\n")
		}

		// Write managed section start marker
		writer.WriteString(managedSectionStart + "\n")

		// Write managed entries
		for _, entry := range entries {
			line := formatHostEntry(entry)
			if _, err := writer.WriteString(line + "\n"); err != nil {
				tempFile.Close()
				return fmt.Errorf("failed to write entry: %w", err)
			}
		}

		// Write managed section end marker
		writer.WriteString(managedSectionEnd + "\n")
	}

	// Flush and close temp file
	if err := writer.Flush(); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	tempFile.Close()

	// Copy temp file to system hosts file
	if err := copyFile(tempPath, s.systemHostsPath); err != nil {
		return fmt.Errorf("failed to copy temp file to hosts file: %w", err)
	}

	return nil
}

// formatHostEntry formats a HostEntry as a hosts file line
func formatHostEntry(entry HostEntry) string {
	line := fmt.Sprintf("%-15s %s", entry.IP, entry.Domain)
	if entry.Type != "" {
		line += fmt.Sprintf(" # type:%s", entry.Type)
	}
	return line
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Write with appropriate permissions (644 for hosts file)
	if err := os.WriteFile(dst, input, 0644); err != nil {
		if os.IsPermission(err) {
			return ErrPermissionDenied
		}
		return err
	}

	return nil
}

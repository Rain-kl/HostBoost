package hostsync

import (
	"fmt"
	"os/exec"
	"runtime"
)

// FlushDNSCache flushes the DNS cache based on the operating system
func (s *Syncer) FlushDNSCache() error {
	var cmd *exec.Cmd
	var cmdDesc string

	switch runtime.GOOS {
	case "darwin": // macOS
		// macOS DNS cache flush command
		cmd = exec.Command("dscacheutil", "-flushcache")
		cmdDesc = "dscacheutil -flushcache"

		// For macOS 10.10.4+, also need to restart mDNSResponder
		if err := cmd.Run(); err == nil {
			killCmd := exec.Command("killall", "-HUP", "mDNSResponder")
			if killErr := killCmd.Run(); killErr != nil {
				// Log but don't fail if killall fails
				fmt.Printf("Warning: failed to restart mDNSResponder: %v\n", killErr)
			}
		} else {
			return fmt.Errorf("failed to flush DNS cache with %s: %w", cmdDesc, err)
		}
		return nil

	case "linux":
		// Linux: Try multiple common DNS cache services
		// systemd-resolved (Ubuntu 18.04+, many modern distros)
		cmd = exec.Command("systemd-resolve", "--flush-caches")
		if err := cmd.Run(); err == nil {
			return nil
		}

		// Try resolvectl (newer systemd)
		cmd = exec.Command("resolvectl", "flush-caches")
		if err := cmd.Run(); err == nil {
			return nil
		}

		// nscd (some older systems)
		cmd = exec.Command("nscd", "-i", "hosts")
		if err := cmd.Run(); err == nil {
			return nil
		}

		// dnsmasq (if using dnsmasq as DNS cache)
		cmd = exec.Command("killall", "-HUP", "dnsmasq")
		if err := cmd.Run(); err == nil {
			return nil
		}

		// If all failed, return a note (not an error as Linux might not use DNS cache)
		fmt.Println("Note: Could not flush DNS cache. Your system might not use a DNS cache service, or it may require manual restart.")
		return nil

	case "windows":
		// Windows DNS cache flush command
		cmd = exec.Command("ipconfig", "/flushdns")
		cmdDesc = "ipconfig /flushdns"
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to flush DNS cache with %s: %w", cmdDesc, err)
		}
		return nil

	default:
		// Unknown OS, just note it
		fmt.Printf("Note: DNS cache flush not implemented for OS: %s\n", runtime.GOOS)
		return nil
	}
}

// flushDNSCacheWithOutput flushes DNS cache and returns the output
func (s *Syncer) flushDNSCacheWithOutput() (string, error) {
	var cmd *exec.Cmd
	var cmdDesc string

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("dscacheutil", "-flushcache")
		cmdDesc = "dscacheutil -flushcache"
		output, err := cmd.CombinedOutput()
		if err != nil {
			return string(output), fmt.Errorf("failed to flush DNS cache with %s: %w", cmdDesc, err)
		}

		// Also restart mDNSResponder
		killCmd := exec.Command("killall", "-HUP", "mDNSResponder")
		killOutput, killErr := killCmd.CombinedOutput()
		if killErr != nil {
			return string(output) + "\n" + string(killOutput), fmt.Errorf("failed to restart mDNSResponder: %w", killErr)
		}

		return string(output) + "\nDNS cache flushed successfully", nil

	case "linux":
		// Try different commands
		commands := [][]string{
			{"systemd-resolve", "--flush-caches"},
			{"resolvectl", "flush-caches"},
			{"nscd", "-i", "hosts"},
		}

		for _, cmdArgs := range commands {
			cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
			output, err := cmd.CombinedOutput()
			if err == nil {
				return string(output) + "\nDNS cache flushed successfully", nil
			}
		}

		return "", fmt.Errorf("could not flush DNS cache - no supported DNS cache service found")

	case "windows":
		cmd = exec.Command("ipconfig", "/flushdns")
		cmdDesc = "ipconfig /flushdns"
		output, err := cmd.CombinedOutput()
		if err != nil {
			return string(output), fmt.Errorf("failed to flush DNS cache with %s: %w", cmdDesc, err)
		}
		return string(output), nil

	default:
		return "", fmt.Errorf("DNS cache flush not implemented for OS: %s", runtime.GOOS)
	}
}

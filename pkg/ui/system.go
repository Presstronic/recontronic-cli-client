package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/presstronic/recontronic-cli-client/pkg/config"
)

// ToolStatus represents the status of a recon tool
type ToolStatus struct {
	Name      string
	Installed bool
	Version   string
	Type      string // "built-in", "external"
}

// SystemStatus represents overall system health
type SystemStatus struct {
	Tools          []ToolStatus
	ServerStatus   string // "connected", "disconnected", "unknown"
	AuthStatus     string // "authenticated", "not_authenticated"
	StorageUsed    int64
	ToolsAvailable int
	ToolsTotal     int
}

// GetSystemStatus checks tool availability and system health
func GetSystemStatus(cfg *config.Config) (*SystemStatus, error) {
	status := &SystemStatus{
		Tools: []ToolStatus{
			{Name: "crt.sh", Installed: true, Version: "built-in", Type: "built-in"},
			checkExternalTool("subfinder"),
			checkExternalTool("amass"),
			checkExternalTool("assetfinder"),
			checkExternalTool("httpx"),
			checkExternalTool("nuclei"),
		},
	}

	// Count available tools
	for _, tool := range status.Tools {
		status.ToolsTotal++
		if tool.Installed {
			status.ToolsAvailable++
		}
	}

	// Check server status
	if cfg != nil && cfg.Server != "" {
		status.ServerStatus = "disconnected" // Would need actual check
	} else {
		status.ServerStatus = "not_configured"
	}

	// Check auth status
	if cfg != nil && cfg.APIKey != "" {
		status.AuthStatus = "authenticated"
	} else {
		status.AuthStatus = "not_authenticated"
	}

	// Get storage used
	stats, err := GatherStats()
	if err == nil {
		status.StorageUsed = stats.StorageUsed
	}

	return status, nil
}

// checkExternalTool checks if an external tool is installed
func checkExternalTool(name string) ToolStatus {
	status := ToolStatus{
		Name:      name,
		Installed: false,
		Type:      "external",
	}

	path, err := exec.LookPath(name)
	if err != nil {
		return status
	}

	status.Installed = true

	// Try to get version
	versionCmd := exec.Command(name, "-version")
	if output, err := versionCmd.CombinedOutput(); err == nil {
		version := strings.TrimSpace(string(output))
		// Extract version number (usually first line)
		lines := strings.Split(version, "\n")
		if len(lines) > 0 {
			status.Version = strings.TrimSpace(lines[0])
			// Limit version string length
			if len(status.Version) > 30 {
				status.Version = status.Version[:30]
			}
		}
	} else {
		// Try -h or --version
		versionCmd = exec.Command(name, "--version")
		if output, err := versionCmd.CombinedOutput(); err == nil {
			version := strings.TrimSpace(string(output))
			lines := strings.Split(version, "\n")
			if len(lines) > 0 {
				status.Version = strings.TrimSpace(lines[0])
				if len(status.Version) > 30 {
					status.Version = status.Version[:30]
				}
			}
		} else {
			// If no version info, just show path
			status.Version = "installed"
		}
	}

	// Fallback: just show it's installed
	if status.Version == "" {
		status.Version = fmt.Sprintf("found at %s", path)
		if len(status.Version) > 30 {
			status.Version = "installed"
		}
	}

	return status
}

package recon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/config"
)

// StorageFormat represents the output format for results
type StorageFormat int

const (
	FormatJSON StorageFormat = iota
	FormatText
)

// GetResultsDir returns the base results directory
func GetResultsDir() (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "results"), nil
}

// GetDomainResultsDir returns the results directory for a specific domain
func GetDomainResultsDir(domain string) (string, error) {
	resultsDir, err := GetResultsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(resultsDir, domain), nil
}

// EnsureDomainResultsDir creates the results directory for a domain if it doesn't exist
func EnsureDomainResultsDir(domain string) error {
	domainDir, err := GetDomainResultsDir(domain)
	if err != nil {
		return err
	}

	// Create directory with 0700 permissions (owner only)
	if err := os.MkdirAll(domainDir, 0700); err != nil {
		return fmt.Errorf("failed to create results directory: %w", err)
	}

	return nil
}

// SaveResults saves results to a file in the domain's results directory
func SaveResults(domain, toolName string, data interface{}, format StorageFormat) (string, error) {
	if err := EnsureDomainResultsDir(domain); err != nil {
		return "", err
	}

	domainDir, err := GetDomainResultsDir(domain)
	if err != nil {
		return "", err
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	var filename string

	switch format {
	case FormatJSON:
		filename = fmt.Sprintf("%s_%s.json", toolName, timestamp)
	case FormatText:
		filename = fmt.Sprintf("%s_%s.txt", toolName, timestamp)
	default:
		return "", fmt.Errorf("unsupported format: %d", format)
	}

	filePath := filepath.Join(domainDir, filename)

	// Marshal data based on format
	var fileData []byte
	switch format {
	case FormatJSON:
		fileData, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", fmt.Errorf("failed to marshal JSON: %w", err)
		}
	case FormatText:
		// Assume data is already a string or []byte
		switch v := data.(type) {
		case string:
			fileData = []byte(v)
		case []byte:
			fileData = v
		default:
			return "", fmt.Errorf("text format requires string or []byte data")
		}
	}

	// Write file with secure permissions
	if err := os.WriteFile(filePath, fileData, 0600); err != nil {
		return "", fmt.Errorf("failed to write results file: %w", err)
	}

	return filePath, nil
}

// LoadLatestResult loads the most recent result file for a tool
func LoadLatestResult(domain, toolName string, result interface{}) error {
	domainDir, err := GetDomainResultsDir(domain)
	if err != nil {
		return err
	}

	// Find latest file matching pattern
	pattern := filepath.Join(domainDir, fmt.Sprintf("%s_*.json", toolName))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to search for results: %w", err)
	}

	if len(matches) == 0 {
		return fmt.Errorf("no results found for %s on %s", toolName, domain)
	}

	// Get the latest file (files are timestamped, so last alphabetically is latest)
	latestFile := matches[len(matches)-1]

	// Read and unmarshal
	data, err := os.ReadFile(latestFile)
	if err != nil {
		return fmt.Errorf("failed to read results file: %w", err)
	}

	if err := json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("failed to unmarshal results: %w", err)
	}

	return nil
}

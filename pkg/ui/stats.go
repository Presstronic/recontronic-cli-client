package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/config"
)

// DashboardStats represents overall statistics
type DashboardStats struct {
	TotalDomains    int
	TotalSubdomains int
	TotalAlive      int
	ScansLast24h    int
	ScansLast7d     int
	StorageUsed     int64
	LastUpdated     time.Time
}

// SubdomainResult represents the structure of subdomain JSON files
type SubdomainResult struct {
	Domain      string      `json:"domain"`
	Timestamp   time.Time   `json:"timestamp"`
	TotalUnique int         `json:"total_unique"`
	TotalAlive  int         `json:"total_alive,omitempty"`
	Subdomains  []Subdomain `json:"subdomains"`
}

// Subdomain represents a single subdomain entry
type Subdomain struct {
	Name         string              `json:"name"`
	DiscoveredBy []string            `json:"discovered_by"`
	Verified     *VerificationResult `json:"verified,omitempty"`
}

// VerificationResult represents verification data
type VerificationResult struct {
	Status string `json:"status"` // "alive", "dead", "error"
}

// GatherStats collects statistics from the results directory
func GatherStats() (*DashboardStats, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	resultsDir := filepath.Join(configDir, "results")

	stats := &DashboardStats{
		LastUpdated: time.Now(),
	}

	// Check if results directory exists
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
		return stats, nil
	}

	// Count domains
	domains, err := os.ReadDir(resultsDir)
	if err != nil {
		return stats, nil
	}

	stats.TotalDomains = len(domains)

	// Analyze each domain's results
	for _, domain := range domains {
		if !domain.IsDir() {
			continue
		}

		domainPath := filepath.Join(resultsDir, domain.Name())
		files, err := os.ReadDir(domainPath)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			filePath := filepath.Join(domainPath, file.Name())

			// Get file size for storage stats
			info, err := os.Stat(filePath)
			if err == nil {
				stats.StorageUsed += info.Size()
			}

			// Parse subdomain JSON files
			if filepath.Ext(file.Name()) == ".json" &&
				len(file.Name()) > 11 &&
				file.Name()[:11] == "subdomains_" {

				data, err := os.ReadFile(filePath)
				if err != nil {
					continue
				}

				var result SubdomainResult
				if err := json.Unmarshal(data, &result); err != nil {
					continue
				}

				stats.TotalSubdomains += result.TotalUnique

				// Count alive subdomains
				for _, sub := range result.Subdomains {
					if sub.Verified != nil && sub.Verified.Status == "alive" {
						stats.TotalAlive++
					}
				}

				// Count scans by age
				age := time.Since(result.Timestamp)
				if age < 24*time.Hour {
					stats.ScansLast24h++
				}
				if age < 7*24*time.Hour {
					stats.ScansLast7d++
				}
			}
		}
	}

	return stats, nil
}

// FormatBytes formats bytes as human-readable size
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

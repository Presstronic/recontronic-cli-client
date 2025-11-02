package export

import (
	"os"
	"path/filepath"

	"github.com/presstronic/recontronic-cli-client/pkg/config"
	"github.com/presstronic/recontronic-cli-client/pkg/recon"
)

// ExportFormat represents the output format for exports
type ExportFormat string

const (
	FormatCSV      ExportFormat = "csv"
	FormatJSON     ExportFormat = "json"
	FormatMarkdown ExportFormat = "markdown"
)

// ExportOptions configures export behavior
type ExportOptions struct {
	Format     ExportFormat
	OutputPath string
	AliveOnly  bool
	DeadOnly   bool
	StatusCode int
	Source     string
}

// GetExportsDir returns the default exports directory
func GetExportsDir() (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	exportsDir := filepath.Join(configDir, "exports")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(exportsDir, 0700); err != nil {
		return "", err
	}

	return exportsDir, nil
}

// filterSubdomains applies export options to filter subdomains
func filterSubdomains(subdomains []recon.Subdomain, options ExportOptions) []recon.Subdomain {
	var filtered []recon.Subdomain

	for _, sub := range subdomains {
		// Apply filters
		if options.AliveOnly && (sub.Verified == nil || sub.Verified.Status != "alive") {
			continue
		}

		if options.DeadOnly && (sub.Verified == nil || sub.Verified.Status != "dead") {
			continue
		}

		if options.StatusCode != 0 {
			if sub.Verified == nil || sub.Verified.HTTP == nil || sub.Verified.HTTP.StatusCode != options.StatusCode {
				continue
			}
		}

		if options.Source != "" {
			found := false
			for _, source := range sub.DiscoveredBy {
				if source == options.Source {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		filtered = append(filtered, sub)
	}

	return filtered
}

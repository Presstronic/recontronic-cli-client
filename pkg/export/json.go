package export

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/presstronic/recontronic-cli-client/pkg/recon"
)

// ExportToJSON exports subdomain results to JSON format
func ExportToJSON(result *recon.SubdomainResults, options ExportOptions) (string, error) {
	filePath := options.OutputPath
	if filePath == "" {
		filePath = fmt.Sprintf("%s_subdomains.json", result.Domain)
	}

	// Filter subdomains based on options
	filtered := *result
	filtered.Subdomains = filterSubdomains(result.Subdomains, options)
	filtered.TotalUnique = len(filtered.Subdomains)

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to write JSON file: %w", err)
	}

	return filePath, nil
}

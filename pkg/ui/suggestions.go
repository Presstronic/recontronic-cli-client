package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/config"
)

// Suggestion represents an actionable suggestion for the user
type Suggestion struct {
	Message  string
	Action   string // Command to run
	Priority int    // 1=high, 2=medium, 3=low
}

// GenerateSuggestions analyzes state and suggests next actions
func GenerateSuggestions() ([]Suggestion, error) {
	var suggestions []Suggestion

	configDir, err := config.GetConfigDir()
	if err != nil {
		return suggestions, nil
	}

	resultsDir := filepath.Join(configDir, "results")

	// Check if results directory exists
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
		suggestions = append(suggestions, Suggestion{
			Message:  "No scans yet - try 'recon subdomain example.com'",
			Action:   "recon subdomain example.com",
			Priority: 2,
		})
		return suggestions, nil
	}

	// Analyze each domain's results
	domains, err := os.ReadDir(resultsDir)
	if err != nil {
		return suggestions, nil
	}

	for _, domain := range domains {
		if !domain.IsDir() {
			continue
		}

		domainName := domain.Name()
		domainPath := filepath.Join(resultsDir, domainName)
		files, err := os.ReadDir(domainPath)
		if err != nil {
			continue
		}

		var latestSubdomainFile string
		var latestSubdomainTime time.Time
		hasUnverified := false
		unverifiedCount := 0

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			filePath := filepath.Join(domainPath, file.Name())

			// Find subdomain files
			if filepath.Ext(file.Name()) == ".json" &&
				len(file.Name()) > 11 &&
				file.Name()[:11] == "subdomains_" {

				info, err := os.Stat(filePath)
				if err != nil {
					continue
				}

				if latestSubdomainFile == "" || info.ModTime().After(latestSubdomainTime) {
					latestSubdomainFile = filePath
					latestSubdomainTime = info.ModTime()
				}

				// Check if it has unverified subdomains
				data, err := os.ReadFile(filePath)
				if err != nil {
					continue
				}

				var result SubdomainResult
				if err := json.Unmarshal(data, &result); err != nil {
					continue
				}

				unverified := 0
				for _, sub := range result.Subdomains {
					if sub.Verified == nil {
						unverified++
					}
				}

				if unverified > 0 {
					hasUnverified = true
					unverifiedCount = unverified
				}

				// Check if scan is old (> 7 days)
				age := time.Since(result.Timestamp)
				if age > 7*24*time.Hour {
					suggestions = append(suggestions, Suggestion{
						Message:  fmt.Sprintf("%s not scanned in %dd - consider re-scanning", domainName, int(age.Hours()/24)),
						Action:   fmt.Sprintf("recon subdomain %s", domainName),
						Priority: 3,
					})
				}
			}
		}

		// Suggest verification if unverified results exist
		if hasUnverified {
			suggestions = append(suggestions, Suggestion{
				Message:  fmt.Sprintf("%s has %d unverified subdomains", domainName, unverifiedCount),
				Action:   fmt.Sprintf("recon verify %s", domainName),
				Priority: 1,
			})
		}
	}

	// Suggest installing missing tools
	systemStatus, err := GetSystemStatus(nil)
	if err == nil {
		missingTools := []string{}
		for _, tool := range systemStatus.Tools {
			if !tool.Installed && tool.Type == "external" {
				missingTools = append(missingTools, tool.Name)
			}
		}

		if len(missingTools) > 0 && len(missingTools) <= 3 {
			toolList := ""
			for i, tool := range missingTools {
				if i > 0 {
					toolList += ", "
				}
				toolList += tool
			}
			suggestions = append(suggestions, Suggestion{
				Message:  fmt.Sprintf("Install %s for better coverage", toolList),
				Action:   "",
				Priority: 3,
			})
		}
	}

	// Limit to top 5 suggestions
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return suggestions, nil
}

package recon

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ResultInfo represents metadata about a stored result file
type ResultInfo struct {
	Domain      string
	ToolName    string
	Timestamp   time.Time
	FilePath    string
	FileSize    int64
	TotalCount  int
	AliveCount  int
	DeadCount   int
	Verified    bool
	SourcesUsed []string
}

// QueryOptions configures result filtering
type QueryOptions struct {
	AliveOnly  bool
	DeadOnly   bool
	StatusCode int
	Source     string
}

// ListResults lists all stored results grouped by domain
func ListResults() (map[string][]ResultInfo, error) {
	resultsDir, err := GetResultsDir()
	if err != nil {
		return nil, err
	}

	// Check if results directory exists
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
		return make(map[string][]ResultInfo), nil
	}

	// Read all domain directories
	entries, err := os.ReadDir(resultsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read results directory: %w", err)
	}

	resultsByDomain := make(map[string][]ResultInfo)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		domain := entry.Name()
		domainResults, err := ListResultsForDomain(domain)
		if err != nil {
			continue
		}

		if len(domainResults) > 0 {
			resultsByDomain[domain] = domainResults
		}
	}

	return resultsByDomain, nil
}

// ListResultsForDomain lists all results for a specific domain
func ListResultsForDomain(domain string) ([]ResultInfo, error) {
	domainDir, err := GetDomainResultsDir(domain)
	if err != nil {
		return nil, err
	}

	// Check if directory exists
	if _, err := os.Stat(domainDir); os.IsNotExist(err) {
		return []ResultInfo{}, nil
	}

	// Find all JSON files
	pattern := filepath.Join(domainDir, "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search for results: %w", err)
	}

	var results []ResultInfo

	for _, filePath := range matches {
		// Parse filename to extract tool name and timestamp
		filename := filepath.Base(filePath)
		parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")

		if len(parts) < 3 {
			continue
		}

		toolName := parts[0]
		timestampStr := strings.Join(parts[1:], "_")

		// Parse timestamp
		timestamp, err := time.Parse("20060102_150405", timestampStr)
		if err != nil {
			continue
		}

		// Get file size
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		// Load file to get counts
		resultInfo := ResultInfo{
			Domain:    domain,
			ToolName:  toolName,
			Timestamp: timestamp,
			FilePath:  filePath,
			FileSize:  fileInfo.Size(),
		}

		// Parse file to get statistics
		if toolName == "subdomains" {
			var subResult SubdomainResults
			if err := loadJSONFile(filePath, &subResult); err == nil {
				resultInfo.TotalCount = subResult.TotalUnique
				resultInfo.SourcesUsed = subResult.SourcesUsed

				// Count alive and dead
				aliveCount := 0
				deadCount := 0
				hasVerified := false

				for _, sub := range subResult.Subdomains {
					if sub.Verified != nil {
						hasVerified = true
						if sub.Verified.Status == "alive" {
							aliveCount++
						} else if sub.Verified.Status == "dead" {
							deadCount++
						}
					}
				}

				resultInfo.AliveCount = aliveCount
				resultInfo.DeadCount = deadCount
				resultInfo.Verified = hasVerified
			}
		}

		results = append(results, resultInfo)
	}

	// Sort by timestamp (newest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.After(results[j].Timestamp)
	})

	return results, nil
}

// LoadSubdomainResult loads a specific subdomain result by timestamp
func LoadSubdomainResult(domain string, timestamp time.Time) (*SubdomainResults, error) {
	domainDir, err := GetDomainResultsDir(domain)
	if err != nil {
		return nil, err
	}

	// Build expected filename
	timestampStr := timestamp.Format("20060102_150405")
	filename := fmt.Sprintf("subdomains_%s.json", timestampStr)
	filePath := filepath.Join(domainDir, filename)

	var result SubdomainResults
	if err := loadJSONFile(filePath, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetLatestSubdomainResult loads the most recent subdomain scan for a domain
func GetLatestSubdomainResult(domain string) (*SubdomainResults, error) {
	results, err := ListResultsForDomain(domain)
	if err != nil {
		return nil, err
	}

	// Filter to only subdomain results
	var subdomainResults []ResultInfo
	for _, r := range results {
		if r.ToolName == "subdomains" {
			subdomainResults = append(subdomainResults, r)
		}
	}

	if len(subdomainResults) == 0 {
		return nil, fmt.Errorf("no subdomain results found for %s", domain)
	}

	// Load the latest (results are already sorted by timestamp descending)
	latest := subdomainResults[0]
	var result SubdomainResults
	if err := loadJSONFile(latest.FilePath, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// QuerySubdomains filters subdomains based on query options
func QuerySubdomains(domain string, options QueryOptions) ([]Subdomain, error) {
	result, err := GetLatestSubdomainResult(domain)
	if err != nil {
		return nil, err
	}

	var filtered []Subdomain

	for _, sub := range result.Subdomains {
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

	return filtered, nil
}

// loadJSONFile is a helper to load and unmarshal a JSON file
func loadJSONFile(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}

// FormatFileSize formats a file size in human-readable format
func FormatFileSize(bytes int64) string {
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

package recon

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// SubdomainResults represents the complete subdomain enumeration results
type SubdomainResults struct {
	Domain       string                 `json:"domain"`
	Timestamp    time.Time              `json:"timestamp"`
	SourcesUsed  []string               `json:"sources_used"`
	TotalUnique  int                    `json:"total_unique"`
	Subdomains   []Subdomain            `json:"subdomains"`
	Summary      map[string]int         `json:"summary"`
}

// Subdomain represents a single subdomain entry
type Subdomain struct {
	Name         string                 `json:"name"`
	DiscoveredBy []string               `json:"discovered_by"`
	FirstSeen    time.Time              `json:"first_seen"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// SubdomainSource interface for enumeration tools
type SubdomainSource interface {
	Name() string
	IsAvailable() bool
	Enumerate(domain string) ([]string, error)
}

// EnumerateSubdomains runs all available sources and aggregates results
func EnumerateSubdomains(domain string, sources []SubdomainSource) (*SubdomainResults, error) {
	results := &SubdomainResults{
		Domain:      domain,
		Timestamp:   time.Now(),
		SourcesUsed: []string{},
		Subdomains:  []Subdomain{},
		Summary:     make(map[string]int),
	}

	// Map to track which sources found each subdomain
	subdomainMap := make(map[string]*Subdomain)

	// Run each source
	for _, source := range sources {
		if !source.IsAvailable() {
			continue
		}

		sourceName := source.Name()
		results.SourcesUsed = append(results.SourcesUsed, sourceName)

		// Show progress
		fmt.Printf("Running %s... ", sourceName)
		startTime := time.Now()

		// Enumerate subdomains
		subdomains, err := source.Enumerate(domain)
		duration := time.Since(startTime)

		if err != nil {
			// Log error but continue with other sources
			fmt.Printf("✗ failed after %s: %v\n", duration.Round(time.Second), err)
			continue
		}

		fmt.Printf("✓ completed in %s\n", duration.Round(time.Second))

		// Clean the results
		subdomains = CleanDomains(subdomains)
		results.Summary[sourceName] = len(subdomains)

		// Merge into results
		for _, sub := range subdomains {
			if existing, found := subdomainMap[sub]; found {
				// Subdomain already found by another source
				existing.DiscoveredBy = append(existing.DiscoveredBy, sourceName)
			} else {
				// New subdomain
				subdomainMap[sub] = &Subdomain{
					Name:         sub,
					DiscoveredBy: []string{sourceName},
					FirstSeen:    time.Now(),
					Metadata:     make(map[string]interface{}),
				}
			}
		}
	}

	// Convert map to slice
	for _, subdomain := range subdomainMap {
		results.Subdomains = append(results.Subdomains, *subdomain)
	}

	// Sort subdomains
	sortedNames := make([]string, len(results.Subdomains))
	for i, sub := range results.Subdomains {
		sortedNames[i] = sub.Name
	}
	sortedNames = SortDomains(sortedNames)

	// Rebuild subdomains in sorted order
	sortedSubdomains := make([]Subdomain, len(sortedNames))
	for i, name := range sortedNames {
		sortedSubdomains[i] = *subdomainMap[name]
	}
	results.Subdomains = sortedSubdomains

	results.TotalUnique = len(results.Subdomains)

	return results, nil
}

// CrtShSource implements SubdomainSource for crt.sh certificate transparency
type CrtShSource struct{}

func (s *CrtShSource) Name() string {
	return "crt.sh"
}

func (s *CrtShSource) IsAvailable() bool {
	return true // Always available (API-based)
}

func (s *CrtShSource) Enumerate(domain string) ([]string, error) {
	// Query crt.sh API
	url := fmt.Sprintf("https://crt.sh/?q=%%.%s&output=json", domain)
	result, err := ExecuteWithTimeout("curl", 2*time.Minute, "-s", url)
	if err != nil {
		return nil, fmt.Errorf("crt.sh query failed: %w", err)
	}

	// Parse JSON response
	var crtResults []struct {
		NameValue string `json:"name_value"`
	}

	if err := json.Unmarshal([]byte(result.Stdout), &crtResults); err != nil {
		return nil, fmt.Errorf("failed to parse crt.sh response: %w", err)
	}

	// Extract subdomains
	var subdomains []string
	for _, entry := range crtResults {
		// name_value can contain multiple subdomains separated by newlines
		names := strings.Split(entry.NameValue, "\n")
		for _, name := range names {
			name = strings.TrimSpace(name)
			if name != "" {
				subdomains = append(subdomains, name)
			}
		}
	}

	return subdomains, nil
}

// SubfinderSource implements SubdomainSource for subfinder
type SubfinderSource struct{}

func (s *SubfinderSource) Name() string {
	return "subfinder"
}

func (s *SubfinderSource) IsAvailable() bool {
	return IsToolAvailable("subfinder")
}

func (s *SubfinderSource) Enumerate(domain string) ([]string, error) {
	// Run subfinder with JSON output
	result, err := ExecuteWithTimeout("subfinder", 5*time.Minute, "-d", domain, "-silent", "-json")
	if err != nil {
		return nil, fmt.Errorf("subfinder execution failed: %w", err)
	}

	// Parse JSON lines output
	var subdomains []string
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		var subfinderResult struct {
			Host string `json:"host"`
		}

		if err := json.Unmarshal([]byte(line), &subfinderResult); err != nil {
			// If JSON parsing fails, treat line as plain subdomain
			subdomains = append(subdomains, strings.TrimSpace(line))
			continue
		}

		if subfinderResult.Host != "" {
			subdomains = append(subdomains, subfinderResult.Host)
		}
	}

	return subdomains, nil
}

// AmassSource implements SubdomainSource for amass
type AmassSource struct{}

func (s *AmassSource) Name() string {
	return "amass"
}

func (s *AmassSource) IsAvailable() bool {
	return IsToolAvailable("amass")
}

func (s *AmassSource) Enumerate(domain string) ([]string, error) {
	// Run amass in passive mode
	result, err := ExecuteWithTimeout("amass", 10*time.Minute, "enum", "-passive", "-d", domain, "-nocolor")
	if err != nil {
		return nil, fmt.Errorf("amass execution failed: %w", err)
	}

	// Parse output (one subdomain per line)
	var subdomains []string
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			subdomains = append(subdomains, line)
		}
	}

	return subdomains, nil
}

// AssetfinderSource implements SubdomainSource for assetfinder
type AssetfinderSource struct{}

func (s *AssetfinderSource) Name() string {
	return "assetfinder"
}

func (s *AssetfinderSource) IsAvailable() bool {
	return IsToolAvailable("assetfinder")
}

func (s *AssetfinderSource) Enumerate(domain string) ([]string, error) {
	// Run assetfinder with subs-only flag
	result, err := ExecuteWithTimeout("assetfinder", 5*time.Minute, "--subs-only", domain)
	if err != nil {
		return nil, fmt.Errorf("assetfinder execution failed: %w", err)
	}

	// Parse output (one subdomain per line)
	var subdomains []string
	lines := strings.Split(strings.TrimSpace(result.Stdout), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			subdomains = append(subdomains, line)
		}
	}

	return subdomains, nil
}

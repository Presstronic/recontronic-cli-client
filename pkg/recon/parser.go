package recon

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateDomain checks if a domain is valid
func ValidateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	// Basic domain validation regex
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	return nil
}

// CleanDomains removes duplicates, wildcards, and sorts domains
func CleanDomains(domains []string) []string {
	// Remove wildcards first
	cleaned := RemoveWildcards(domains)

	// Deduplicate
	deduped := Deduplicate(cleaned)

	// Sort alphabetically
	return SortDomains(deduped)
}

// RemoveWildcards strips wildcard prefixes from domains
func RemoveWildcards(domains []string) []string {
	result := make([]string, 0, len(domains))

	for _, domain := range domains {
		// Remove leading *. wildcard
		cleaned := strings.TrimPrefix(domain, "*.")

		// Also handle other wildcard patterns
		cleaned = strings.TrimSpace(cleaned)

		if cleaned != "" {
			result = append(result, cleaned)
		}
	}

	return result
}

// Deduplicate removes duplicate entries (case-insensitive)
func Deduplicate(items []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(items))

	for _, item := range items {
		// Normalize to lowercase for comparison
		key := strings.ToLower(item)

		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return result
}

// SortDomains sorts domains alphabetically (case-insensitive)
func SortDomains(domains []string) []string {
	// Simple bubble sort for now (can optimize later if needed)
	result := make([]string, len(domains))
	copy(result, domains)

	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if strings.ToLower(result[i]) > strings.ToLower(result[j]) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// FilterByPattern filters domains matching a pattern
func FilterByPattern(domains []string, pattern string) []string {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return domains // Return original on error
	}

	result := make([]string, 0, len(domains))
	for _, domain := range domains {
		if regex.MatchString(domain) {
			result = append(result, domain)
		}
	}

	return result
}

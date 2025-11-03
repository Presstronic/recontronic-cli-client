package recon

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// WhoisInfo represents parsed WHOIS information for a domain
type WhoisInfo struct {
	Domain       string    `json:"domain"`
	Registrar    string    `json:"registrar,omitempty"`
	CreatedDate  string    `json:"created_date,omitempty"`
	UpdatedDate  string    `json:"updated_date,omitempty"`
	ExpiryDate   string    `json:"expiry_date,omitempty"`
	NameServers  []string  `json:"name_servers,omitempty"`
	Status       []string  `json:"status,omitempty"`
	RegistrarURL string    `json:"registrar_url,omitempty"`
	WhoisServer  string    `json:"whois_server,omitempty"`
	RawOutput    string    `json:"raw_output"`
	LookedUpAt   time.Time `json:"looked_up_at"`
}

// WhoisResults represents the complete WHOIS lookup results
type WhoisResults struct {
	Domain     string    `json:"domain"`
	Info       WhoisInfo `json:"info"`
	LookedUpAt time.Time `json:"looked_up_at"`
	Error      string    `json:"error,omitempty"`
}

// LookupWhois performs a WHOIS lookup for the given domain
func LookupWhois(ctx context.Context, domain string, timeout time.Duration) (*WhoisInfo, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Execute whois command
	cmd := exec.CommandContext(ctx, "whois", domain)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("whois command failed: %w", err)
	}

	rawOutput := string(output)

	// Parse the WHOIS output
	info := parseWhoisOutput(domain, rawOutput)
	info.LookedUpAt = time.Now()

	return &info, nil
}

// parseWhoisOutput parses raw WHOIS output into structured data
func parseWhoisOutput(domain, rawOutput string) WhoisInfo {
	info := WhoisInfo{
		Domain:      domain,
		RawOutput:   rawOutput,
		NameServers: []string{},
		Status:      []string{},
	}

	// Use maps to track seen nameservers and statuses for deduplication
	seenNameServers := make(map[string]bool)
	seenStatuses := make(map[string]bool)

	lines := strings.Split(rawOutput, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "%") || strings.HasPrefix(line, ">>>") {
			continue
		}

		// Split on first colon
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Parse based on key (case-insensitive matching)
		keyLower := strings.ToLower(key)

		switch {
		case strings.Contains(keyLower, "registrar") && !strings.Contains(keyLower, "whois") &&
			!strings.Contains(keyLower, "url") && !strings.Contains(keyLower, "iana") &&
			!strings.Contains(keyLower, "abuse") && info.Registrar == "":
			info.Registrar = value

		case strings.Contains(keyLower, "creation date") || strings.Contains(keyLower, "created"):
			if info.CreatedDate == "" {
				info.CreatedDate = value
			}

		case strings.Contains(keyLower, "updated date") || strings.Contains(keyLower, "last updated"):
			if info.UpdatedDate == "" {
				info.UpdatedDate = value
			}

		case strings.Contains(keyLower, "expiry date") || strings.Contains(keyLower, "expiration date") ||
			strings.Contains(keyLower, "registry expiry"):
			if info.ExpiryDate == "" {
				info.ExpiryDate = value
			}

		case strings.Contains(keyLower, "name server"):
			// Extract just the nameserver hostname
			ns := strings.Fields(value)
			if len(ns) > 0 {
				nsLower := strings.ToLower(ns[0])
				if !seenNameServers[nsLower] {
					seenNameServers[nsLower] = true
					info.NameServers = append(info.NameServers, nsLower)
				}
			}

		case strings.Contains(keyLower, "domain status") || key == "Status":
			// Extract status value
			statusParts := strings.Fields(value)
			if len(statusParts) > 0 {
				// Normalize status to lowercase for deduplication
				statusLower := strings.ToLower(statusParts[0])
				if !seenStatuses[statusLower] {
					seenStatuses[statusLower] = true
					info.Status = append(info.Status, statusParts[0])
				}
			}

		case strings.Contains(keyLower, "registrar url"):
			if info.RegistrarURL == "" {
				info.RegistrarURL = value
			}

		case strings.Contains(keyLower, "registrar whois server"):
			if info.WhoisServer == "" {
				info.WhoisServer = value
			}
		}
	}

	return info
}

// SaveWhoisResults saves WHOIS results to a JSON file
func SaveWhoisResults(domain string, info *WhoisInfo) error {
	results := WhoisResults{
		Domain:     domain,
		Info:       *info,
		LookedUpAt: time.Now(),
	}

	_, err := SaveResults(domain, "whois", results, FormatJSON)
	return err
}

// LoadWhoisResults loads the latest WHOIS results for a domain
func LoadWhoisResults(domain string) (*WhoisResults, error) {
	var results WhoisResults
	if err := LoadLatestResult(domain, "whois", &results); err != nil {
		return nil, err
	}
	return &results, nil
}

// FormatWhoisInfo returns a human-readable string representation of WHOIS info
func FormatWhoisInfo(info *WhoisInfo) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("Domain: %s\n", info.Domain))

	if info.Registrar != "" {
		b.WriteString(fmt.Sprintf("Registrar: %s\n", info.Registrar))
	}

	if info.CreatedDate != "" {
		b.WriteString(fmt.Sprintf("Created: %s\n", info.CreatedDate))
	}

	if info.UpdatedDate != "" {
		b.WriteString(fmt.Sprintf("Updated: %s\n", info.UpdatedDate))
	}

	if info.ExpiryDate != "" {
		b.WriteString(fmt.Sprintf("Expires: %s\n", info.ExpiryDate))
	}

	if len(info.NameServers) > 0 {
		b.WriteString("\nName Servers:\n")
		for _, ns := range info.NameServers {
			b.WriteString(fmt.Sprintf("  - %s\n", ns))
		}
	}

	if len(info.Status) > 0 {
		b.WriteString("\nStatus:\n")
		for _, status := range info.Status {
			b.WriteString(fmt.Sprintf("  - %s\n", status))
		}
	}

	if info.RegistrarURL != "" {
		b.WriteString(fmt.Sprintf("\nRegistrar URL: %s\n", info.RegistrarURL))
	}

	return b.String()
}

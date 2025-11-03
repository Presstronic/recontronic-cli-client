package recon

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// DNSRecord represents a single DNS record
type DNSRecord struct {
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   uint32 `json:"ttl,omitempty"`
}

// DNSInfo represents all DNS information for a subdomain
type DNSInfo struct {
	Subdomain      string    `json:"subdomain"`
	A              []string  `json:"a_records,omitempty"`
	AAAA           []string  `json:"aaaa_records,omitempty"`
	CNAME          []string  `json:"cname_records,omitempty"`
	MX             []string  `json:"mx_records,omitempty"`
	TXT            []string  `json:"txt_records,omitempty"`
	NS             []string  `json:"ns_records,omitempty"`
	CloudProvider  string    `json:"cloud_provider,omitempty"`
	TakeoverRisk   bool      `json:"takeover_risk"`
	TakeoverReason string    `json:"takeover_reason,omitempty"`
	QueryTime      time.Time `json:"query_time"`
	Error          string    `json:"error,omitempty"`
}

// DNSResults represents the complete DNS enumeration results
type DNSResults struct {
	Domain       string     `json:"domain"`
	Records      []DNSInfo  `json:"records"`
	TotalQueried int        `json:"total_queried"`
	Summary      DNSSummary `json:"summary"`
	EnumeratedAt time.Time  `json:"enumerated_at"`
}

// DNSSummary provides statistics about DNS enumeration
type DNSSummary struct {
	TotalA         int      `json:"total_a"`
	TotalAAAA      int      `json:"total_aaaa"`
	TotalMX        int      `json:"total_mx"`
	TotalTXT       int      `json:"total_txt"`
	TotalCNAME     int      `json:"total_cname"`
	TotalNS        int      `json:"total_ns"`
	TakeoverRisks  int      `json:"takeover_risks"`
	CloudProviders []string `json:"cloud_providers"`
	UniqueIPs      int      `json:"unique_ips"`
}

// DNSEnumerationOptions configures DNS enumeration
type DNSEnumerationOptions struct {
	AliveOnly     bool
	RecordTypes   []string // A, AAAA, MX, TXT, NS, CNAME
	Concurrency   int
	Timeout       time.Duration
	CheckTakeover bool
}

// Common subdomain takeover signatures
var takeoverSignatures = map[string][]string{
	"herokuapp.com":     {"No such app", "There's nothing here"},
	"github.io":         {"404", "There isn't a GitHub Pages site here"},
	"azurewebsites.net": {"404", "Error 404"},
	"cloudfront.net":    {"ERROR: The request could not be satisfied"},
	"s3.amazonaws.com":  {"NoSuchBucket", "The specified bucket does not exist"},
	"bitbucket.io":      {"Repository not found"},
	"ghost.io":          {"The thing you were looking for is no longer here"},
	"pantheonsite.io":   {"404 error unknown site"},
	"zendesk.com":       {"Help Center Closed"},
	"uservoice.com":     {"This UserVoice subdomain is currently available"},
	"surge.sh":          {"project not found"},
	"tumblr.com":        {"Whatever you were looking for doesn't currently exist"},
	"wordpress.com":     {"Do you want to register"},
	"statuspage.io":     {"You are being redirected"},
	"hubspot.net":       {"404"},
}

// Cloud provider IP ranges and patterns
var cloudProviders = map[string][]string{
	"AWS":          {"amazonaws.com", "cloudfront.net", "awsglobalaccelerator.com"},
	"Azure":        {"azurewebsites.net", "cloudapp.azure.com", "azure.com"},
	"GCP":          {"googleapis.com", "googleusercontent.com", "cloud.google.com"},
	"Cloudflare":   {"cloudflare.com", "cloudflare.net", "cloudflaressl.com"},
	"Akamai":       {"akamai.net", "akamaitechnologies.com", "akamaiedge.net"},
	"Fastly":       {"fastly.net", "fastlylb.net"},
	"DigitalOcean": {"digitaloceanspaces.com"},
	"Heroku":       {"herokuapp.com", "herokussl.com"},
}

// EnumerateDNS performs DNS enumeration for all subdomains
func EnumerateDNS(ctx context.Context, domain string, options DNSEnumerationOptions) (*DNSResults, error) {
	// Load latest subdomain results
	var subdomainResults SubdomainResults
	if err := LoadLatestResult(domain, "subdomains", &subdomainResults); err != nil {
		return nil, fmt.Errorf("failed to load subdomain results: %w", err)
	}

	// Filter subdomains if needed
	var subdomainsToQuery []Subdomain
	if options.AliveOnly {
		for _, sub := range subdomainResults.Subdomains {
			if sub.Verified != nil && sub.Verified.Status == "alive" {
				subdomainsToQuery = append(subdomainsToQuery, sub)
			}
		}
	} else {
		subdomainsToQuery = subdomainResults.Subdomains
	}

	if len(subdomainsToQuery) == 0 {
		return nil, fmt.Errorf("no subdomains to query")
	}

	// Set defaults
	if options.Concurrency == 0 {
		options.Concurrency = 10
	}
	if options.Timeout == 0 {
		options.Timeout = 5 * time.Second
	}
	if len(options.RecordTypes) == 0 {
		options.RecordTypes = []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS"}
	}

	// Create results structure
	results := &DNSResults{
		Domain:       domain,
		Records:      make([]DNSInfo, 0, len(subdomainsToQuery)),
		TotalQueried: len(subdomainsToQuery),
		EnumeratedAt: time.Now(),
	}

	// Concurrent DNS enumeration
	semaphore := make(chan struct{}, options.Concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, subdomain := range subdomainsToQuery {
		wg.Add(1)
		go func(sub Subdomain) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			info := queryDNSInfo(ctx, sub.Name, options)

			mu.Lock()
			results.Records = append(results.Records, info)
			mu.Unlock()
		}(subdomain)
	}

	wg.Wait()

	// Calculate summary
	results.Summary = calculateDNSSummary(results.Records)

	return results, nil
}

// queryDNSInfo queries all DNS records for a single subdomain
func queryDNSInfo(ctx context.Context, subdomain string, options DNSEnumerationOptions) DNSInfo {
	info := DNSInfo{
		Subdomain: subdomain,
		QueryTime: time.Now(),
	}

	resolver := &net.Resolver{
		PreferGo: true,
	}

	// Query A records
	if contains(options.RecordTypes, "A") {
		ips, err := resolver.LookupIP(ctx, "ip4", subdomain)
		if err == nil {
			for _, ip := range ips {
				info.A = append(info.A, ip.String())
			}
		}
	}

	// Query AAAA records
	if contains(options.RecordTypes, "AAAA") {
		ips, err := resolver.LookupIP(ctx, "ip6", subdomain)
		if err == nil {
			for _, ip := range ips {
				info.AAAA = append(info.AAAA, ip.String())
			}
		}
	}

	// Query CNAME records
	if contains(options.RecordTypes, "CNAME") {
		cname, err := resolver.LookupCNAME(ctx, subdomain)
		if err == nil && cname != subdomain && cname != subdomain+"." {
			info.CNAME = []string{strings.TrimSuffix(cname, ".")}

			// Check for subdomain takeover
			if options.CheckTakeover {
				checkSubdomainTakeover(&info, cname)
			}
		}
	}

	// Query MX records
	if contains(options.RecordTypes, "MX") {
		mxRecords, err := resolver.LookupMX(ctx, subdomain)
		if err == nil {
			for _, mx := range mxRecords {
				info.MX = append(info.MX, strings.TrimSuffix(mx.Host, "."))
			}
		}
	}

	// Query TXT records
	if contains(options.RecordTypes, "TXT") {
		txtRecords, err := resolver.LookupTXT(ctx, subdomain)
		if err == nil {
			info.TXT = txtRecords
		}
	}

	// Query NS records
	if contains(options.RecordTypes, "NS") {
		nsRecords, err := resolver.LookupNS(ctx, subdomain)
		if err == nil {
			for _, ns := range nsRecords {
				info.NS = append(info.NS, strings.TrimSuffix(ns.Host, "."))
			}
		}
	}

	// Identify cloud provider
	info.CloudProvider = identifyCloudProvider(info)

	return info
}

// checkSubdomainTakeover checks if a CNAME points to a potentially vulnerable service
func checkSubdomainTakeover(info *DNSInfo, cname string) {
	cname = strings.ToLower(cname)

	for service, _ := range takeoverSignatures {
		if strings.Contains(cname, service) {
			info.TakeoverRisk = true
			info.TakeoverReason = fmt.Sprintf("CNAME points to %s (potential takeover)", service)
			return
		}
	}
}

// identifyCloudProvider identifies the cloud provider based on DNS records
func identifyCloudProvider(info DNSInfo) string {
	// Check CNAME records
	for _, cname := range info.CNAME {
		cnameLower := strings.ToLower(cname)
		for provider, patterns := range cloudProviders {
			for _, pattern := range patterns {
				if strings.Contains(cnameLower, pattern) {
					return provider
				}
			}
		}
	}

	// Check NS records
	for _, ns := range info.NS {
		nsLower := strings.ToLower(ns)
		for provider, patterns := range cloudProviders {
			for _, pattern := range patterns {
				if strings.Contains(nsLower, pattern) {
					return provider
				}
			}
		}
	}

	return ""
}

// calculateDNSSummary calculates statistics from DNS records
func calculateDNSSummary(records []DNSInfo) DNSSummary {
	summary := DNSSummary{}
	uniqueIPs := make(map[string]bool)
	cloudProvidersMap := make(map[string]bool)

	for _, record := range records {
		summary.TotalA += len(record.A)
		summary.TotalAAAA += len(record.AAAA)
		summary.TotalMX += len(record.MX)
		summary.TotalTXT += len(record.TXT)
		summary.TotalCNAME += len(record.CNAME)
		summary.TotalNS += len(record.NS)

		if record.TakeoverRisk {
			summary.TakeoverRisks++
		}

		if record.CloudProvider != "" && !cloudProvidersMap[record.CloudProvider] {
			cloudProvidersMap[record.CloudProvider] = true
			summary.CloudProviders = append(summary.CloudProviders, record.CloudProvider)
		}

		for _, ip := range record.A {
			uniqueIPs[ip] = true
		}
		for _, ip := range record.AAAA {
			uniqueIPs[ip] = true
		}
	}

	summary.UniqueIPs = len(uniqueIPs)

	return summary
}

// SaveDNSResults saves DNS results to a JSON file
func SaveDNSResults(domain string, results *DNSResults) error {
	_, err := SaveResults(domain, "dns", results, FormatJSON)
	return err
}

// LoadDNSResults loads the latest DNS results for a domain
func LoadDNSResults(domain string) (*DNSResults, error) {
	var results DNSResults
	if err := LoadLatestResult(domain, "dns", &results); err != nil {
		return nil, err
	}
	return &results, nil
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}

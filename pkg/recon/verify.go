package recon

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// VerificationResult represents the verification status of a subdomain
type VerificationResult struct {
	Timestamp time.Time   `json:"timestamp"`
	Status    string      `json:"status"` // "alive", "dead", "error"
	DNS       *DNSResult  `json:"dns,omitempty"`
	HTTP      *HTTPResult `json:"http,omitempty"`
}

// DNSResult represents DNS resolution results
type DNSResult struct {
	Resolves bool     `json:"resolves"`
	IPs      []string `json:"ips,omitempty"`
	Error    string   `json:"error,omitempty"`
}

// HTTPResult represents HTTP probe results
type HTTPResult struct {
	Accessible     bool     `json:"accessible"`
	URL            string   `json:"url"`
	StatusCode     int      `json:"status_code,omitempty"`
	Title          string   `json:"title,omitempty"`
	RedirectChain  []string `json:"redirect_chain,omitempty"`
	FinalURL       string   `json:"final_url,omitempty"`
	ContentLength  int64    `json:"content_length,omitempty"`
	ResponseTimeMs int64    `json:"response_time_ms,omitempty"`
}

// VerifyOptions configures verification behavior
type VerifyOptions struct {
	Concurrency int           // Parallel probes (default: 10)
	Timeout     time.Duration // Per-probe timeout (default: 10s)
	UserAgent   string        // Custom user agent
}

// DefaultVerifyOptions returns default verification options
func DefaultVerifyOptions() VerifyOptions {
	return VerifyOptions{
		Concurrency: 10,
		Timeout:     10 * time.Second,
		UserAgent:   "Mozilla/5.0 (compatible; Recontronic/1.0)",
	}
}

// VerifySubdomain verifies a single subdomain
func VerifySubdomain(subdomain string, options VerifyOptions) (*VerificationResult, error) {
	result := &VerificationResult{
		Timestamp: time.Now(),
		Status:    "dead",
	}

	// Step 1: DNS Resolution
	dnsResult := resolveDNS(subdomain)
	result.DNS = dnsResult

	if !dnsResult.Resolves {
		result.Status = "dead"
		return result, nil
	}

	// Step 2: HTTP Probe
	httpResult := probeHTTP(subdomain, dnsResult.IPs, options)
	result.HTTP = httpResult

	if httpResult != nil && httpResult.Accessible {
		result.Status = "alive"
	}

	return result, nil
}

// VerifySubdomains verifies multiple subdomains concurrently
func VerifySubdomains(subdomains []Subdomain, options VerifyOptions) ([]Subdomain, error) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, options.Concurrency)
	resultsChan := make(chan struct {
		index  int
		result *VerificationResult
	}, len(subdomains))

	// Verify each subdomain concurrently
	for i, sub := range subdomains {
		wg.Add(1)
		go func(index int, subdomain Subdomain) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Verify subdomain
			result, err := VerifySubdomain(subdomain.Name, options)
			if err != nil {
				// Log error but don't fail
				fmt.Printf("Warning: failed to verify %s: %v\n", subdomain.Name, err)
				return
			}

			// Send result
			resultsChan <- struct {
				index  int
				result *VerificationResult
			}{index: index, result: result}
		}(i, sub)
	}

	// Close results channel when all done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Update subdomains with results
	verified := make([]Subdomain, len(subdomains))
	copy(verified, subdomains)

	for res := range resultsChan {
		verified[res.index].Verified = res.result
	}

	return verified, nil
}

// resolveDNS checks if a subdomain resolves
func resolveDNS(subdomain string) *DNSResult {
	result := &DNSResult{
		Resolves: false,
	}

	// Resolve with timeout
	resolver := &net.Resolver{
		PreferGo: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ips, err := resolver.LookupIP(ctx, "ip", subdomain)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	if len(ips) == 0 {
		result.Error = "no IP addresses found"
		return result
	}

	result.Resolves = true
	for _, ip := range ips {
		result.IPs = append(result.IPs, ip.String())
	}

	return result
}

// probeHTTP attempts to connect via HTTP/HTTPS
func probeHTTP(subdomain string, ips []string, options VerifyOptions) *HTTPResult {
	result := &HTTPResult{
		Accessible: false,
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: options.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip cert validation for recon
			},
			DisableKeepAlives: true,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Try HTTPS first, then HTTP
	protocols := []string{"https", "http"}

	for _, protocol := range protocols {
		url := fmt.Sprintf("%s://%s", protocol, subdomain)

		startTime := time.Now()
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", options.UserAgent)

		resp, err := client.Do(req)
		responseTime := time.Since(startTime)

		if err != nil {
			continue
		}
		defer resp.Body.Close()

		// Success!
		result.Accessible = true
		result.URL = url
		result.StatusCode = resp.StatusCode
		result.ResponseTimeMs = responseTime.Milliseconds()
		result.ContentLength = resp.ContentLength

		// Extract title from HTML
		if strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
			body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // Read max 1MB
			if err == nil {
				result.Title = extractTitle(string(body))
			}
		}

		// Track redirects
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			if location := resp.Header.Get("Location"); location != "" {
				result.FinalURL = location
			}
		}

		return result
	}

	return result
}

// extractTitle extracts the <title> tag from HTML
func extractTitle(html string) string {
	// Simple regex to extract title
	re := regexp.MustCompile(`(?i)<title[^>]*>([^<]+)</title>`)
	matches := re.FindStringSubmatch(html)

	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		// Limit title length
		if len(title) > 100 {
			title = title[:100] + "..."
		}
		return title
	}

	return ""
}

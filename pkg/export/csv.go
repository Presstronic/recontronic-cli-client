package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/presstronic/recontronic-cli-client/pkg/recon"
)

// ExportToCSV exports subdomain results to CSV format
func ExportToCSV(result *recon.SubdomainResults, options ExportOptions) (string, error) {
	filePath := options.OutputPath
	if filePath == "" {
		filePath = fmt.Sprintf("%s_subdomains.csv", result.Domain)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Filter subdomains based on options
	subdomains := filterSubdomains(result.Subdomains, options)

	// Determine if we have verification data
	hasVerification := false
	for _, sub := range subdomains {
		if sub.Verified != nil {
			hasVerification = true
			break
		}
	}

	// Write header
	var header []string
	if hasVerification {
		header = []string{
			"Subdomain",
			"Status",
			"DNS Resolves",
			"IP Addresses",
			"HTTP Accessible",
			"HTTP URL",
			"Status Code",
			"Title",
			"Response Time (ms)",
			"Content Length",
			"Discovered By",
			"First Seen",
		}
	} else {
		header = []string{
			"Subdomain",
			"Discovered By",
			"First Seen",
		}
	}

	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write rows
	for _, sub := range subdomains {
		var row []string

		if hasVerification && sub.Verified != nil {
			ips := "-"
			if sub.Verified.DNS != nil && len(sub.Verified.DNS.IPs) > 0 {
				ips = strings.Join(sub.Verified.DNS.IPs, ";")
			}

			dnsResolves := "false"
			if sub.Verified.DNS != nil && sub.Verified.DNS.Resolves {
				dnsResolves = "true"
			}

			httpAccessible := "false"
			httpURL := "-"
			statusCode := "-"
			title := "-"
			responseTime := "-"
			contentLength := "-"

			if sub.Verified.HTTP != nil {
				if sub.Verified.HTTP.Accessible {
					httpAccessible = "true"
				}
				httpURL = sub.Verified.HTTP.URL
				if sub.Verified.HTTP.StatusCode > 0 {
					statusCode = strconv.Itoa(sub.Verified.HTTP.StatusCode)
				}
				if sub.Verified.HTTP.Title != "" {
					title = sub.Verified.HTTP.Title
				}
				if sub.Verified.HTTP.ResponseTimeMs > 0 {
					responseTime = strconv.FormatInt(sub.Verified.HTTP.ResponseTimeMs, 10)
				}
				if sub.Verified.HTTP.ContentLength > 0 {
					contentLength = strconv.FormatInt(sub.Verified.HTTP.ContentLength, 10)
				}
			}

			row = []string{
				sub.Name,
				sub.Verified.Status,
				dnsResolves,
				ips,
				httpAccessible,
				httpURL,
				statusCode,
				title,
				responseTime,
				contentLength,
				strings.Join(sub.DiscoveredBy, ";"),
				sub.FirstSeen.Format("2006-01-02 15:04:05"),
			}
		} else {
			row = []string{
				sub.Name,
				strings.Join(sub.DiscoveredBy, ";"),
				sub.FirstSeen.Format("2006-01-02 15:04:05"),
			}
		}

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return filePath, nil
}

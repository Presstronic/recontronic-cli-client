package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/recon"
	"github.com/presstronic/recontronic-cli-client/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	dnsAliveOnly     bool
	dnsRecordTypes   string
	dnsConcurrency   int
	dnsTimeout       time.Duration
	dnsCheckTakeover bool
)

var reconDNSCmd = &cobra.Command{
	Use:   "dns <domain>",
	Short: "Enumerate DNS records for subdomains",
	Long: `Enumerate DNS records for all discovered subdomains including:
  - A records (IPv4 addresses)
  - AAAA records (IPv6 addresses)
  - CNAME records (aliases)
  - MX records (mail servers)
  - TXT records (SPF, DMARC, verification records)
  - NS records (name servers)

This command also:
  - Identifies cloud providers (AWS, Azure, GCP, Cloudflare, Akamai)
  - Detects potential subdomain takeover opportunities
  - Maps subdomains to IP addresses for port scanning

Results are automatically saved to ~/.recon-cli/results/<domain>/dns_<timestamp>.json

Examples:
  recon dns example.com
  recon dns example.com --alive-only
  recon dns example.com --types A,AAAA,MX
  recon dns example.com --check-takeover
  recon dns example.com --concurrency 20 --timeout 10s`,
	Args: cobra.ExactArgs(1),
	RunE: runReconDNS,
}

func init() {
	reconDNSCmd.Flags().BoolVar(&dnsAliveOnly, "alive-only", true, "Only query DNS for alive subdomains")
	reconDNSCmd.Flags().StringVar(&dnsRecordTypes, "types", "A,AAAA,CNAME,MX,TXT,NS", "DNS record types to query (comma-separated)")
	reconDNSCmd.Flags().IntVar(&dnsConcurrency, "concurrency", 10, "Number of concurrent DNS queries")
	reconDNSCmd.Flags().DurationVar(&dnsTimeout, "timeout", 5*time.Second, "Timeout per DNS query")
	reconDNSCmd.Flags().BoolVar(&dnsCheckTakeover, "check-takeover", true, "Check for subdomain takeover opportunities")
	reconCmd.AddCommand(reconDNSCmd)
}

func runReconDNS(cmd *cobra.Command, args []string) error {
	domain := args[0]

	// Validate domain
	if err := recon.ValidateDomain(domain); err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}

	fmt.Printf("Enumerating DNS records for %s\n", domain)
	fmt.Println("Mode: Passive DNS enumeration")

	// Parse record types
	recordTypes := strings.Split(dnsRecordTypes, ",")
	for i, rt := range recordTypes {
		recordTypes[i] = strings.TrimSpace(strings.ToUpper(rt))
	}

	// Setup options
	options := recon.DNSEnumerationOptions{
		AliveOnly:     dnsAliveOnly,
		RecordTypes:   recordTypes,
		Concurrency:   dnsConcurrency,
		Timeout:       dnsTimeout,
		CheckTakeover: dnsCheckTakeover,
	}

	ctx := context.Background()
	startTime := time.Now()

	// Start progress indicator
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(startTime)
				fmt.Printf("\rProgress: Querying DNS records... [%s elapsed]", elapsed.Round(time.Second))
			case <-done:
				return
			}
		}
	}()

	// Perform DNS enumeration
	results, err := recon.EnumerateDNS(ctx, domain, options)
	done <- true
	fmt.Printf("\r\033[K") // Clear progress line

	if err != nil {
		return fmt.Errorf("DNS enumeration failed: %w", err)
	}

	duration := time.Since(startTime)

	// Save results
	if err := recon.SaveDNSResults(domain, results); err != nil {
		fmt.Printf("Warning: Failed to save results: %v\n", err)
	} else {
		fmt.Printf("\n✓ Results saved to ~/.recon-cli/results/%s/\n", domain)
	}

	// Display summary
	displayDNSSummary(results, duration)

	// Display key findings
	displayKeyFindings(results)

	// Log activity
	activityResult := fmt.Sprintf("%d IPs, %d CNAMEs", results.Summary.UniqueIPs, results.Summary.TotalCNAME)
	if results.Summary.TakeoverRisks > 0 {
		activityResult += fmt.Sprintf(", %d takeover risks", results.Summary.TakeoverRisks)
	}

	ui.LogActivity(ui.ActivityEntry{
		Timestamp: time.Now(),
		Domain:    domain,
		Action:    "dns",
		Status:    "completed",
		Result:    activityResult,
	})

	return nil
}

func displayDNSSummary(results *recon.DNSResults, duration time.Duration) {
	fmt.Println("\nSummary:")
	fmt.Printf("  Subdomains queried: %d\n", results.TotalQueried)
	fmt.Printf("  A records: %d\n", results.Summary.TotalA)
	fmt.Printf("  AAAA records: %d\n", results.Summary.TotalAAAA)
	fmt.Printf("  CNAME records: %d\n", results.Summary.TotalCNAME)
	fmt.Printf("  MX records: %d\n", results.Summary.TotalMX)
	fmt.Printf("  TXT records: %d\n", results.Summary.TotalTXT)
	fmt.Printf("  NS records: %d\n", results.Summary.TotalNS)
	fmt.Printf("  Unique IPs: %d\n", results.Summary.UniqueIPs)
	fmt.Printf("  Duration: %s\n", duration.Round(time.Second))
}

func displayKeyFindings(results *recon.DNSResults) {
	fmt.Println("\nKey Findings:")

	// Subdomain takeover risks
	if results.Summary.TakeoverRisks > 0 {
		fmt.Printf("  ⚠️  Potential subdomain takeovers: %d\n", results.Summary.TakeoverRisks)

		// Show first few takeover risks
		count := 0
		for _, record := range results.Records {
			if record.TakeoverRisk && count < 5 {
				fmt.Printf("      - %s → %s\n", record.Subdomain, record.TakeoverReason)
				count++
			}
		}
		if results.Summary.TakeoverRisks > 5 {
			fmt.Printf("      ... and %d more (see JSON results)\n", results.Summary.TakeoverRisks-5)
		}
	} else {
		fmt.Println("  ✓ No obvious subdomain takeover risks detected")
	}

	// Cloud providers
	if len(results.Summary.CloudProviders) > 0 {
		fmt.Printf("  ☁️  Cloud providers detected: %s\n", strings.Join(results.Summary.CloudProviders, ", "))
	}

	// Mail servers
	if results.Summary.TotalMX > 0 {
		fmt.Printf("  📧 Mail servers found: %d MX records\n", results.Summary.TotalMX)

		// Show unique mail server domains
		mailServers := make(map[string]bool)
		for _, record := range results.Records {
			for _, mx := range record.MX {
				// Extract domain from mail server
				parts := strings.Split(mx, ".")
				if len(parts) >= 2 {
					mailDomain := parts[len(parts)-2] + "." + parts[len(parts)-1]
					mailServers[mailDomain] = true
				}
			}
		}
		if len(mailServers) > 0 {
			var domains []string
			for domain := range mailServers {
				domains = append(domains, domain)
			}
			fmt.Printf("      Providers: %s\n", strings.Join(domains, ", "))
		}
	}

	// Security records
	hasSecurityRecords := false
	hasSPF := false
	hasDMARC := false
	hasDKIM := false

	for _, record := range results.Records {
		for _, txt := range record.TXT {
			txtLower := strings.ToLower(txt)
			if strings.HasPrefix(txtLower, "v=spf1") {
				hasSPF = true
				hasSecurityRecords = true
			}
			if strings.HasPrefix(txtLower, "v=dmarc1") {
				hasDMARC = true
				hasSecurityRecords = true
			}
			if strings.Contains(txtLower, "dkim") {
				hasDKIM = true
				hasSecurityRecords = true
			}
		}
	}

	if hasSecurityRecords {
		fmt.Printf("  🔒 Security records: SPF (%v), DMARC (%v), DKIM (%v)\n",
			formatBool(hasSPF), formatBool(hasDMARC), formatBool(hasDKIM))
	}

	// Sample records
	fmt.Println("\nSample DNS Records:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "  SUBDOMAIN\tRECORD TYPE\tVALUE\tCLOUD")

	count := 0
	for _, record := range results.Records {
		if count >= 10 {
			break
		}

		// Show A records
		if len(record.A) > 0 {
			cloud := ""
			if record.CloudProvider != "" {
				cloud = record.CloudProvider
			}
			fmt.Fprintf(w, "  %s\tA\t%s\t%s\n", record.Subdomain, record.A[0], cloud)
			count++
		}

		// Show CNAME records
		if len(record.CNAME) > 0 {
			cloud := ""
			if record.CloudProvider != "" {
				cloud = record.CloudProvider
			}
			risk := ""
			if record.TakeoverRisk {
				risk = " ⚠️"
			}
			fmt.Fprintf(w, "  %s\tCNAME\t%s%s\t%s\n", record.Subdomain, record.CNAME[0], risk, cloud)
			count++
		}
	}

	w.Flush()

	if len(results.Records) > 10 {
		fmt.Printf("\n  ... and %d more records (see JSON results for complete data)\n", len(results.Records)-10)
	}
}

func formatBool(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

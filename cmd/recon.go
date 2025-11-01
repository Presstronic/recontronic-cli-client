package cmd

import (
	"fmt"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/recon"
	"github.com/presstronic/recontronic-cli-client/pkg/ui"
	"github.com/spf13/cobra"
)

var reconCmd = &cobra.Command{
	Use:   "recon",
	Short: "Reconnaissance tools",
	Long: `Passive reconnaissance tools for domain enumeration and information gathering.

Available subcommands:
  subdomain - Find subdomains using multiple sources
  verify    - Verify which subdomains are alive
  dns       - Enumerate DNS records
  whois     - Lookup WHOIS information`,
}

var reconSubdomainCmd = &cobra.Command{
	Use:   "subdomain <domain>",
	Short: "Find subdomains using multiple sources",
	Long: `Find subdomains for a target domain using multiple enumeration sources.

Available sources:
  - subfinder (if installed)
  - amass (if installed - future)
  - assetfinder (if installed - future)
  - crt.sh (built-in - future)

The tool will automatically detect which sources are available and use them all.`,
	Args: cobra.ExactArgs(1),
	RunE: runReconSubdomain,
}

var (
	subdomainSources []string
)

func init() {
	rootCmd.AddCommand(reconCmd)
	reconCmd.AddCommand(reconSubdomainCmd)

	// Flags for subdomain command
	reconSubdomainCmd.Flags().StringSliceVar(&subdomainSources, "sources", []string{}, "Specific sources to use (comma-separated)")
}

func runReconSubdomain(cmd *cobra.Command, args []string) error {
	domain := args[0]

	// Validate domain
	if err := recon.ValidateDomain(domain); err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}

	fmt.Printf("Finding subdomains for %s\n", domain)
	fmt.Println("Mode: Passive reconnaissance (safe, no active scanning)\n")

	// Detect available sources (in order of speed/reliability)
	var sources []recon.SubdomainSource

	// crt.sh - always available (API-based)
	crtshSource := &recon.CrtShSource{}
	if crtshSource.IsAvailable() {
		sources = append(sources, crtshSource)
	}

	// subfinder - fast and comprehensive
	subfinderSource := &recon.SubfinderSource{}
	if subfinderSource.IsAvailable() {
		sources = append(sources, subfinderSource)
	}

	// assetfinder - additional coverage
	assetfinderSource := &recon.AssetfinderSource{}
	if assetfinderSource.IsAvailable() {
		sources = append(sources, assetfinderSource)
	}

	// amass - most comprehensive but slowest
	amassSource := &recon.AmassSource{}
	if amassSource.IsAvailable() {
		sources = append(sources, amassSource)
	}

	// Check if any sources are available
	if len(sources) == 0 {
		return fmt.Errorf("no enumeration tools available. At minimum, curl must be installed for crt.sh")
	}

	// Show which sources will be used
	fmt.Println("Sources:")
	for _, source := range sources {
		fmt.Printf("  ✓ %s\n", source.Name())
	}
	fmt.Println()

	// Run enumeration
	startTime := time.Now()
	results, err := recon.EnumerateSubdomains(domain, sources)
	if err != nil {
		return fmt.Errorf("enumeration failed: %w", err)
	}
	duration := time.Since(startTime)

	// Display summary
	fmt.Println("Results:")
	for source, count := range results.Summary {
		fmt.Printf("  %s: %d subdomains\n", source, count)
	}
	fmt.Printf("\nTotal unique: %d subdomains\n", results.TotalUnique)
	fmt.Printf("Time taken: %s\n\n", duration.Round(time.Second))

	// Save results
	filePath, err := recon.SaveResults(domain, "subdomains", results, recon.FormatJSON)
	if err != nil {
		return fmt.Errorf("failed to save results: %w", err)
	}

	fmt.Printf("Saved to: %s\n\n", filePath)

	// Show first 10 subdomains
	if len(results.Subdomains) > 0 {
		fmt.Println("Sample subdomains (first 10):")
		limit := 10
		if len(results.Subdomains) < limit {
			limit = len(results.Subdomains)
		}
		for i := 0; i < limit; i++ {
			sub := results.Subdomains[i]
			sources := fmt.Sprintf("[%s]", sub.DiscoveredBy[0])
			if len(sub.DiscoveredBy) > 1 {
				sources = fmt.Sprintf("[%d sources]", len(sub.DiscoveredBy))
			}
			fmt.Printf("  %s %s\n", sub.Name, sources)
		}

		if len(results.Subdomains) > 10 {
			fmt.Printf("  ... and %d more\n", len(results.Subdomains)-10)
		}
	}

	// Log activity
	activityResult := fmt.Sprintf("%d found", results.TotalUnique)
	if err := ui.LogActivity(ui.ActivityEntry{
		Timestamp: time.Now(),
		Domain:    domain,
		Action:    "subdomain enum",
		Status:    "completed",
		Result:    activityResult,
	}); err != nil {
		// Don't fail if logging fails
		fmt.Printf("Warning: failed to log activity: %v\n", err)
	}

	fmt.Println("\nNext: Run 'recon verify", domain, "' to check which subdomains are alive (coming soon)")

	return nil
}

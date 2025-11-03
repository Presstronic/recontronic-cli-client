package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/recon"
	"github.com/presstronic/recontronic-cli-client/pkg/ui"
	"github.com/spf13/cobra"
)

var (
	whoisTimeout time.Duration
	whoisRaw     bool
	whoisJSON    bool
)

var reconWhoisCmd = &cobra.Command{
	Use:   "whois <domain>",
	Short: "Lookup WHOIS information for a domain",
	Long: `Lookup WHOIS information for a target domain including:
  - Registrar information
  - Registration and expiration dates
  - Name servers
  - Domain status
  - Contact information (if available)

Results are automatically saved to ~/.recon-cli/results/<domain>/whois_<timestamp>.json

Examples:
  recon whois example.com
  recon whois example.com --timeout 30s
  recon whois example.com --json
  recon whois example.com --raw`,
	Args: cobra.ExactArgs(1),
	RunE: runReconWhois,
}

func init() {
	reconWhoisCmd.Flags().DurationVar(&whoisTimeout, "timeout", 30*time.Second, "Timeout for WHOIS lookup")
	reconWhoisCmd.Flags().BoolVar(&whoisRaw, "raw", false, "Show raw WHOIS output")
	reconWhoisCmd.Flags().BoolVar(&whoisJSON, "json", false, "Output results as JSON")
	reconCmd.AddCommand(reconWhoisCmd)
}

func runReconWhois(cmd *cobra.Command, args []string) error {
	domain := args[0]

	// Validate domain
	if err := recon.ValidateDomain(domain); err != nil {
		return fmt.Errorf("invalid domain: %w", err)
	}

	fmt.Printf("Looking up WHOIS information for %s\n", domain)
	fmt.Println("Mode: Passive reconnaissance (WHOIS query)")

	ctx := context.Background()

	// Perform WHOIS lookup
	info, err := recon.LookupWhois(ctx, domain, whoisTimeout)
	if err != nil {
		return fmt.Errorf("WHOIS lookup failed: %w", err)
	}

	// Save results
	if err := recon.SaveWhoisResults(domain, info); err != nil {
		fmt.Printf("Warning: Failed to save results: %v\n", err)
	} else {
		fmt.Printf("\n✓ Results saved to ~/.recon-cli/results/%s/\n", domain)
	}

	// Log activity
	result := fmt.Sprintf("Registrar: %s", info.Registrar)
	if info.ExpiryDate != "" {
		result = fmt.Sprintf("%s, Expires: %s", result, info.ExpiryDate)
	}
	ui.LogActivity(ui.ActivityEntry{
		Timestamp: time.Now(),
		Domain:    domain,
		Action:    "whois",
		Status:    "completed",
		Result:    result,
	})

	// Display results based on flags
	if whoisJSON {
		// Output as JSON
		jsonData, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(jsonData))
	} else if whoisRaw {
		// Output raw WHOIS data
		fmt.Println("\n" + info.RawOutput)
	} else {
		// Output formatted summary
		fmt.Println("\n" + recon.FormatWhoisInfo(info))
	}

	return nil
}

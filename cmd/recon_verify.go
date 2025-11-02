package cmd

import (
	"fmt"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/recon"
	"github.com/presstronic/recontronic-cli-client/pkg/ui"
	"github.com/spf13/cobra"
)

var reconVerifyCmd = &cobra.Command{
	Use:   "verify <domain>",
	Short: "Verify which subdomains are alive",
	Long: `Verify which discovered subdomains are actually alive and responding.

This command:
1. Loads the latest subdomain results for the domain
2. Performs DNS resolution checks
3. Probes HTTP/HTTPS endpoints
4. Updates the results file with verification data

The verification process is passive and only checks if subdomains respond.`,
	Args: cobra.ExactArgs(1),
	RunE: runReconVerify,
}

var (
	verifyConcurrency int
	verifyTimeout     int
)

func init() {
	reconCmd.AddCommand(reconVerifyCmd)

	// Flags for verify command
	reconVerifyCmd.Flags().IntVar(&verifyConcurrency, "concurrency", 10, "Number of parallel probes")
	reconVerifyCmd.Flags().IntVar(&verifyTimeout, "timeout", 10, "Timeout per probe in seconds")
}

func runReconVerify(cmd *cobra.Command, args []string) error {
	domain := args[0]

	fmt.Printf("Verifying subdomains for %s\n", domain)
	fmt.Println("Mode: Passive verification (DNS + HTTP probing)")

	// Load latest subdomain results
	var results recon.SubdomainResults
	if err := recon.LoadLatestResult(domain, "subdomains", &results); err != nil {
		return fmt.Errorf("failed to load subdomain results: %w\nRun 'recon subdomain %s' first", err, domain)
	}

	fmt.Printf("Loaded %d subdomains from previous scan\n", len(results.Subdomains))
	fmt.Printf("Starting verification (concurrency: %d, timeout: %ds)\n\n", verifyConcurrency, verifyTimeout)

	// Set up verification options
	options := recon.DefaultVerifyOptions()
	options.Concurrency = verifyConcurrency
	options.Timeout = time.Duration(verifyTimeout) * time.Second

	// Track progress
	startTime := time.Now()
	total := len(results.Subdomains)
	verified := 0
	alive := 0

	// Progress ticker
	progressTicker := time.NewTicker(2 * time.Second)
	defer progressTicker.Stop()

	// Channel to track completion
	done := make(chan bool)

	// Show progress in background
	go func() {
		for {
			select {
			case <-progressTicker.C:
				if verified > 0 {
					pct := float64(verified) / float64(total) * 100
					fmt.Printf("\rProgress: %d/%d (%.1f%%) | Alive: %d", verified, total, pct, alive)
				}
			case <-done:
				return
			}
		}
	}()

	// Verify subdomains with progress tracking
	verifiedSubdomains := make([]recon.Subdomain, 0, len(results.Subdomains))
	batchSize := options.Concurrency

	for i := 0; i < len(results.Subdomains); i += batchSize {
		end := i + batchSize
		if end > len(results.Subdomains) {
			end = len(results.Subdomains)
		}

		batch := results.Subdomains[i:end]
		verifiedBatch, err := recon.VerifySubdomains(batch, options)
		if err != nil {
			done <- true
			return fmt.Errorf("verification failed: %w", err)
		}

		for _, sub := range verifiedBatch {
			verifiedSubdomains = append(verifiedSubdomains, sub)
			verified++
			if sub.Verified != nil && sub.Verified.Status == "alive" {
				alive++
			}
		}
	}

	done <- true
	duration := time.Since(startTime)

	// Clear progress line
	fmt.Print("\r" + string(make([]byte, 80)) + "\r")

	// Update results with verification data
	results.Subdomains = verifiedSubdomains

	// Add verification summary to results
	dead := verified - alive
	if results.Summary == nil {
		results.Summary = make(map[string]int)
	}
	results.Summary["verified_total"] = verified
	results.Summary["verified_alive"] = alive
	results.Summary["verified_dead"] = dead

	// Save updated results
	filePath, err := recon.SaveResults(domain, "subdomains", results, recon.FormatJSON)
	if err != nil {
		return fmt.Errorf("failed to save results: %w", err)
	}

	// Display summary
	fmt.Println("\nVerification Complete!")
	fmt.Printf("Time taken: %s\n\n", duration.Round(time.Second))

	fmt.Println("Results:")
	fmt.Printf("  Total verified: %d subdomains\n", verified)
	fmt.Printf("  Alive:          %d (%.1f%%)\n", alive, float64(alive)/float64(verified)*100)
	fmt.Printf("  Dead:           %d (%.1f%%)\n", dead, float64(dead)/float64(verified)*100)
	fmt.Printf("\nUpdated: %s\n\n", filePath)

	// Show sample alive subdomains
	if alive > 0 {
		fmt.Println("Sample alive subdomains (first 10):")
		count := 0
		for _, sub := range verifiedSubdomains {
			if sub.Verified != nil && sub.Verified.Status == "alive" && count < 10 {
				statusCode := ""
				if sub.Verified.HTTP != nil {
					statusCode = fmt.Sprintf(" [%d]", sub.Verified.HTTP.StatusCode)
				}
				title := ""
				if sub.Verified.HTTP != nil && sub.Verified.HTTP.Title != "" {
					title = fmt.Sprintf(" - %s", sub.Verified.HTTP.Title)
					if len(title) > 50 {
						title = title[:50] + "..."
					}
				}
				fmt.Printf("  %s%s%s\n", sub.Verified.HTTP.URL, statusCode, title)
				count++
			}
		}
	}

	// Log activity
	activityResult := fmt.Sprintf("%d/%d alive", alive, verified)
	if err := ui.LogActivity(ui.ActivityEntry{
		Timestamp: time.Now(),
		Domain:    domain,
		Action:    "verify",
		Status:    "completed",
		Result:    activityResult,
	}); err != nil {
		fmt.Printf("Warning: failed to log activity: %v\n", err)
	}

	return nil
}

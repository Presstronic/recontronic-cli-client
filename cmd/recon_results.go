package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/presstronic/recontronic-cli-client/pkg/recon"
	"github.com/spf13/cobra"
)

var reconResultsCmd = &cobra.Command{
	Use:   "results",
	Short: "Manage reconnaissance results",
	Long: `View, filter, and manage stored reconnaissance results.

Available subcommands:
  list  - List all stored results
  view  - View specific result details`,
}

var reconResultsListCmd = &cobra.Command{
	Use:   "list [domain]",
	Short: "List all stored results",
	Long: `List all stored reconnaissance results grouped by domain.

If a domain is specified, only results for that domain will be shown.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runReconResultsList,
}

var reconResultsViewCmd = &cobra.Command{
	Use:   "view <domain>",
	Short: "View subdomain results for a domain",
	Long: `View the most recent subdomain results for a domain.

Supports filtering options to narrow down results.`,
	Args: cobra.ExactArgs(1),
	RunE: runReconResultsView,
}

var (
	viewAliveOnly  bool
	viewDeadOnly   bool
	viewStatusCode int
	viewSource     string
	viewLimit      int
)

func init() {
	reconCmd.AddCommand(reconResultsCmd)
	reconResultsCmd.AddCommand(reconResultsListCmd)
	reconResultsCmd.AddCommand(reconResultsViewCmd)

	// Flags for view command
	reconResultsViewCmd.Flags().BoolVar(&viewAliveOnly, "alive-only", false, "Show only alive subdomains")
	reconResultsViewCmd.Flags().BoolVar(&viewDeadOnly, "dead-only", false, "Show only dead subdomains")
	reconResultsViewCmd.Flags().IntVar(&viewStatusCode, "status", 0, "Filter by HTTP status code")
	reconResultsViewCmd.Flags().StringVar(&viewSource, "source", "", "Filter by discovery source")
	reconResultsViewCmd.Flags().IntVarP(&viewLimit, "limit", "n", 0, "Limit number of results shown (0 = all)")
}

func runReconResultsList(cmd *cobra.Command, args []string) error {
	// If domain specified, list only that domain
	if len(args) == 1 {
		domain := args[0]
		return listResultsForDomain(domain)
	}

	// List all results
	resultsByDomain, err := recon.ListResults()
	if err != nil {
		return fmt.Errorf("failed to list results: %w", err)
	}

	if len(resultsByDomain) == 0 {
		fmt.Println("No results found.")
		fmt.Println("\nRun 'recon subdomain <domain>' to start collecting data.")
		return nil
	}

	fmt.Println("Stored Results:")
	fmt.Println()

	// Sort domains alphabetically
	domains := make([]string, 0, len(resultsByDomain))
	for domain := range resultsByDomain {
		domains = append(domains, domain)
	}

	for _, domain := range domains {
		results := resultsByDomain[domain]

		fmt.Printf("📁 %s\n", domain)

		for _, result := range results {
			// Format timestamp
			timeStr := formatTimeAgo(result.Timestamp)

			// Format verification status
			var status string
			if result.Verified {
				status = fmt.Sprintf("✓ verified (%d alive, %d dead)", result.AliveCount, result.DeadCount)
			} else {
				status = "⚠ not verified"
			}

			// Format file size
			sizeStr := recon.FormatFileSize(result.FileSize)

			fmt.Printf("  %s  %s  (%d total)  %s  [%s]\n",
				timeStr,
				result.ToolName,
				result.TotalCount,
				status,
				sizeStr,
			)
		}
		fmt.Println()
	}

	return nil
}

func listResultsForDomain(domain string) error {
	results, err := recon.ListResultsForDomain(domain)
	if err != nil {
		return fmt.Errorf("failed to list results for %s: %w", domain, err)
	}

	if len(results) == 0 {
		fmt.Printf("No results found for %s\n", domain)
		fmt.Printf("\nRun 'recon subdomain %s' to start collecting data.\n", domain)
		return nil
	}

	fmt.Printf("Results for %s:\n\n", domain)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tTOOL\tTOTAL\tALIVE\tDEAD\tSTATUS\tSIZE")
	fmt.Fprintln(w, "─────────\t────\t─────\t─────\t────\t──────\t────")

	for _, result := range results {
		timeStr := result.Timestamp.Format("2006-01-02 15:04")

		var status string
		if result.Verified {
			status = "verified"
		} else {
			status = "unverified"
		}

		aliveStr := "-"
		deadStr := "-"
		if result.Verified {
			aliveStr = fmt.Sprintf("%d", result.AliveCount)
			deadStr = fmt.Sprintf("%d", result.DeadCount)
		}

		sizeStr := recon.FormatFileSize(result.FileSize)

		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			timeStr,
			result.ToolName,
			result.TotalCount,
			aliveStr,
			deadStr,
			status,
			sizeStr,
		)
	}

	w.Flush()
	fmt.Println()

	return nil
}

func runReconResultsView(cmd *cobra.Command, args []string) error {
	domain := args[0]

	// Build query options
	options := recon.QueryOptions{
		AliveOnly:  viewAliveOnly,
		DeadOnly:   viewDeadOnly,
		StatusCode: viewStatusCode,
		Source:     viewSource,
	}

	// Load and filter subdomains
	subdomains, err := recon.QuerySubdomains(domain, options)
	if err != nil {
		return fmt.Errorf("failed to load results: %w", err)
	}

	if len(subdomains) == 0 {
		fmt.Printf("No results found for %s", domain)
		if viewAliveOnly || viewDeadOnly || viewStatusCode != 0 || viewSource != "" {
			fmt.Print(" matching filters")
		}
		fmt.Println()
		return nil
	}

	// Get result metadata for header
	resultInfo, err := recon.ListResultsForDomain(domain)
	if err != nil {
		return err
	}

	if len(resultInfo) > 0 {
		latest := resultInfo[0]
		fmt.Printf("Results for %s\n", domain)
		fmt.Printf("Scanned: %s (%s)\n", latest.Timestamp.Format("2006-01-02 15:04:05"), formatTimeAgo(latest.Timestamp))
		if len(latest.SourcesUsed) > 0 {
			fmt.Printf("Sources: %s\n", strings.Join(latest.SourcesUsed, ", "))
		}
		fmt.Printf("Total: %d subdomains", latest.TotalCount)
		if latest.Verified {
			fmt.Printf(" (%d alive, %d dead)", latest.AliveCount, latest.DeadCount)
		}
		fmt.Println()
		fmt.Println()
	}

	// Apply limit
	if viewLimit > 0 && len(subdomains) > viewLimit {
		subdomains = subdomains[:viewLimit]
	}

	// Display results
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// Determine if we need verification columns
	hasVerification := false
	for _, sub := range subdomains {
		if sub.Verified != nil {
			hasVerification = true
			break
		}
	}

	// Print header
	if hasVerification {
		fmt.Fprintln(w, "SUBDOMAIN\tSTATUS\tHTTP\tTITLE\tSOURCES")
		fmt.Fprintln(w, "─────────\t──────\t────\t─────\t───────")
	} else {
		fmt.Fprintln(w, "SUBDOMAIN\tSOURCES")
		fmt.Fprintln(w, "─────────\t───────")
	}

	// Print subdomains
	for _, sub := range subdomains {
		sources := strings.Join(sub.DiscoveredBy, ",")

		if hasVerification && sub.Verified != nil {
			status := sub.Verified.Status

			httpInfo := "-"
			title := "-"

			if sub.Verified.HTTP != nil && sub.Verified.HTTP.Accessible {
				httpInfo = fmt.Sprintf("%d", sub.Verified.HTTP.StatusCode)
				if sub.Verified.HTTP.Title != "" {
					title = sub.Verified.HTTP.Title
					// Truncate long titles
					if len(title) > 40 {
						title = title[:37] + "..."
					}
				}
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				sub.Name,
				status,
				httpInfo,
				title,
				sources,
			)
		} else {
			fmt.Fprintf(w, "%s\t%s\n", sub.Name, sources)
		}
	}

	w.Flush()

	// Show totals
	fmt.Printf("\nShowing %d subdomain(s)", len(subdomains))
	if viewLimit > 0 {
		fmt.Printf(" (limited to %d)", viewLimit)
	}
	fmt.Println()

	// Show next steps
	if !hasVerification {
		fmt.Printf("\nNext: Run 'recon verify %s' to check which subdomains are alive\n", domain)
	}

	return nil
}

package ui

import (
	"fmt"
	"strings"

	"github.com/presstronic/recontronic-cli-client/pkg/config"
)

// DisplayDashboard shows the main dashboard
func DisplayDashboard(cfg *config.Config) error {
	// Try to display rich dashboard, fallback to simple if it fails
	if err := displaySimpleDashboard(cfg); err != nil {
		return err
	}
	return nil
}

// displaySimpleDashboard shows a simple text-based dashboard
func displaySimpleDashboard(cfg *config.Config) error {
	// Gather all data
	stats, err := GatherStats()
	if err != nil {
		stats = &DashboardStats{} // Use empty stats on error
	}

	systemStatus, err := GetSystemStatus(cfg)
	if err != nil {
		return fmt.Errorf("failed to get system status: %w", err)
	}

	activities, err := GetRecentActivity(5)
	if err != nil {
		activities = []ActivityEntry{} // Use empty activities on error
	}

	suggestions, err := GenerateSuggestions()
	if err != nil {
		suggestions = []Suggestion{} // Use empty suggestions on error
	}

	// Print dashboard
	printHeader(cfg, systemStatus)
	fmt.Println()
	printQuickStats(stats)
	fmt.Println()
	printRecentActivity(activities)
	fmt.Println()
	printSystemStatus(systemStatus)
	fmt.Println()
	if len(suggestions) > 0 {
		printSuggestions(suggestions)
		fmt.Println()
	}
	printFooter()
	fmt.Println()

	return nil
}

func printHeader(cfg *config.Config, status *SystemStatus) {
	line := strings.Repeat("═", 80)
	fmt.Println("╔" + line + "╗")

	// Title and status line
	title := " Recontronic CLI"
	serverInfo := ""
	if cfg != nil && cfg.Server != "" {
		serverInfo = fmt.Sprintf(" Server: %s", cfg.Server)
		if status.ServerStatus == "connected" {
			serverInfo += " [Connected]"
		} else {
			serverInfo += " [Offline]"
		}
	}

	authInfo := ""
	if status.AuthStatus == "authenticated" {
		authInfo = " | Authenticated"
	}

	toolsInfo := fmt.Sprintf(" | Tools: %d/%d available", status.ToolsAvailable, status.ToolsTotal)

	headerLine := title + serverInfo + authInfo + toolsInfo
	padding := 82 - len(headerLine)
	if padding < 0 {
		padding = 0
	}

	fmt.Printf("║%s%s║\n", headerLine, strings.Repeat(" ", padding))
	fmt.Println("╠" + line + "╣")
}

func printQuickStats(stats *DashboardStats) {
	fmt.Println("║ 📊 QUICK STATISTICS")
	fmt.Println("║ ┌────────────────────────────────────────────────────────────────────────────┐")

	fmt.Printf("║ │ Domains Scanned:  %-60d │\n", stats.TotalDomains)
	fmt.Printf("║ │ Subdomains Found: %-60d │\n", stats.TotalSubdomains)
	fmt.Printf("║ │ Alive Targets:    %-60d │\n", stats.TotalAlive)
	fmt.Printf("║ │ Last 24h Scans:   %-60d │\n", stats.ScansLast24h)
	fmt.Printf("║ │ Storage Used:     %-60s │\n", FormatBytes(stats.StorageUsed))

	fmt.Println("║ └────────────────────────────────────────────────────────────────────────────┘")
}

func printRecentActivity(activities []ActivityEntry) {
	fmt.Println("║ 🔍 RECENT ACTIVITY")
	fmt.Println("║ ┌────────────────────────────────────────────────────────────────────────────┐")

	if len(activities) == 0 {
		fmt.Println("║ │ No recent activity                                                         │")
	} else {
		for _, activity := range activities {
			timeAgo := FormatTimeAgo(activity.Timestamp)
			statusIcon := "✓"
			if activity.Status == "failed" {
				statusIcon = "✗"
			} else if activity.Status == "in_progress" {
				statusIcon = "⋯"
			}

			line := fmt.Sprintf(" %s  %s  %s - %s (%s)",
				statusIcon,
				timeAgo,
				activity.Domain,
				activity.Action,
				activity.Result)

			// Truncate if too long
			if len(line) > 76 {
				line = line[:73] + "..."
			}

			padding := 78 - len(line)
			if padding < 0 {
				padding = 0
			}

			fmt.Printf("║ │%s%s│\n", line, strings.Repeat(" ", padding))
		}
	}

	fmt.Println("║ └────────────────────────────────────────────────────────────────────────────┘")
}

func printSystemStatus(status *SystemStatus) {
	fmt.Println("║ ⚙️  SYSTEM STATUS")
	fmt.Println("║ ┌────────────────────────────────────────────────────────────────────────────┐")

	for _, tool := range status.Tools {
		icon := "✓"
		if !tool.Installed {
			icon = "✗"
		}

		var line string
		if tool.Installed {
			versionInfo := tool.Version
			if len(versionInfo) > 40 {
				versionInfo = versionInfo[:40]
			}
			line = fmt.Sprintf(" %s %-15s  %s", icon, tool.Name, versionInfo)
		} else {
			line = fmt.Sprintf(" %s %-15s  (not installed)", icon, tool.Name)
		}

		padding := 78 - len(line)
		if padding < 0 {
			padding = 0
		}

		fmt.Printf("║ │%s%s│\n", line, strings.Repeat(" ", padding))
	}

	fmt.Println("║ └────────────────────────────────────────────────────────────────────────────┘")
}

func printSuggestions(suggestions []Suggestion) {
	fmt.Println("║ 💡 SUGGESTIONS")
	fmt.Println("║ ┌────────────────────────────────────────────────────────────────────────────┐")

	if len(suggestions) == 0 {
		fmt.Println("║ │ No suggestions at this time                                                │")
	} else {
		for _, sug := range suggestions {
			line := fmt.Sprintf(" • %s", sug.Message)

			// Truncate if too long
			if len(line) > 76 {
				line = line[:73] + "..."
			}

			padding := 78 - len(line)
			if padding < 0 {
				padding = 0
			}

			fmt.Printf("║ │%s%s│\n", line, strings.Repeat(" ", padding))
		}
	}

	fmt.Println("║ └────────────────────────────────────────────────────────────────────────────┘")
}

func printFooter() {
	line := strings.Repeat("═", 80)
	fmt.Println("║")
	fmt.Println("║ Type 'help' for commands, 'dash' to refresh, or 'exit' to quit...")
	fmt.Println("╚" + line + "╝")
}

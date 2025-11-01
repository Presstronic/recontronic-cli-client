package cmd

import (
	"fmt"

	"github.com/presstronic/recontronic-cli-client/pkg/ui"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	Aliases: []string{"dash"},
	Short:   "Display the dashboard",
	Long: `Display the interactive dashboard showing recent activity, statistics,
system status, and actionable suggestions.`,
	RunE: runDashboard,
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
}

func runDashboard(cmd *cobra.Command, args []string) error {
	if err := ui.DisplayDashboard(cfg); err != nil {
		return fmt.Errorf("failed to display dashboard: %w", err)
	}
	return nil
}

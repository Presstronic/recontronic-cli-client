package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display version, commit hash, and build date for the Recontronic CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Recontronic CLI\n")
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Commit:     %s\n", commit)
		fmt.Printf("Built:      %s\n", date)
	},
}

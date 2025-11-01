package cmd

import (
	"fmt"
	"os"

	"github.com/presstronic/recontronic-cli-client/pkg/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	debug   bool
	output  string

	// Global config instance
	cfg *config.Config
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "recon-cli",
	Short: "Recontronic CLI - Bug bounty reconnaissance platform client",
	Long: `Recontronic CLI is a command-line interface for managing continuous
reconnaissance and anomaly detection for bug bounty programs.

The CLI provides tools for:
- User authentication and API key management
- Bug bounty program management
- Reconnaissance scan control and monitoring
- Security anomaly tracking and review
- Real-time dashboards and statistics`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Override with command-line flags if provided
		if output != "" {
			cfg.OutputFormat = output
		}
		if debug {
			cfg.LogLevel = "debug"
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is provided, start interactive mode
		return startInteractiveMode()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.recon-cli/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output format (table|json|yaml)")

	// Add subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
}

// GetConfig returns the loaded configuration
func GetConfig() *config.Config {
	return cfg
}

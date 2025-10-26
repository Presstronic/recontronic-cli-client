package cmd

import (
	"fmt"

	"github.com/presstronic/recontronic-cli-client/pkg/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long:  `View and modify CLI configuration settings.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  server         - Server URL (e.g., http://localhost:8080)
  grpc-server    - gRPC server address (e.g., localhost:9090)
  api-key        - API key for authentication
  timeout        - Request timeout (e.g., 30s, 1m)
  output-format  - Output format (table, json, yaml)
  log-level      - Log level (debug, info, warn, error)`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		if err := config.Set(key, value); err != nil {
			return err
		}

		fmt.Printf("✓ Configuration updated: %s = %s\n", key, value)

		// Show config file location
		configPath, _ := config.GetConfigPath()
		fmt.Printf("  Config file: %s\n", configPath)

		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long:  `Display the current value of a configuration key.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		value, err := config.Get(key)
		if err != nil {
			return err
		}

		// Mask sensitive values
		if (key == "api-key" || key == "api_key") && value != "" {
			if len(value) > 12 {
				value = value[:8] + "..." + value[len(value)-4:]
			}
		}

		fmt.Printf("%s: %s\n", key, value)
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	Long:  `Display all current configuration settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load("")
		if err != nil {
			return err
		}

		fmt.Println("Configuration:")
		fmt.Printf("  server:         %s\n", cfg.Server)
		fmt.Printf("  grpc-server:    %s\n", cfg.GRPCServer)

		// Mask API key
		apiKey := cfg.APIKey
		if apiKey != "" && len(apiKey) > 12 {
			apiKey = apiKey[:8] + "..." + apiKey[len(apiKey)-4:]
		} else if apiKey == "" {
			apiKey = "(not set)"
		}
		fmt.Printf("  api-key:        %s\n", apiKey)

		fmt.Printf("  timeout:        %s\n", cfg.Timeout)
		fmt.Printf("  output-format:  %s\n", cfg.OutputFormat)
		fmt.Printf("  log-level:      %s\n", cfg.LogLevel)

		// Show config file location
		configPath, _ := config.GetConfigPath()
		fmt.Printf("\nConfig file: %s\n", configPath)

		return nil
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long:  `Create a new configuration file with default values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := config.GetConfigPath()
		if err != nil {
			return err
		}

		// Check if config already exists
		force, _ := cmd.Flags().GetBool("force")
		if !force {
			if _, err := config.Load(""); err == nil {
				return fmt.Errorf("configuration file already exists. Use --force to overwrite")
			}
		}

		// Prompt for server URL
		fmt.Print("Server URL [http://localhost:8080]: ")
		var serverURL string
		fmt.Scanln(&serverURL)
		if serverURL == "" {
			serverURL = "http://localhost:8080"
		}

		// Create default config
		cfg := config.DefaultConfig()
		cfg.Server = serverURL

		// Save config
		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("✓ Configuration initialized\n")
		fmt.Printf("  Location: %s\n", configPath)
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  1. Register: recon-cli auth register\n")
		fmt.Printf("  2. Login:    recon-cli auth login\n")

		return nil
	},
}

func init() {
	// Add subcommands
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configInitCmd)

	// Flags for init command
	configInitCmd.Flags().Bool("force", false, "overwrite existing configuration")
}

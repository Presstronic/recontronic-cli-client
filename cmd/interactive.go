package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/presstronic/recontronic-cli-client/pkg/config"
	"github.com/presstronic/recontronic-cli-client/pkg/ui"
	"github.com/spf13/cobra"
)

// startInteractiveMode starts the interactive REPL session
func startInteractiveMode() error {
	// Display dashboard on startup
	if err := ui.DisplayDashboard(cfg); err != nil {
		// Fallback to simple welcome message if dashboard fails
		fmt.Println("Recontronic CLI - Interactive Mode")
		fmt.Println("Type 'help' for available commands, 'exit' or 'quit' to leave")
		fmt.Println()
	}

	// Configure readline with history
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "> ",
		HistoryFile:     os.ExpandEnv("$HOME/.recon-cli/history"),
		HistoryLimit:    20,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer rl.Close()

	for {
		// Read input with history support
		line, err := rl.Readline()
		if err != nil {
			if err == io.EOF || err == readline.ErrInterrupt {
				fmt.Println("\nGoodbye!")
				break
			}
			return fmt.Errorf("error reading input: %w", err)
		}

		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle exit commands
		if line == "exit" || line == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Handle clear screen
		if line == "clear" || line == "cls" {
			fmt.Print("\033[H\033[2J")
			continue
		}

		// Handle dashboard refresh
		if line == "dash" || line == "dashboard" || line == "refresh" {
			if err := ui.DisplayDashboard(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error displaying dashboard: %v\n", err)
			}
			continue
		}

		// Execute the command
		if err := executeInteractiveCommand(line); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		fmt.Println() // Add blank line after command output
	}

	return nil
}

// executeInteractiveCommand executes a command string in the interactive session
func executeInteractiveCommand(input string) error {
	// Parse the input into arguments
	args := parseCommandLine(input)
	if len(args) == 0 {
		return nil
	}

	// Create a new root command for this execution
	// We need to reset the command tree for each execution
	cmd := buildRootCommand()

	// Set the arguments
	cmd.SetArgs(args)

	// Suppress usage on error in interactive mode
	cmd.SilenceUsage = true
	cmd.SilenceErrors = false

	// Execute the command
	return cmd.Execute()
}

// buildRootCommand creates a fresh root command with all subcommands
func buildRootCommand() *cobra.Command {
	cmd := &cobra.Command{
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
	}

	// Add persistent flags
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.recon-cli/config.yaml)")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
	cmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output format (table|json|yaml)")

	// Add all subcommands
	cmd.AddCommand(authCmd)
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(configCmd)
	cmd.AddCommand(reconCmd)
	cmd.AddCommand(dashboardCmd)

	return cmd
}

// parseCommandLine splits a command line into arguments, respecting quotes
func parseCommandLine(input string) []string {
	var args []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range input {
		switch {
		case r == '"' || r == '\'':
			if inQuote {
				if r == quoteChar {
					inQuote = false
					quoteChar = 0
				} else {
					current.WriteRune(r)
				}
			} else {
				inQuote = true
				quoteChar = r
			}
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

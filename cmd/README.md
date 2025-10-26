# Commands (cmd)

This directory contains all CLI command implementations using the Cobra framework.

## Structure

- `root.go` - Root command and global flags
- `program.go` - Program management commands
- `scan.go` - Scan control commands
- `anomalies.go` - Anomaly viewing and management commands
- `dashboard.go` - Interactive TUI dashboard
- `stats.go` - Statistics viewing command
- `config.go` - Configuration management commands

## Adding a New Command

1. Create a new file for your command group
2. Initialize the command using Cobra
3. Register with root command
4. Add tests in corresponding `_test.go` file

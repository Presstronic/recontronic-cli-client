# Recontronic CLI Client

A powerful command-line interface for managing continuous reconnaissance and anomaly detection for bug bounty programs.

## Overview

The Recontronic CLI is the primary interface for interacting with the Recontronic Platform - a continuous reconnaissance and anomaly detection system purpose-built for bug bounty hunting. The CLI provides real-time monitoring, intelligent alerting, and comprehensive program management capabilities.

## Features

- **Program Management**: Add, list, and manage bug bounty programs
- **Scan Control**: Trigger manual scans and monitor progress in real-time
- **Anomaly Tracking**: Query and review detected security anomalies
- **Live Dashboard**: Real-time TUI dashboard with streaming updates
- **Statistics**: Live platform metrics and performance data
- **Configuration**: Simple configuration management for server endpoints and credentials

## Installation

### Prerequisites

- Go 1.21 or higher
- Access to a Recontronic Platform server (REST API and gRPC endpoints)
- API key for authentication

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/recontronic-cli-client.git
cd recontronic-cli-client

# Build the binary (using Makefile - recommended)
make build

# Or build manually
go build -o recon-cli main.go

# Install to GOPATH/bin (using Makefile)
make install

# Or move to your PATH manually
sudo mv recon-cli /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/yourusername/recontronic-cli-client@latest
```

## Quick Start

### 1. Configure the CLI

```bash
# Set server endpoints
recon-cli config set server http://your-server:8080
recon-cli config set grpc-server your-server:9090
recon-cli config set api-key your-api-key-here
```

### 2. Add a Program

```bash
recon-cli program add \
  --name "Example Corp" \
  --platform hackerone \
  --scope "*.example.com,*.example.io" \
  --frequency 1h
```

### 3. Trigger a Scan

```bash
recon-cli scan trigger --program-id 1 --type passive
```

### 4. View Anomalies

```bash
recon-cli anomalies list --min-priority 70 --unreviewed
```

### 5. Launch Dashboard

```bash
recon-cli dashboard
```

## Usage

### Program Commands

```bash
# Add a new program
recon-cli program add --name "Company" --platform hackerone --scope "*.example.com"

# List all programs
recon-cli program list

# Get program details
recon-cli program get --id 1

# Delete a program
recon-cli program delete --id 1
```

### Scan Commands

```bash
# Trigger a new scan
recon-cli scan trigger --program-id 1 --type passive

# Watch scan progress
recon-cli scan watch --scan-id 42

# List recent scans
recon-cli scan list --program-id 1 --limit 10
```

### Anomaly Commands

```bash
# List anomalies
recon-cli anomalies list

# List high-priority unreviewed anomalies
recon-cli anomalies list --min-priority 80 --unreviewed

# View anomaly details
recon-cli anomalies view --id 12

# Mark anomaly as reviewed
recon-cli anomalies review --id 12 --notes "Investigated, false positive"
```

### Dashboard & Monitoring

```bash
# Launch interactive dashboard
recon-cli dashboard

# View live statistics
recon-cli stats

# Stream anomalies in real-time
recon-cli anomalies stream --min-priority 70
```

### Configuration Commands

```bash
# Set configuration values
recon-cli config set <key> <value>

# Get configuration value
recon-cli config get <key>

# List all configuration
recon-cli config list

# Initialize config file
recon-cli config init
```

## Configuration

The CLI stores configuration in `~/.recon-cli/config.yaml`:

```yaml
server: http://localhost:8080
grpc_server: localhost:9090
api_key: your-api-key-here
timeout: 30s
output_format: table  # table, json, yaml
log_level: info
```

### Environment Variables

Configuration can also be set via environment variables:

```bash
export RECON_SERVER="http://localhost:8080"
export RECON_GRPC_SERVER="localhost:9090"
export RECON_API_KEY="your-api-key"
```

## Development

### Project Structure

```
recontronic-cli-client/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command
│   ├── program.go         # Program subcommands
│   ├── scan.go            # Scan subcommands
│   ├── anomalies.go       # Anomaly subcommands
│   ├── dashboard.go       # Dashboard TUI
│   ├── stats.go           # Stats command
│   └── config.go          # Config management
├── pkg/                   # Reusable packages
│   ├── client/           # API clients
│   │   ├── rest.go       # REST API client
│   │   └── grpc.go       # gRPC client
│   ├── config/           # Configuration handling
│   │   └── config.go
│   ├── ui/               # User interface components
│   │   └── dashboard.go  # Bubble Tea TUI
│   └── models/           # Data models
│       └── types.go
├── proto/                # Protocol buffer definitions
│   └── recon/v1/
├── scripts/              # Build and utility scripts
├── docs/                 # Documentation
├── main.go              # Entry point
└── go.mod
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests (requires server)
go test -tags=integration ./...
```

## Troubleshooting

### Connection Issues

If you're having trouble connecting to the server:

1. Verify server endpoint: `recon-cli config get server`
2. Test connectivity: `curl http://your-server:8080/health`
3. Check API key is set: `recon-cli config get api-key`
4. Verify firewall rules allow outbound connections

### Authentication Errors

- Ensure your API key is valid
- Check if the key has expired
- Verify the key has proper permissions

### gRPC Streaming Issues

- Ensure gRPC port (9090) is accessible
- Check for firewall blocking gRPC traffic
- Verify TLS/SSL configuration matches server

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines and how to submit contributions.

## License

MIT License - See [LICENSE](LICENSE) for details.

## Support

- Documentation: [docs/](docs/)
- Issues: [GitHub Issues](https://github.com/yourusername/recontronic-cli-client/issues)
- Vision Document: [early-vision-doc.md](early-vision-doc.md)

## Related Projects

- [Recontronic Platform Server](https://github.com/yourusername/recontronic-platform) - Backend platform
- [Recontronic Web UI](https://github.com/yourusername/recontronic-web) - Web interface (future)

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [gRPC](https://grpc.io/) - RPC framework

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Recontronic CLI is a command-line interface for managing continuous reconnaissance and anomaly detection for bug bounty programs. It's written in Go and interfaces with the Recontronic Platform server via REST API and gRPC.

**Current Status:**
- Authentication system fully implemented (register, login, API key management)
- Interactive REPL mode with dashboard
- Activity logging and smart suggestions
- Reconnaissance tools: subdomain enumeration, HTTP verification, WHOIS lookup

## Build & Development Commands

**NOTE:** The project now includes a comprehensive Makefile. Run `make help` to see all available targets.

### Building with Makefile (Recommended)
```bash
# Build for current platform
make build

# Build for all platforms (Linux, macOS, Windows)
make build-all

# Build for specific platform
make build-linux
make build-darwin
make build-windows

# Clean build artifacts
make clean

# Install to GOPATH/bin
make install
```

### Testing with Makefile (Recommended)
```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Run fmt, vet, and test
make check
```

### Code Quality
```bash
# Format code
make fmt

# Run go vet
make vet

# Run golangci-lint (if installed)
make lint
```

### Manual Build Commands (Alternative)
```bash
# Build for current platform
go build -o recon-cli main.go

# Run without building
go run main.go [command]
make run  # Or use Makefile

# Clean build artifacts
rm -f recon-cli

# Create sample data for dashboard testing
./scripts/create-sample-data.sh
```

### Manual Testing (Alternative)
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests (requires server running)
go test -tags=integration ./...
```

### Running
```bash
# After building
./recon-cli [command]

# Examples - Authentication
./recon-cli auth register
./recon-cli auth login
./recon-cli auth whoami
./recon-cli config set server http://localhost:8080

# Examples - Reconnaissance
./recon-cli recon subdomain example.com
./recon-cli recon verify example.com
./recon-cli recon whois example.com
./recon-cli recon results list
./recon-cli recon results export example.com --format csv
```

## Architecture

### Interactive Dashboard

**Location:** `pkg/ui/dashboard.go`, `pkg/ui/stats.go`, `pkg/ui/activity.go`, `pkg/ui/system.go`, `pkg/ui/suggestions.go`

The CLI features an interactive dashboard that displays on startup, providing immediate context about recent activity and system status.

**Dashboard Components:**

1. **Header Bar** - Server connection, auth status, tool availability count
2. **Quick Statistics** - Domains scanned, subdomains found, alive targets, recent scans, storage usage
3. **Recent Activity** - Last 5 activities with human-readable timestamps (e.g., "5m ago", "2h ago")
4. **System Status** - Installed tools detection (crt.sh, subfinder, amass, assetfinder, httpx, nuclei)
5. **Smart Suggestions** - Actionable recommendations based on workspace state

**Activity Logging:**
- Location: `~/.recon-cli/activity.log`
- Format: JSON Lines (one JSON object per line)
- Auto-logged: All recon operations (subdomain enum, verify, dns, whois)
- Functions: `ui.LogActivity()`, `ui.GetRecentActivity()`

**Statistics Gathering:**
- Scans results directory: `~/.recon-cli/results/`
- Parses JSON result files
- Counts: domains, subdomains, alive targets, recent scans
- Calculates storage usage

**Smart Suggestions Engine:**
- Detects unverified subdomains → suggests running verify
- Detects old scans (>7 days) → suggests re-scanning
- Detects missing tools → suggests installation
- Priority levels: 1=high, 2=medium, 3=low

**Usage:**
```bash
# Dashboard auto-displays on startup
./recon-cli

# Manual refresh in interactive mode
> dash
> dashboard
> refresh

# Direct command
./recon-cli dashboard
```

### Command Structure (Cobra Framework)

**Entry Point:** `main.go` → `cmd.Execute()` → Cobra command tree

**Root Command:** `cmd/root.go` defines the base `recon-cli` command with global flags:
- `--config`: Override default config file path
- `--debug`: Enable debug logging
- `--output, -o`: Set output format (table|json|yaml)

**PersistentPreRunE Hook:** Loads configuration before every command runs and stores it in the global `cfg` variable (accessible via `GetConfig()`).

**Subcommands:**
- `cmd/auth.go` - Authentication commands (register, login, whoami, keys)
- `cmd/config_cmd.go` - Configuration management
- `cmd/version.go` - Version information
- `cmd/dashboard.go` - Interactive dashboard display
- `cmd/interactive.go` - REPL mode with command parser
- `cmd/recon.go` - Reconnaissance command root
- `cmd/recon_subdomain.go` - Subdomain enumeration (integrated into recon.go)
- `cmd/recon_verify.go` - HTTP/HTTPS verification
- `cmd/recon_whois.go` - WHOIS domain lookup
- `cmd/recon_results.go` - Results management (list, view, export)

### Configuration System

**Location:** `~/.recon-cli/config.yaml`

**Implementation:** `pkg/config/config.go` using Viper for configuration management

**Environment Variables:** All config values can be overridden with `RECON_` prefix:
- `RECON_SERVER`
- `RECON_GRPC_SERVER`
- `RECON_API_KEY`

**Security:** Config files are automatically set to `0600` permissions (owner read/write only) because they contain API keys.

**Key Functions:**
- `Load(cfgFile string)` - Loads configuration from file or uses defaults
- `Save(cfg *Config)` - Writes configuration to file with secure permissions
- `Set(key, value string)` - Updates a single config value
- `SaveAPIKey(apiKey string)` - Specifically for saving API keys after login

### REST Client Architecture

**Location:** `pkg/client/rest.go`

**Design Pattern:** The `RestClient` struct wraps HTTP operations with:
- Base URL management (trailing slash removal)
- API key authentication (Bearer token in Authorization header)
- Timeout configuration
- Debug logging (sanitizes API keys in output)
- Consistent error handling via `APIError` type

**Authentication Flow:**
1. Unauthenticated endpoints (register, login): `authenticated=false` in `doRequest()`
2. Login returns API key in `LoginResponse.APIKey`
3. API key is saved to config via `config.SaveAPIKey()`
4. Subsequent requests use `authenticated=true` to include Bearer token

**Error Classification Helpers:**
- `IsAuthError(err)` - 401 Unauthorized
- `IsNotFoundError(err)` - 404 Not Found
- `IsValidationError(err)` - 400 Bad Request

### Data Models

**Location:** `pkg/models/types.go`

All API request/response structures are defined here:
- `User`, `APIKey` - Authentication models
- `RegisterRequest`, `LoginRequest`, `LoginResponse` - Auth payloads
- `Program`, `Scan`, `Anomaly` - Future features (server not yet implemented)

**API Key Format:** All API keys use the prefix `rct_` (e.g., `rct_AbCdEf123456789...`)

### Reconnaissance Package

**Location:** `pkg/recon/`

Provides standalone reconnaissance tools for passive information gathering:

**Core Files:**
- `executor.go` - Safe command execution with timeouts and context cancellation
- `storage.go` - JSON file storage in `~/.recon-cli/results/<domain>/`
- `parser.go` - Domain validation, cleaning, deduplication, and sorting
- `subdomain.go` - Multi-source subdomain enumeration (crt.sh, subfinder, amass, assetfinder)
- `verify.go` - DNS resolution and HTTP/HTTPS verification with title extraction
- `whois.go` - WHOIS domain information lookup and parsing
- `results.go` - Results management, filtering, and querying

**Key Features:**
- **Subdomain Enumeration:** Multi-source passive discovery with tool detection
- **HTTP Verification:** Concurrent DNS + HTTP probing with status codes and HTML titles
- **WHOIS Lookup:** Domain registration info (registrar, dates, nameservers, status)
- **Results Management:** List, view, filter, and export results in CSV/JSON/Markdown
- **Storage:** JSON files with timestamps, organized by domain
- **Activity Logging:** All operations logged to `~/.recon-cli/activity.log`

**WHOIS Implementation:**
- Executes system `whois` command with configurable timeout (default: 30s)
- Parses output to extract: registrar, creation/expiry/update dates, nameservers, domain status
- Deduplicates nameservers and status entries (case-insensitive)
- Supports `--json` flag for raw JSON output, `--raw` for unparsed WHOIS data
- Results saved to `~/.recon-cli/results/<domain>/whois_<timestamp>.json`

**Example Usage:**
```bash
# WHOIS lookup
./recon-cli recon whois example.com
./recon-cli recon whois example.com --json
./recon-cli recon whois example.com --timeout 60s
```

### User Input System

**Location:** `pkg/ui/input.go`

**Key Functions:**
- `ReadPassword(prompt)` - Reads password without echo (uses golang.org/x/term)
- `ReadPasswordWithConfirm(prompt, confirmPrompt)` - Password with confirmation
- `ReadInput(prompt)` - Standard text input
- `Confirm(prompt)` - Yes/no confirmation

**Terminal Detection:** `ReadPassword()` detects if stdin is a terminal and falls back to normal input for testing/piping.

**Validation Functions:**
- `ValidateEmail(email)` - Basic email format check
- `ValidateUsername(username)` - 3-50 chars, alphanumeric only
- `ValidatePassword(password)` - 8-72 characters

## Server Integration

**Server Repository:** https://github.com/Presstronic/recontronic-server

**API Base URL:** `/api/v1/auth/` for authentication endpoints

**Authentication Endpoints:**
- `POST /api/v1/auth/register` - User registration (public)
- `POST /api/v1/auth/login` - User login, returns API key (public)
- `GET /api/v1/auth/me` - Get current user (requires auth)
- `POST /api/v1/auth/keys` - Create new API key (requires auth)
- `GET /api/v1/auth/keys` - List API keys (requires auth)
- `DELETE /api/v1/auth/keys/{id}` - Revoke API key (requires auth)

**Not Yet Implemented on Server:**
- Program management endpoints
- Scan management endpoints
- Anomaly detection endpoints
- gRPC streaming

**Testing Against Local Server:**
```bash
# In server repository
docker-compose up -d
make run  # Server runs on http://localhost:8080

# In CLI
./recon-cli config set server http://localhost:8080
./recon-cli auth register
./recon-cli auth login
```

## Code Patterns & Conventions

### Command Implementation Pattern

Every command follows this structure:

```go
var myCmd = &cobra.Command{
    Use:   "mycommand [args]",
    Short: "One-line description",
    Long:  `Multi-line description with usage details`,
    RunE:  runMyCommand,  // Use RunE for error handling
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    ctx := context.Background()

    // Get global config
    cfg := GetConfig()

    // Create REST client
    restClient := client.NewRestClient(cfg.Server, cfg.APIKey, cfg.Timeout)
    if debug {
        restClient.SetDebug(true)
    }

    // Make API call
    result, err := restClient.SomeMethod(ctx, ...)
    if err != nil {
        // Use typed error checking
        if client.IsAuthError(err) {
            return fmt.Errorf("authentication failed: ...")
        }
        return fmt.Errorf("operation failed: %w", err)
    }

    // Format output
    fmt.Printf("Success: %s\n", result.Field)
    return nil
}
```

### Error Handling Pattern

**User-Facing Errors:** Commands return descriptive errors with actionable guidance:
```go
if cfg.APIKey == "" {
    return fmt.Errorf("not authenticated: please run 'recon-cli auth login' first")
}
```

**API Error Classification:**
```go
if err != nil {
    if client.IsAuthError(err) {
        return fmt.Errorf("authentication failed: your API key may be invalid or expired\nPlease run 'recon-cli auth login' to get a new key")
    }
    if client.IsValidationError(err) {
        return fmt.Errorf("validation failed - please check your inputs: %w", err)
    }
    return fmt.Errorf("operation failed: %w", err)
}
```

### Output Formatting

**Tabular Output:** Use `text/tabwriter` for aligned columns:
```go
w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
fmt.Fprintln(w, "ID\tNAME\tSTATUS")
fmt.Fprintln(w, "──\t────\t──────")
for _, item := range items {
    fmt.Fprintf(w, "%d\t%s\t%s\n", item.ID, item.Name, item.Status)
}
w.Flush()
```

**Time Formatting:** Use helper functions like `formatTimeAgo()` and `formatExpiresAt()` in `cmd/auth.go` for human-readable relative times.

### Flag Management

**Persistent Flags:** Defined in `cmd/root.go`, available to all commands (config, debug, output)

**Command-Specific Flags:** Defined in each command's `init()` function:
```go
func init() {
    myCmd.Flags().StringVarP(&flagVar, "name", "n", "", "Description")
    myCmd.Flags().BoolVar(&anotherFlag, "force", false, "Description")
}
```

## Testing Strategy

### Unit Tests
- Test individual functions in isolation
- Mock external dependencies (HTTP client, filesystem)
- Use table-driven tests for multiple scenarios

### Integration Tests
- Tag with `// +build integration`
- Require `RECON_SERVER` environment variable
- Test against live server instance
- See example in `AUTHENTICATION-INTEGRATION-SUMMARY.md:304`

## Security Considerations

**API Key Storage:**
- Config file permissions: `0600` (owner read/write only)
- Config directory permissions: `0700` (owner only)
- Validation: Keys must start with `rct_` and be at least 20 chars

**Password Handling:**
- Never echo passwords to terminal
- Use `golang.org/x/term` for secure input
- Always confirm passwords during registration
- Validate length (8-72 chars) before sending to server

**Debug Output Sanitization:**
- API keys are truncated in debug logs: `rct_AbCd...9876` (first 8 + last 4)
- Implemented in `pkg/client/rest.go:74-81`

## Common Development Tasks

### Adding a New Command

1. Create command in `cmd/mycommand.go`
2. Define `cobra.Command` with `Use`, `Short`, `Long`, `RunE`
3. Add flags in `init()` function
4. Register in `cmd/root.go`: `rootCmd.AddCommand(myCmd)`
5. Implement `runMyCommand()` function following the pattern above

### Adding a New API Endpoint

1. Add method to `pkg/client/rest.go`
2. Define request/response models in `pkg/models/types.go`
3. Use `doRequest()` helper with appropriate authentication flag
4. Add error type checking as needed

### Updating Configuration Schema

1. Update `Config` struct in `pkg/config/config.go`
2. Add defaults in `DefaultConfig()` or `Load()`
3. Add case in `Set()` for command-line updates
4. Add case in `Get()` for retrieval
5. Update validation if needed

## Project Dependencies

**Core Libraries:**
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration management
- `golang.org/x/term` - Terminal password input

**Future Libraries (not yet used):**
- `github.com/charmbracelet/bubbletea` - TUI framework for dashboard
- gRPC libraries - Real-time streaming

## Module Information

**Module Path:** `github.com/presstronic/recontronic-cli-client`

**Go Version:** 1.25.3

**Import Paths:**
- Commands: `github.com/presstronic/recontronic-cli-client/cmd`
- Config: `github.com/presstronic/recontronic-cli-client/pkg/config`
- Client: `github.com/presstronic/recontronic-cli-client/pkg/client`
- Models: `github.com/presstronic/recontronic-cli-client/pkg/models`
- UI: `github.com/presstronic/recontronic-cli-client/pkg/ui`

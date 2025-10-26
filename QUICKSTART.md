# Quick Start Guide for Developers

This guide will help you get started developing the Recontronic CLI Client.

## Prerequisites

- Go 1.21 or higher installed
- Git
- Basic understanding of Go and CLI development
- (Optional) golangci-lint for code quality checks

## Initial Setup

### 1. Verify Go Installation

```bash
go version
# Should output: go version go1.21.x or higher
```

### 2. Project Structure

The project structure is already created:

```
recontronic-cli-client/
├── cmd/                    # Cobra commands (empty - to be implemented)
├── pkg/                    # Reusable packages (empty - to be implemented)
│   ├── client/            # REST and gRPC clients
│   ├── config/            # Configuration management
│   ├── models/            # Data models
│   └── ui/                # UI components
├── proto/                  # Protocol buffer definitions
├── scripts/                # Build and utility scripts
├── docs/                   # Additional documentation
├── test/                   # Integration tests
├── examples/               # Usage examples
├── README.md              # Main documentation
├── CONTRIBUTING.md        # Contribution guidelines
├── mvp-issues.csv         # 50 MVP issues/stories
└── MVP-ISSUES-SUMMARY.md  # Issue overview
```

## Development Workflow

### Step 1: Initialize Go Module (Issue RECON-001)

This is your first task. You need to:

```bash
# Initialize Go module
go mod init github.com/yourusername/recontronic-cli-client

# This creates go.mod file
```

### Step 2: Install Core Dependencies

```bash
# Install Cobra for CLI framework
go get github.com/spf13/cobra@latest

# Install Viper for configuration
go get github.com/spf13/viper@latest

# Install Bubble Tea for TUI (Terminal UI)
go get github.com/charmbracelet/bubbletea@latest

# Install gRPC
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf/cmd/protoc-gen-go@latest
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Step 3: Create main.go

Create a basic `main.go` file:

```go
package main

import (
    "fmt"
    "os"
)

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    // TODO: Initialize root command
    fmt.Println("Recontronic CLI Client")
    fmt.Printf("Version: %s\nCommit: %s\nBuilt: %s\n", version, commit, date)
    os.Exit(0)
}
```

Test it works:

```bash
go run main.go
# Should print version info
```

### Step 4: Create Makefile (Issue RECON-029)

Create a `Makefile` for common tasks:

```makefile
.PHONY: build test clean install lint

# Binary name
BINARY_NAME=recon-cli

# Version info
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse HEAD)
DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Ldflags for version injection
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

# Build the binary
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} main.go

# Install to $GOPATH/bin
install:
	go install ${LDFLAGS}

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -f coverage.out coverage.html

# Run the application
run:
	go run main.go

# Build for all platforms
build-all:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe
```

Test the Makefile:

```bash
make build
./recon-cli
```

## Issue Workflow

### Understanding the Issues

1. Open `mvp-issues.csv` - contains all 50 issues
2. Review `MVP-ISSUES-SUMMARY.md` - provides overview and implementation order
3. Start with Phase 1 (Foundation) issues

### Recommended First 5 Issues

#### 1. RECON-001: Initialize Go Module (DONE ABOVE)
- ✅ Go module created
- ✅ Basic main.go created
- ✅ Build works

#### 2. RECON-029: Create Makefile (DONE ABOVE)
- ✅ Makefile created
- ✅ Build, test, clean targets work

#### 3. RECON-003: Implement Root Command with Cobra

Create `cmd/root.go`:

```go
package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    cfgFile string
    debug   bool
    output  string
)

var rootCmd = &cobra.Command{
    Use:   "recon-cli",
    Short: "Recontronic CLI - Bug bounty reconnaissance platform client",
    Long:  `A powerful CLI for managing continuous reconnaissance and anomaly detection.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.recon-cli/config.yaml)")
    rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
    rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "table", "output format (table|json|yaml)")
}

func initConfig() {
    // TODO: Implement in RECON-002
}
```

Update `main.go`:

```go
package main

import "github.com/yourusername/recontronic-cli-client/cmd"

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    cmd.Execute()
}
```

#### 4. RECON-002: Configuration Management

Create `pkg/config/config.go`:

```go
package config

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/viper"
)

type Config struct {
    Server      string `mapstructure:"server"`
    GRPCServer  string `mapstructure:"grpc_server"`
    APIKey      string `mapstructure:"api_key"`
    Timeout     string `mapstructure:"timeout"`
    OutputFormat string `mapstructure:"output_format"`
    LogLevel    string `mapstructure:"log_level"`
}

func Load(cfgFile string) (*Config, error) {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        if err != nil {
            return nil, err
        }

        configDir := filepath.Join(home, ".recon-cli")
        viper.AddConfigPath(configDir)
        viper.SetConfigType("yaml")
        viper.SetConfigName("config")
    }

    // Environment variable support
    viper.SetEnvPrefix("RECON")
    viper.AutomaticEnv()

    // Defaults
    viper.SetDefault("timeout", "30s")
    viper.SetDefault("output_format", "table")
    viper.SetDefault("log_level", "info")

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

#### 5. RECON-008: Define Data Models

Create `pkg/models/types.go`:

```go
package models

import "time"

type Program struct {
    ID            int64     `json:"id"`
    Name          string    `json:"name"`
    Platform      string    `json:"platform"`
    Scope         []string  `json:"scope"`
    ScanFrequency string    `json:"scan_frequency"`
    CreatedAt     time.Time `json:"created_at"`
    LastScannedAt *time.Time `json:"last_scanned_at,omitempty"`
    IsActive      bool      `json:"is_active"`
    Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

type Scan struct {
    ID            int64     `json:"id"`
    ProgramID     int64     `json:"program_id"`
    ScanType      string    `json:"scan_type"`
    Status        string    `json:"status"`
    Progress      int       `json:"progress"`
    AssetsFound   int       `json:"assets_found"`
    StartedAt     *time.Time `json:"started_at,omitempty"`
    CompletedAt   *time.Time `json:"completed_at,omitempty"`
}

type Anomaly struct {
    ID                  int64     `json:"id"`
    ProgramID           int64     `json:"program_id"`
    ProgramName         string    `json:"program_name"`
    Type                string    `json:"type"`
    Description         string    `json:"description"`
    PriorityScore       float64   `json:"priority_score"`
    DetectedAt          time.Time `json:"detected_at"`
    IsReviewed          bool      `json:"is_reviewed"`
    Metadata            map[string]interface{} `json:"metadata,omitempty"`
}
```

## Testing Your Progress

After implementing the first few issues:

```bash
# Build the project
make build

# Run the CLI
./recon-cli

# Should show help text from Cobra
./recon-cli --help

# Try version command (once implemented)
./recon-cli version
```

## Getting Help

1. **Documentation**: Check README.md and CONTRIBUTING.md
2. **Issues**: Refer to mvp-issues.csv for detailed acceptance criteria
3. **Vision**: Review early-vision-doc.md for overall context
4. **Summary**: Check MVP-ISSUES-SUMMARY.md for implementation order

## Recommended Reading

Before starting development:

1. [Cobra Documentation](https://cobra.dev/)
2. [Viper Documentation](https://github.com/spf13/viper)
3. [Bubble Tea Tutorial](https://github.com/charmbracelet/bubbletea)
4. [Go Project Layout](https://github.com/golang-standards/project-layout)

## Common Mistakes to Avoid

1. ❌ Don't implement features out of order (follow dependencies)
2. ❌ Don't skip writing tests
3. ❌ Don't commit sensitive data (API keys, credentials)
4. ❌ Don't ignore linting errors
5. ❌ Don't forget to update documentation

## Success Criteria

You'll know you're on track when:

✅ `go build` succeeds without errors
✅ `make test` runs (even if tests are minimal)
✅ `./recon-cli --help` shows useful output
✅ Code passes linting (`make lint`)
✅ Directory structure is clean and organized

## Next Steps

1. ✅ You've read this guide
2. ⬜ Initialize Go module (RECON-001)
3. ⬜ Create Makefile (RECON-029)
4. ⬜ Implement root command (RECON-003)
5. ⬜ Implement configuration (RECON-002)
6. ⬜ Define data models (RECON-008)
7. ⬜ Continue with Phase 1 issues...

---

**Ready to start?** Begin with RECON-001 and work your way through the issues in order!

Good luck! 🚀

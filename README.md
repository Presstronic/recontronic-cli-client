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

## Bug Bounty Reconnaissance Workflow

The CLI provides a complete reconnaissance workflow for bug bounty hunting. Follow these steps in order for best results:

### Phase 1: Discovery ✅ IMPLEMENTED

#### Step 1: Subdomain Enumeration
**Find all subdomains for your target domain**

```bash
# Basic subdomain enumeration
./recon-cli recon subdomain example.com

# With custom timeout for slow tools
./recon-cli recon subdomain example.com --timeout 15m
```

**Sample Output:**
```
Finding subdomains for example.com
Mode: Passive reconnaissance (safe, no active scanning)

Running crt.sh...
✓ crt.sh found 445 subdomains (3.2s)

Running subfinder...
✓ subfinder found 789 subdomains (12.4s)

Running assetfinder...
✓ assetfinder found 234 subdomains (8.1s)

Summary:
  Total unique subdomains: 808
  Sources used: crt.sh, subfinder, assetfinder

✓ Results saved to ~/.recon-cli/results/example.com/subdomains_20251103_120534.json
```

- **Tools Used:** crt.sh, subfinder, amass, assetfinder
- **Output:** List of all discovered subdomains
- **Saved To:** `~/.recon-cli/results/example.com/subdomains_*.json`
- **Typical Results:** 100-1000+ subdomains depending on target size

#### Step 2: HTTP/HTTPS Verification
**Find which subdomains are alive and accessible**

```bash
# Basic verification (default: 10 concurrent, 10s timeout)
./recon-cli recon verify example.com

# Faster scanning with higher concurrency
./recon-cli recon verify example.com --concurrency 50 --timeout 5s

# Conservative scanning (slower but more reliable)
./recon-cli recon verify example.com --concurrency 5 --timeout 30s
```

**Sample Output:**
```
Verifying subdomains for example.com
Mode: Passive verification (DNS + HTTP probing)

Progress: 156/808 verified (19.3%) [29.4s elapsed]

Summary:
  Total subdomains: 808
  Alive: 156 (19.3%)
  Dead: 652 (80.7%)
  Duration: 45.2s

Sample alive hosts:
  https://www.example.com - 200 OK - "Example Domain - Official Site"
  https://api.example.com - 403 Forbidden - "Access Denied"
  https://mail.example.com - 200 OK - "Webmail Login"
  https://admin.example.com - 401 Unauthorized - "Authentication Required"
  http://dev.example.com - 200 OK - "Development Server"

✓ Results updated in ~/.recon-cli/results/example.com/subdomains_20251103_120534.json
```

- **What It Does:**
  - DNS resolution check
  - HTTP/HTTPS probing (tries HTTPS first, falls back to HTTP)
  - Extracts HTTP status codes
  - Captures HTML page titles
  - Measures response times
- **Output:** List of alive hosts with status codes and titles
- **Typical Results:** 10-30% of subdomains are usually alive

#### Step 3: WHOIS Lookup
**Get domain registration and infrastructure information**

```bash
# Basic WHOIS lookup
./recon-cli recon whois example.com

# Output as JSON for parsing
./recon-cli recon whois example.com --json

# Show raw WHOIS output
./recon-cli recon whois example.com --raw

# Custom timeout (default: 30s)
./recon-cli recon whois example.com --timeout 60s
```

**Sample Output:**
```
Looking up WHOIS information for example.com
Mode: Passive reconnaissance (WHOIS query)

✓ Results saved to ~/.recon-cli/results/example.com/

Domain: example.com
Registrar: MarkMonitor Inc.
Created: 1995-08-14T04:00:00Z
Updated: 2024-08-13T07:01:38Z
Expires: 2025-08-13T04:00:00Z

Name Servers:
  - a.iana-servers.net
  - b.iana-servers.net

Status:
  - clientDeleteProhibited
  - clientTransferProhibited
  - clientUpdateProhibited
  - serverDeleteProhibited
  - serverTransferProhibited
  - serverUpdateProhibited

Registrar URL: http://www.markmonitor.com
```

- **Information Gathered:**
  - Registrar details
  - Creation, update, and expiry dates
  - Authoritative nameservers
  - Domain status (locked, unlocked, etc.)
- **Why It Matters:** Helps validate scope and understand domain infrastructure

#### Step 4: View & Export Results
**Organize and export your findings**

```bash
# List all results for all domains
./recon-cli recon results list

# View all subdomains for a specific domain
./recon-cli recon results view example.com

# View only alive subdomains
./recon-cli recon results view example.com --alive-only

# Filter by HTTP status code
./recon-cli recon results view example.com --status 200

# Filter by discovery source
./recon-cli recon results view example.com --source subfinder

# Limit results
./recon-cli recon results view example.com --alive-only --limit 50

# Export to CSV (great for spreadsheet analysis)
./recon-cli recon results export example.com --format csv --alive-only

# Export to JSON (for tool integration)
./recon-cli recon results export example.com --format json --alive-only

# Export to Markdown (for reporting)
./recon-cli recon results export example.com --format markdown

# Export with custom output path
./recon-cli recon results export example.com --format csv --output ~/reports/example.csv

# Combine multiple filters
./recon-cli recon results export example.com --format csv --alive-only --status 200
```

**Sample Output (list):**
```
Results for all domains:

example.com/
  2025-11-03 12:05  subdomains  (808 total, 156 alive)  ✓ verified
  2025-11-03 12:15  whois       Registrar: MarkMonitor Inc.

tesla.com/
  2025-11-01 01:40  subdomains  (808 total)  ⚠ not verified

basecamp.com/
  2025-11-01 02:50  subdomains  (133 total, 23 alive)  ✓ verified
```

**Sample Output (export CSV):**
```
Exporting results for example.com...
Format: CSV
Filters: alive-only

✓ Exported 156 subdomains to ~/.recon-cli/exports/example.com_20251103_121830.csv
```

---

### Phase 2: Deep Enumeration ✅ IMPLEMENTED

#### Step 5: DNS Enumeration
**Get detailed DNS records for all alive subdomains**

```bash
# DNS enumeration for all alive subdomains (default)
./recon-cli recon dns example.com

# Query specific record types
./recon-cli recon dns example.com --types A,AAAA,MX,TXT

# DNS enumeration for all subdomains (not just alive)
./recon-cli recon dns example.com --alive-only=false

# Export DNS results to CSV
./recon-cli recon results export example.com --type dns --format csv

# High-speed scanning
./recon-cli recon dns example.com --concurrency 50 --timeout 3s

# Check for subdomain takeover opportunities (enabled by default)
./recon-cli recon dns example.com --check-takeover
```

**Sample Output:**
```
Enumerating DNS records for basecamp.com
Mode: Passive DNS enumeration

✓ Results saved to ~/.recon-cli/results/basecamp.com/

Summary:
  Subdomains queried: 23
  A records: 50
  AAAA records: 0
  CNAME records: 3
  MX records: 3
  TXT records: 8
  NS records: 10
  Unique IPs: 18
  Duration: 3s

Key Findings:
  ✓ No obvious subdomain takeover risks detected
  ☁️  Cloud providers detected: Cloudflare
  📧 Mail servers found: 3 MX records
      Providers: basecamp.com
  🔒 Security records: SPF (yes), DMARC (no), DKIM (no)

Sample DNS Records:
  SUBDOMAIN                 RECORD TYPE  VALUE                     CLOUD
  3.basecamp.com            A            104.18.12.81              Cloudflare
  www.updates.basecamp.com  CNAME        ext-cust.squarespace.com
  storage.basecamp.com      A            104.18.17.127             Cloudflare
  ... and 13 more records (see JSON results for complete data)
```

**What This Provides:**
- **A/AAAA Records:** IP addresses (IPv4/IPv6) - map subdomains to actual hosts for port scanning
- **MX Records:** Mail servers - identify email infrastructure targets
- **TXT Records:** SPF, DMARC, DKIM verification records - find security misconfigurations
- **NS Records:** Authoritative nameservers - understand DNS infrastructure
- **CNAME Records:** Subdomain aliases - **detect potential subdomain takeover opportunities** 💰
- **Cloud Providers:** Automatic identification (AWS, Azure, GCP, Cloudflare, Akamai, Fastly)
- **Takeover Detection:** Checks for 15+ vulnerable services (herokuapp, github.io, s3, azurewebsites, etc.)
- **Security Analysis:** Detects SPF, DMARC, DKIM configurations

**Why This Matters:**
- 🎯 Maps subdomains to IP addresses ready for port scanning
- 🔍 Identifies shared infrastructure (multiple domains on same IP = similar attack surface)
- ☁️ Discovers cloud providers (AWS, Azure, GCP = different security models)
- 💰 Finds dangling CNAMEs = potential subdomain takeovers (HIGH/CRITICAL severity)
- 📧 Identifies mail infrastructure for email security testing
- 🚨 Detects security misconfigurations (missing DMARC = email spoofing)

---

### Phase 3: Active Scanning 📋 PLANNED

#### Step 6: Port Scanning
**Identify open ports and running services**

```bash
# Port scanning (PLANNED - not yet implemented)
# Use external tools for now:

# Fast port scan with naabu (recommended)
cat ~/.recon-cli/exports/example.com_alive.txt | naabu -p 80,443,8080,8443 -o ports.txt

# Full port scan with nmap
nmap -iL alive_hosts.txt -p- -oA nmap_results

# Quick common ports scan
masscan -iL alive_ips.txt -p 80,443,8080,8443,3000,8000,9000 --rate 1000
```

**Common Targets:**
- **Web:** 80, 443, 8080, 8443, 8000, 3000
- **Admin panels:** 9000, 10000
- **Databases:** 3306 (MySQL), 5432 (PostgreSQL), 27017 (MongoDB)
- **APIs:** 8081, 8082, 9090
- **Security Note:** ⚠️ Only scan assets within scope!

#### Step 7: Technology Detection
**Identify frameworks, libraries, and technologies**

```bash
# Technology detection (PLANNED - not yet implemented)
# Use external tools for now:

# Using httpx with technology detection
cat ~/.recon-cli/exports/example.com_alive.txt | httpx -td -title -status-code -o tech_stack.txt

# Using wappalyzer
wappalyzer https://example.com

# Check headers and identify technologies
curl -I https://example.com
```

**What You'll Discover:**
- Web frameworks (React, Vue, Django, Rails)
- Server software (nginx, Apache, IIS)
- CDNs and WAFs (Cloudflare, Akamai, AWS CloudFront)
- CMS platforms (WordPress, Drupal, Joomla)

#### Step 8: Vulnerability Scanning
**Run automated security checks**

```bash
# Vulnerability scanning with Nuclei (PLANNED integration)
# Nuclei is already installed - use it directly for now:

# Scan all alive hosts
./recon-cli recon results export example.com --format txt --alive-only -o alive.txt
nuclei -l alive.txt -t ~/nuclei-templates/ -o vulnerabilities.txt

# Scan for specific vulnerability types
nuclei -l alive.txt -t ~/nuclei-templates/cves/ -severity critical,high

# Scan for misconfigurations
nuclei -l alive.txt -t ~/nuclei-templates/misconfiguration/

# Scan for exposed panels
nuclei -l alive.txt -t ~/nuclei-templates/exposed-panels/
```

**Check For:**
- ✅ Known CVEs
- ✅ Misconfigurations
- ✅ Exposed admin panels
- ✅ Default credentials
- ✅ Information disclosure
- ✅ Missing security headers

#### Step 9: Content Discovery
**Find hidden endpoints and files**

```bash
# Content discovery (PLANNED - not yet implemented)
# Use external tools for now:

# Fast fuzzing with ffuf
ffuf -u https://example.com/FUZZ -w /path/to/wordlist.txt -o discovery.json

# Directory brute force with gobuster
gobuster dir -u https://example.com -w /path/to/wordlist.txt -o dirs.txt

# Recursive discovery with feroxbuster
feroxbuster -u https://example.com -w /path/to/wordlist.txt
```

**Discover:**
- 🔍 Admin panels, API endpoints, backup files
- 📁 `.git`, `.env`, config files, `backup.sql`
- 🧪 Development/staging endpoints
- 📝 Documentation, changelogs, READMEs

#### Step 10: Visual Reconnaissance
**Screenshot all alive hosts for quick review**

```bash
# Visual reconnaissance (PLANNED - not yet implemented)
# Use external tools for now:

# Screenshots with gowitness (recommended)
gowitness file -f alive.txt -P screenshots/

# Screenshots with aquatone
cat alive.txt | aquatone -out aquatone_results/

# Screenshots with eyewitness
eyewitness -f alive.txt -d eyewitness_results/
```

**Benefits:**
- 📸 Quickly identify interesting targets visually
- 🎯 Find login panels, admin interfaces, custom apps
- 👀 Spot unusual pages that deserve manual testing
- 📊 Generate visual reports for client deliverables

---

### Current Workflow Example

Here's a complete reconnaissance session:

```bash
# 1. Discover all subdomains
./recon-cli recon subdomain tesla.com
# Output: Found 808 unique subdomains

# 2. Find which ones are alive
./recon-cli recon verify tesla.com
# Output: 156 alive (19.3%), 652 dead (80.7%)

# 3. Get domain registration info
./recon-cli recon whois tesla.com
# Output: Registrar, nameservers, expiry date

# 4. Get DNS records for alive hosts
./recon-cli recon dns tesla.com
# Output: 142 A records, 18 unique IPs, cloud providers detected
#         Potential subdomain takeover opportunities identified

# 5. Export everything for further testing
./recon-cli recon results export tesla.com --format csv --alive-only
# Output: tesla_alive_hosts.csv with 156 entries + all DNS data

# 6. Ready to attack!
# You now have:
# - 156 alive subdomains
# - 142 with IP addresses
# - Cloud provider information
# - Potential subdomain takeovers to exploit
# - Mail infrastructure to test
```

**What You'll Have:**
- ✅ Complete subdomain inventory (808 total)
- ✅ List of alive/accessible hosts with status codes (156 alive)
- ✅ Domain registration information
- ✅ DNS records and IP mappings (142 IPs)
- ✅ Cloud provider identification (AWS, Azure, GCP, Cloudflare, etc.)
- ✅ Subdomain takeover opportunities detected
- ✅ Mail server infrastructure mapped
- ✅ Security configuration analysis (SPF, DMARC, DKIM)
- ✅ Exportable data for tools like Burp Suite, nuclei, nmap, etc.

**You're now ready for Phase 3: Active Scanning (port scanning, vulnerability scanning, etc.)**

---

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

# Dashboard Demo

## Overview

The Recontronic CLI now features an **interactive dashboard** that displays on startup, giving you immediate visibility into your reconnaissance activities.

## What You See

### 1. Header Bar
```
Recontronic CLI    Server: localhost:8080 [Offline]
Authenticated | Tools: 1/6 available
```
- Server connection status
- Authentication status
- Tool availability count

### 2. Quick Statistics
```
📊 QUICK STATISTICS
  Domains Scanned:  2
  Subdomains Found: 198
  Alive Targets:    89
  Last 24h Scans:   2
  Storage Used:     1.2 KB
```
- Overall metrics from your workspace
- Real-time calculation from results directory

### 3. Recent Activity
```
🔍 RECENT ACTIVITY
  ✓  5m ago  example.com - subdomain enum (156 found)
  ✓  3m ago  example.com - verify (89/156 alive)
  ✓  2h ago  target.org - subdomain enum (42 found)
  ✓  1d ago  testsite.com - subdomain enum (28 found)
  ✓  2d ago  hackerone.com - verify (234/456 alive)
```
- Last 5 operations with human-readable timestamps
- Status indicators: ✓ (completed), ✗ (failed), ⋯ (in progress)
- Automatically logged to `~/.recon-cli/activity.log`

### 4. System Status
```
⚙️  SYSTEM STATUS
  ✓ crt.sh          built-in
  ✗ subfinder       (not installed)
  ✗ amass           (not installed)
  ✗ assetfinder     (not installed)
  ✗ httpx           (not installed)
  ✗ nuclei          (not installed)
```
- Auto-detects which recon tools are installed
- Shows version info when available
- Helps you understand your capability

### 5. Smart Suggestions
```
💡 SUGGESTIONS
  • target.org has 1 unverified subdomains
  • example.com not scanned in 7d - consider re-scanning
  • Install subfinder for better coverage
```
- Context-aware recommendations
- Actionable next steps
- Priority-based (high priority shown first)

## Usage

### Automatic Display
```bash
# Dashboard shows automatically on startup
./recon-cli

# Then you're in interactive mode
> help
> auth login
> exit
```

### Manual Refresh
```bash
# In interactive mode, refresh anytime
> dash
> dashboard
> refresh
```

### Standalone Command
```bash
# Show dashboard and exit
./recon-cli dashboard
```

## How It Works

### Activity Logging
Every recon operation is automatically logged to `~/.recon-cli/activity.log`:

```json
{"timestamp":"2025-01-31T14:30:22Z","domain":"example.com","action":"subdomain enum","status":"completed","result":"156 found"}
{"timestamp":"2025-01-31T14:35:10Z","domain":"example.com","action":"verify","status":"completed","result":"89/156 alive"}
```

Format: JSON Lines (one JSON object per line)

### Statistics Gathering
Dashboard scans `~/.recon-cli/results/` directory:
- Counts domains (subdirectories)
- Parses JSON result files
- Calculates totals (subdomains, alive targets)
- Measures storage usage

### Tool Detection
Uses `exec.LookPath()` to detect installed tools:
- Searches system PATH
- Attempts to get version info
- Falls back gracefully if not found

### Smart Suggestions
Analyzes your workspace to suggest:

1. **Unverified Results** (Priority 1: High)
   - Scans results for subdomains without verification
   - Suggests: `recon verify <domain>`

2. **Old Scans** (Priority 3: Low)
   - Finds scans older than 7 days
   - Suggests: `recon subdomain <domain>`

3. **Missing Tools** (Priority 3: Low)
   - Checks for uninstalled tools
   - Suggests installation for better coverage

## Testing

### Create Sample Data
```bash
# Generate sample results and activity log
./scripts/create-sample-data.sh

# Now run the CLI to see populated dashboard
./recon-cli
```

This creates:
- Sample subdomain results for example.com
- Sample subdomain results for target.org
- Activity log with 5 sample entries

### Expected Output
You should see:
- 2 domains scanned
- ~198 subdomains total
- Recent activity populated
- Suggestions to verify target.org
- Tool status showing only crt.sh available

## Integration Points

### When Operations Happen
Future recon commands will automatically:
1. Log activity via `ui.LogActivity()`
2. Save results to `~/.recon-cli/results/<domain>/`
3. Dashboard picks up changes on next refresh

### Adding New Operations
To log a new operation type:

```go
import "github.com/presstronic/recontronic-cli-client/pkg/ui"

// After operation completes
err := ui.LogActivity(ui.ActivityEntry{
    Timestamp: time.Now(),
    Domain:    "example.com",
    Action:    "subdomain enum",
    Status:    "completed",
    Result:    "156 found",
})
```

## Benefits

1. **Immediate Context** - See what you've been working on
2. **Quick Overview** - Understand coverage at a glance
3. **Actionable Insights** - Know what to do next
4. **Tool Awareness** - Understand your capabilities
5. **Progress Tracking** - Watch your findings grow

## Future Enhancements

Planned for Phase 2:
- [ ] Rich TUI with Bubble Tea (colors, interactive)
- [ ] Live progress bars during scans
- [ ] ASCII charts for trend visualization
- [ ] Watch mode (auto-refresh every 5s)
- [ ] Latest findings panel (new alive subdomains)

---

**The dashboard transforms the CLI from a simple command runner into an intelligent recon workspace!**

#!/bin/bash
# Create sample data for dashboard testing

CONFIG_DIR="$HOME/.recon-cli"
RESULTS_DIR="$CONFIG_DIR/results"

# Create directories
mkdir -p "$RESULTS_DIR/example.com"
mkdir -p "$RESULTS_DIR/target.org"

# Create sample subdomain results for example.com
cat > "$RESULTS_DIR/example.com/subdomains_$(date +%Y%m%d_%H%M%S).json" <<'EOF'
{
  "domain": "example.com",
  "timestamp": "2025-01-31T14:30:22Z",
  "sources_used": ["crtsh", "subfinder"],
  "total_unique": 156,
  "total_alive": 89,
  "subdomains": [
    {
      "name": "api.example.com",
      "discovered_by": ["crtsh", "subfinder"],
      "first_seen": "2025-01-31T14:30:22Z",
      "verified": {
        "status": "alive"
      }
    },
    {
      "name": "app.example.com",
      "discovered_by": ["crtsh", "subfinder"],
      "first_seen": "2025-01-31T14:30:22Z",
      "verified": {
        "status": "alive"
      }
    },
    {
      "name": "staging.example.com",
      "discovered_by": ["subfinder"],
      "first_seen": "2025-01-31T14:30:45Z",
      "verified": {
        "status": "alive"
      }
    },
    {
      "name": "old.example.com",
      "discovered_by": ["crtsh"],
      "first_seen": "2025-01-31T14:30:22Z",
      "verified": {
        "status": "dead"
      }
    }
  ]
}
EOF

# Create sample subdomain results for target.org
cat > "$RESULTS_DIR/target.org/subdomains_$(date -v-2d +%Y%m%d_%H%M%S).json" <<'EOF'
{
  "domain": "target.org",
  "timestamp": "2025-01-29T09:15:33Z",
  "sources_used": ["crtsh"],
  "total_unique": 42,
  "subdomains": [
    {
      "name": "www.target.org",
      "discovered_by": ["crtsh"],
      "first_seen": "2025-01-29T09:15:33Z"
    }
  ]
}
EOF

# Create sample activity log
cat > "$CONFIG_DIR/activity.log" <<EOF
{"timestamp":"$(date -u -v-5M +%Y-%m-%dT%H:%M:%SZ)","domain":"example.com","action":"subdomain enum","status":"completed","result":"156 found"}
{"timestamp":"$(date -u -v-3M +%Y-%m-%dT%H:%M:%SZ)","domain":"example.com","action":"verify","status":"completed","result":"89/156 alive"}
{"timestamp":"$(date -u -v-2H +%Y-%m-%dT%H:%M:%SZ)","domain":"target.org","action":"subdomain enum","status":"completed","result":"42 found"}
{"timestamp":"$(date -u -v-1d +%Y-%m-%dT%H:%M:%SZ)","domain":"testsite.com","action":"subdomain enum","status":"completed","result":"28 found"}
{"timestamp":"$(date -u -v-2d +%Y-%m-%dT%H:%M:%SZ)","domain":"hackerone.com","action":"verify","status":"completed","result":"234/456 alive"}
EOF

echo "Sample data created in $CONFIG_DIR"
echo "You can now test the dashboard!"

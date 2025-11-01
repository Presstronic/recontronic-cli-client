#!/bin/bash

# Script to close completed issues in GitHub

REPO="Presstronic/recontronic-cli-client"

echo "Closing completed issues in $REPO..."
echo ""

close_issue() {
    local issue_id="$1"
    local completion_note="$2"

    # Find the GitHub issue number for this RECON ID
    issue_number=$(gh issue list --repo "$REPO" --state all --limit 100 --search "$issue_id" --json number,title --jq ".[] | select(.title | startswith(\"$issue_id:\")) | .number")

    if [ -n "$issue_number" ]; then
        echo "Closing #$issue_number ($issue_id)..."

        # Close the issue with a completion comment
        gh issue close "$issue_number" --repo "$REPO" --comment "✅ **Completed**

$completion_note

Completed as part of the initial authentication system implementation. See commit history for implementation details."

        echo "  ✓ Closed #$issue_number"
    else
        echo "  ✗ Could not find issue for $issue_id"
    fi

    # Small delay to avoid rate limiting
    sleep 0.5
}

# Close each completed issue
close_issue "RECON-001" "Initial Go module setup, project structure, and go.mod/go.sum created"
close_issue "RECON-002" "Configuration management with Viper implemented, supports config file and env vars"
close_issue "RECON-003" "Root Cobra command implemented with global flags (--config, --debug, --output)"
close_issue "RECON-007" "REST API client fully implemented with all auth endpoints, proper error handling, and Bearer token authentication"
close_issue "RECON-008" "All data models defined: User, APIKey, Program, Scan, Anomaly with JSON tags and validation"
close_issue "RECON-051" "User registration command implemented with interactive prompts and validation"
close_issue "RECON-052" "User login command implemented, returns and saves API key automatically"
close_issue "RECON-053" "Whoami command implemented to display current authenticated user information"
close_issue "RECON-054" "Create API keys command implemented with optional name and expiration"
close_issue "RECON-055" "List API keys command implemented with table output showing all key details"
close_issue "RECON-056" "Revoke API key command implemented with confirmation prompt"
close_issue "RECON-058" "Secure password input implemented using golang.org/x/term with no echo"
close_issue "RECON-059" "API key storage implemented with secure 0600 permissions in ~/.recon-cli/config.yaml"
close_issue "RECON-060" "Authentication middleware implemented in REST client with Bearer token support"

echo ""
echo "=========================================="
echo "Completed issue closure process!"
echo "=========================================="

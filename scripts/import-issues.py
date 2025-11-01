#!/usr/bin/env python3
"""
Import issues from mvp-issues.csv to GitHub Issues.
Requires: gh CLI tool to be installed and authenticated.
"""

import csv
import subprocess
import sys
import json

def parse_labels(label_str):
    """Parse semicolon-separated labels."""
    if not label_str or label_str.strip() == "":
        return []
    return [l.strip() for l in label_str.split(';') if l.strip()]

def create_issue(issue_id, issue_type, priority, title, story, acceptance_criteria,
                 epic, points, dependencies, labels):
    """Create a GitHub issue using gh CLI."""

    # Build issue body
    body = f"**Type:** {issue_type}\n"
    body += f"**Priority:** {priority}\n"
    body += f"**Epic:** {epic}\n"
    body += f"**Story Points:** {points}\n"
    if dependencies and dependencies.strip():
        body += f"**Dependencies:** {dependencies}\n"
    body += f"\n## Story\n\n{story}\n"
    body += f"\n## Acceptance Criteria\n\n{acceptance_criteria}\n"

    # Parse labels
    label_list = parse_labels(labels)

    # Add priority as a label if it doesn't exist
    priority_lower = priority.lower()
    if priority_lower == "critical" and "critical" not in label_list:
        label_list.append("critical")

    # Build gh command
    cmd = [
        "gh", "issue", "create",
        "--repo", "Presstronic/recontronic-cli-client",
        "--title", f"{issue_id}: {title}",
        "--body", body
    ]

    # Add labels
    if label_list:
        for label in label_list:
            cmd.extend(["--label", label])

    try:
        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
        issue_url = result.stdout.strip()
        print(f"✓ Created {issue_id}: {title}")
        print(f"  URL: {issue_url}")
        return issue_url
    except subprocess.CalledProcessError as e:
        print(f"✗ Failed to create {issue_id}: {title}")
        print(f"  Error: {e.stderr}")
        return None

def main():
    csv_file = "mvp-issues.csv"

    print(f"Importing issues from {csv_file}...")
    print(f"Repository: Presstronic/recontronic-cli-client\n")

    created_count = 0
    failed_count = 0

    try:
        with open(csv_file, 'r', encoding='utf-8') as f:
            reader = csv.DictReader(f)

            for row in reader:
                issue_id = row['Issue ID']
                issue_type = row['Type']
                priority = row['Priority']
                title = row['Title']
                story = row['Story']
                acceptance = row['Acceptance Criteria']
                epic = row['Epic']
                points = row['Estimated Points']
                deps = row['Dependencies']
                labels = row['Labels']

                result = create_issue(issue_id, issue_type, priority, title,
                                    story, acceptance, epic, points, deps, labels)

                if result:
                    created_count += 1
                else:
                    failed_count += 1

                # Small delay to avoid rate limiting
                import time
                time.sleep(0.5)

    except FileNotFoundError:
        print(f"Error: {csv_file} not found!")
        print("Please run this script from the project root directory.")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

    print(f"\n{'='*60}")
    print(f"Import complete!")
    print(f"  Created: {created_count} issues")
    print(f"  Failed:  {failed_count} issues")
    print(f"{'='*60}")

if __name__ == "__main__":
    main()

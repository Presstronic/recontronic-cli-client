package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/config"
)

// ActivityEntry represents a single activity log entry
type ActivityEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Domain    string    `json:"domain"`
	Action    string    `json:"action"` // "subdomain", "verify", "dns", "whois"
	Status    string    `json:"status"` // "completed", "failed", "in_progress"
	Result    string    `json:"result"` // "156 found", "89/156 alive", etc.
	Error     string    `json:"error,omitempty"`
}

// GetActivityLogPath returns the path to the activity log file
func GetActivityLogPath() (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "activity.log"), nil
}

// LogActivity appends an activity entry to the log
func LogActivity(entry ActivityEntry) error {
	logPath, err := GetActivityLogPath()
	if err != nil {
		return fmt.Errorf("failed to get activity log path: %w", err)
	}

	// Ensure config directory exists
	if err := config.EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to ensure config directory: %w", err)
	}

	// Open file in append mode, create if doesn't exist
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open activity log: %w", err)
	}
	defer file.Close()

	// Marshal entry to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal activity entry: %w", err)
	}

	// Write JSON line
	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write activity entry: %w", err)
	}

	return nil
}

// GetRecentActivity retrieves the last N activity entries
func GetRecentActivity(limit int) ([]ActivityEntry, error) {
	logPath, err := GetActivityLogPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get activity log path: %w", err)
	}

	// Check if log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return []ActivityEntry{}, nil
	}

	// Read entire file
	data, err := os.ReadFile(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read activity log: %w", err)
	}

	// Parse JSON lines
	var entries []ActivityEntry
	lines := string(data)

	// Split by newlines and parse each line
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] == '\n' {
			lineEnd := i + 1
			lineStart := i + 1

			// Find start of line
			for j := i - 1; j >= 0; j-- {
				if lines[j] == '\n' {
					lineStart = j + 1
					break
				}
				if j == 0 {
					lineStart = 0
					break
				}
			}

			line := lines[lineStart:lineEnd]
			if len(line) > 0 && line[0] == '{' {
				var entry ActivityEntry
				if err := json.Unmarshal([]byte(line), &entry); err == nil {
					entries = append(entries, entry)
					if len(entries) >= limit {
						break
					}
				}
			}

			i = lineStart
		}
	}

	return entries, nil
}

// FormatTimeAgo formats a timestamp as "5m ago", "2h ago", etc.
func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "just now"
	}
	if duration < time.Hour {
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", mins)
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1h ago"
		}
		return fmt.Sprintf("%dh ago", hours)
	}
	if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	}

	return t.Format("Jan 2")
}

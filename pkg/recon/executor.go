package recon

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// ExecutionResult represents the result of a command execution
type ExecutionResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
}

// Execute runs a command safely with timeout and context
func Execute(ctx context.Context, name string, args ...string) (*ExecutionResult, error) {
	startTime := time.Now()

	cmd := exec.CommandContext(ctx, name, args...)

	// Capture stdout and stderr
	stdout, err := cmd.Output()
	duration := time.Since(startTime)

	result := &ExecutionResult{
		Stdout:   string(stdout),
		Duration: duration,
	}

	if err != nil {
		// Check if it's an ExitError (command ran but failed)
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.Stderr = string(exitErr.Stderr)
			result.ExitCode = exitErr.ExitCode()
			return result, fmt.Errorf("command failed with exit code %d: %s", exitErr.ExitCode(), result.Stderr)
		}
		// Command couldn't be started
		return nil, fmt.Errorf("failed to execute command: %w", err)
	}

	result.ExitCode = 0
	return result, nil
}

// ExecuteWithTimeout is a convenience wrapper that adds a timeout
func ExecuteWithTimeout(name string, timeout time.Duration, args ...string) (*ExecutionResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := Execute(ctx, name, args...)
	if err != nil {
		// Check if it was a timeout
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("command timed out after %s", timeout)
		}
		return result, err
	}

	return result, nil
}

// IsToolAvailable checks if a command-line tool is available
func IsToolAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

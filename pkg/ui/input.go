package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// ReadPassword reads a password from stdin without echoing
func ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	// Get file descriptor for stdin
	fd := int(os.Stdin.Fd())

	// Check if stdin is a terminal
	if !term.IsTerminal(fd) {
		// Not a terminal, read normally (for testing or piped input)
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			return scanner.Text(), nil
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", fmt.Errorf("no input received")
	}

	// Read password without echo
	password, err := term.ReadPassword(fd)
	fmt.Println() // Print newline after password input

	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	return string(password), nil
}

// ReadPasswordWithConfirm reads a password and asks for confirmation
func ReadPasswordWithConfirm(prompt, confirmPrompt string) (string, error) {
	password, err := ReadPassword(prompt)
	if err != nil {
		return "", err
	}

	if len(password) == 0 {
		return "", fmt.Errorf("password cannot be empty")
	}

	confirmation, err := ReadPassword(confirmPrompt)
	if err != nil {
		return "", err
	}

	if password != confirmation {
		return "", fmt.Errorf("passwords do not match")
	}

	return password, nil
}

// ReadInput reads a line of input from stdin
func ReadInput(prompt string) (string, error) {
	fmt.Print(prompt)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	return "", fmt.Errorf("no input received")
}

// ReadInputWithDefault reads input with a default value
func ReadInputWithDefault(prompt, defaultValue string) (string, error) {
	if defaultValue != "" {
		prompt = fmt.Sprintf("%s [%s]: ", prompt, defaultValue)
	} else {
		prompt = fmt.Sprintf("%s: ", prompt)
	}

	input, err := ReadInput(prompt)
	if err != nil {
		return "", err
	}

	if input == "" && defaultValue != "" {
		return defaultValue, nil
	}

	return input, nil
}

// Confirm asks for yes/no confirmation
func Confirm(prompt string) (bool, error) {
	response, err := ReadInput(fmt.Sprintf("%s [y/N]: ", prompt))
	if err != nil {
		return false, err
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes", nil
}

// ValidateEmail performs basic email validation
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidateUsername validates username according to server rules
func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters")
	}
	if len(username) > 50 {
		return fmt.Errorf("username must be at most 50 characters")
	}
	// Check alphanumeric
	for _, ch := range username {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
			return fmt.Errorf("username must contain only letters and numbers")
		}
	}
	return nil
}

// ValidatePassword validates password according to server rules
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 72 {
		return fmt.Errorf("password must be at most 72 characters")
	}
	return nil
}

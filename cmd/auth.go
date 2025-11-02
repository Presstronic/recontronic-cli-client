package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/client"
	"github.com/presstronic/recontronic-cli-client/pkg/config"
	"github.com/presstronic/recontronic-cli-client/pkg/ui"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication and API key management",
	Long: `Manage user authentication and API keys for the Recontronic platform.

Commands include user registration, login, viewing current user info,
and creating, listing, and revoking API keys.`,
}

var authRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user account",
	Long: `Register a new user account with the Recontronic platform.

You will be prompted for username, email, and password. After successful
registration, use 'recon-cli auth login' to authenticate and receive an API key.`,
	RunE: runAuthRegister,
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login and receive an API key",
	Long: `Authenticate with the Recontronic platform and receive an API key.

The API key will be saved to your configuration file (~/.recon-cli/config.yaml)
and used automatically for all subsequent commands.`,
	RunE: runAuthLogin,
}

var authWhoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Display current authenticated user information",
	Long: `Display information about the currently authenticated user.

Requires a valid API key in your configuration.`,
	RunE: runAuthWhoami,
}

var authKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage API keys",
	Long: `Create, list, and revoke API keys for your account.

API keys are used to authenticate with the Recontronic platform. You can
create multiple keys for different purposes (e.g., development, production,
CI/CD) and revoke them individually.`,
}

var authKeysCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new API key",
	Long: `Create a new API key for your account.

You can optionally specify a name and expiration date for the key.`,
	RunE: runAuthKeysCreate,
}

var authKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API keys",
	Long: `List all API keys associated with your account.

Shows key ID, name, prefix, expiration, last used time, and status.`,
	RunE: runAuthKeysList,
}

var authKeysRevokeCmd = &cobra.Command{
	Use:   "revoke <key-id>",
	Short: "Revoke an API key",
	Long: `Revoke an API key by its ID.

The key will be immediately deactivated and cannot be used for authentication.`,
	Args: cobra.ExactArgs(1),
	RunE: runAuthKeysRevoke,
}

var (
	keyName      string
	keyExpiresIn string
	forceRevoke  bool
)

func init() {
	authCmd.AddCommand(authRegisterCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authWhoamiCmd)
	authCmd.AddCommand(authKeysCmd)

	authKeysCmd.AddCommand(authKeysCreateCmd)
	authKeysCmd.AddCommand(authKeysListCmd)
	authKeysCmd.AddCommand(authKeysRevokeCmd)

	authKeysCreateCmd.Flags().StringVarP(&keyName, "name", "n", "", "Name for the API key")
	authKeysCreateCmd.Flags().StringVar(&keyExpiresIn, "expires-in", "", "Expiration duration (e.g., 90d, 1y)")

	authKeysRevokeCmd.Flags().BoolVarP(&forceRevoke, "force", "f", false, "Skip confirmation prompt")
}

func runAuthRegister(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	fmt.Println("Register a new Recontronic account")

	username, err := ui.ReadInput("Username: ")
	if err != nil {
		return fmt.Errorf("failed to read username: %w", err)
	}

	if err := ui.ValidateUsername(username); err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	email, err := ui.ReadInput("Email: ")
	if err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}

	if err := ui.ValidateEmail(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	password, err := ui.ReadPasswordWithConfirm("Password: ", "Confirm password: ")
	if err != nil {
		return fmt.Errorf("password error: %w", err)
	}

	if err := ui.ValidatePassword(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	restClient := client.NewRestClient(cfg.Server, "", cfg.Timeout)
	if debug {
		restClient.SetDebug(true)
	}

	user, err := restClient.Register(ctx, username, email, password)
	if err != nil {
		if client.IsValidationError(err) {
			return fmt.Errorf("registration failed - please check your inputs: %w", err)
		}
		return fmt.Errorf("registration failed: %w", err)
	}

	fmt.Println("\n✓ Registration successful!")
	fmt.Printf("Account created for: %s\n", user.Username)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Println("\nNext step: Login to get your API key")
	fmt.Println("  $ recon-cli auth login")

	return nil
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	fmt.Println("Login to Recontronic")

	username, err := ui.ReadInput("Username: ")
	if err != nil {
		return fmt.Errorf("failed to read username: %w", err)
	}

	password, err := ui.ReadPassword("Password: ")
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}

	restClient := client.NewRestClient(cfg.Server, "", cfg.Timeout)
	if debug {
		restClient.SetDebug(true)
	}

	loginResp, err := restClient.Login(ctx, username, password)
	if err != nil {
		if client.IsAuthError(err) {
			return fmt.Errorf("login failed: invalid username or password")
		}
		return fmt.Errorf("login failed: %w", err)
	}

	if err := config.SaveAPIKey(loginResp.APIKey); err != nil {
		fmt.Println("\n✓ Login successful!")
		fmt.Printf("\nYour API key: %s\n", loginResp.APIKey)
		fmt.Println("\n⚠️  WARNING: Failed to save API key to config file")
		fmt.Printf("Error: %v\n", err)
		fmt.Println("\nPlease save it manually:")
		fmt.Printf("  $ recon-cli config set api-key %s\n", loginResp.APIKey)
		return nil
	}

	configPath, _ := config.GetConfigPath()

	fmt.Println("\n✓ Login successful!")
	fmt.Printf("\nYour API key: %s\n", loginResp.APIKey)
	fmt.Println("\n⚠️  IMPORTANT: Save this key securely!")
	fmt.Printf("   It has been saved to: %s\n", configPath)
	fmt.Println("   This key will not be shown again.")
	fmt.Println("\nYou're now authenticated and ready to use the CLI.")

	return nil
}

func runAuthWhoami(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if cfg.APIKey == "" {
		return fmt.Errorf("not authenticated: please run 'recon-cli auth login' first")
	}

	restClient := client.NewRestClient(cfg.Server, cfg.APIKey, cfg.Timeout)
	if debug {
		restClient.SetDebug(true)
	}

	user, err := restClient.GetCurrentUser(ctx)
	if err != nil {
		if client.IsAuthError(err) {
			return fmt.Errorf("authentication failed: your API key may be invalid or expired\nPlease run 'recon-cli auth login' to get a new key")
		}
		return fmt.Errorf("failed to get user info: %w", err)
	}

	keyPrefix := "Not available"
	if len(cfg.APIKey) >= 8 {
		keyPrefix = cfg.APIKey[:8] + "..."
	}

	fmt.Printf("Username:     %s\n", user.Username)
	fmt.Printf("Email:        %s\n", user.Email)
	fmt.Printf("Account ID:   %d\n", user.ID)
	fmt.Printf("Status:       %s\n", formatStatus(user.IsActive))
	fmt.Printf("Created:      %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("API Key:      %s\n", keyPrefix)

	return nil
}

func runAuthKeysCreate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if cfg.APIKey == "" {
		return fmt.Errorf("not authenticated: please run 'recon-cli auth login' first")
	}

	var expiresAt *time.Time
	if keyExpiresIn != "" {
		duration, err := parseDuration(keyExpiresIn)
		if err != nil {
			return fmt.Errorf("invalid expiration duration: %w", err)
		}
		expiry := time.Now().Add(duration)
		expiresAt = &expiry
	}

	restClient := client.NewRestClient(cfg.Server, cfg.APIKey, cfg.Timeout)
	if debug {
		restClient.SetDebug(true)
	}

	apiKey, err := restClient.CreateAPIKey(ctx, keyName, expiresAt)
	if err != nil {
		if client.IsAuthError(err) {
			return fmt.Errorf("authentication failed: please run 'recon-cli auth login' first")
		}
		return fmt.Errorf("failed to create API key: %w", err)
	}

	fmt.Println("✓ New API key created!")
	fmt.Printf("\nAPI Key: %s\n", apiKey.PlainKey)
	if apiKey.Name != "" {
		fmt.Printf("Name:    %s\n", apiKey.Name)
	}
	fmt.Printf("ID:      %d\n", apiKey.ID)
	if apiKey.ExpiresAt != nil {
		fmt.Printf("Expires: %s\n", apiKey.ExpiresAt.Format("2006-01-02"))
	} else {
		fmt.Println("Expires: Never")
	}
	fmt.Println("\n⚠️  Save this key! It won't be shown again.")

	return nil
}

func runAuthKeysList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if cfg.APIKey == "" {
		return fmt.Errorf("not authenticated: please run 'recon-cli auth login' first")
	}

	restClient := client.NewRestClient(cfg.Server, cfg.APIKey, cfg.Timeout)
	if debug {
		restClient.SetDebug(true)
	}

	response, err := restClient.ListAPIKeys(ctx)
	if err != nil {
		if client.IsAuthError(err) {
			return fmt.Errorf("authentication failed: please run 'recon-cli auth login' first")
		}
		return fmt.Errorf("failed to list API keys: %w", err)
	}

	if len(response.APIKeys) == 0 {
		fmt.Println("No API keys found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tPREFIX\tLAST USED\tEXPIRES\tSTATUS")
	fmt.Fprintln(w, "──\t────\t──────\t─────────\t───────\t──────")

	for _, key := range response.APIKeys {
		name := key.Name
		if name == "" {
			name = "-"
		}

		lastUsed := "Never"
		if key.LastUsedAt != nil {
			lastUsed = formatTimeAgo(*key.LastUsedAt)
		}

		expires := "Never"
		if key.ExpiresAt != nil {
			expires = formatExpiresAt(*key.ExpiresAt)
		}

		status := formatStatus(key.IsActive)

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			key.ID, name, key.KeyPrefix, lastUsed, expires, status)
	}

	w.Flush()
	fmt.Printf("\nTotal: %d API key(s)\n", response.Total)

	return nil
}

func runAuthKeysRevoke(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if cfg.APIKey == "" {
		return fmt.Errorf("not authenticated: please run 'recon-cli auth login' first")
	}

	keyID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid key ID: %w", err)
	}

	if !forceRevoke {
		confirmed, err := ui.Confirm(fmt.Sprintf("Are you sure you want to revoke API key ID %d?", keyID))
		if err != nil {
			return fmt.Errorf("confirmation failed: %w", err)
		}
		if !confirmed {
			fmt.Println("Revocation cancelled.")
			return nil
		}
	}

	restClient := client.NewRestClient(cfg.Server, cfg.APIKey, cfg.Timeout)
	if debug {
		restClient.SetDebug(true)
	}

	err = restClient.RevokeAPIKey(ctx, keyID)
	if err != nil {
		if client.IsAuthError(err) {
			return fmt.Errorf("authentication failed: please run 'recon-cli auth login' first")
		}
		if client.IsNotFoundError(err) {
			return fmt.Errorf("API key not found (ID: %d)", keyID)
		}
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	fmt.Printf("✓ API key %d revoked successfully\n", keyID)

	return nil
}

func formatStatus(isActive bool) string {
	if isActive {
		return "Active"
	}
	return "Inactive"
}

func formatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "Just now"
	}
	if duration < time.Hour {
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}

	return t.Format("2006-01-02")
}

func formatExpiresAt(t time.Time) string {
	duration := time.Until(t)

	if duration < 0 {
		return "Expired"
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours < 1 {
			return "Expires soon"
		}
		if hours == 1 {
			return "in 1 hour"
		}
		return fmt.Sprintf("in %d hours", hours)
	}
	if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "in 1 day"
		}
		return fmt.Sprintf("in %d days", days)
	}

	return t.Format("2006-01-02")
}

func parseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration format")
	}

	unit := s[len(s)-1:]
	valueStr := s[:len(s)-1]

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %w", err)
	}

	switch unit {
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	case "w":
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	case "m":
		return time.Duration(value) * 30 * 24 * time.Hour, nil
	case "y":
		return time.Duration(value) * 365 * 24 * time.Hour, nil
	default:
		return time.ParseDuration(s)
	}
}

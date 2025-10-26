package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config represents the CLI configuration
type Config struct {
	Server       string        `mapstructure:"server"`
	GRPCServer   string        `mapstructure:"grpc_server"`
	APIKey       string        `mapstructure:"api_key"`
	Timeout      time.Duration `mapstructure:"timeout"`
	OutputFormat string        `mapstructure:"output_format"`
	LogLevel     string        `mapstructure:"log_level"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Server:       "http://localhost:8080",
		GRPCServer:   "localhost:9090",
		Timeout:      30 * time.Second,
		OutputFormat: "table",
		LogLevel:     "info",
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(home, ".recon-cli")
	configFile := filepath.Join(configDir, "config.yaml")

	return configFile, nil
}

// GetConfigDir returns the path to the config directory
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, ".recon-cli"), nil
}

// EnsureConfigDir creates the config directory with secure permissions
func EnsureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create directory with 0700 permissions (owner only)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Verify permissions
	info, err := os.Stat(configDir)
	if err != nil {
		return fmt.Errorf("failed to stat config directory: %w", err)
	}

	// Check permissions (on Unix-like systems)
	if info.Mode().Perm() != 0700 {
		// Try to fix permissions
		if err := os.Chmod(configDir, 0700); err != nil {
			return fmt.Errorf("config directory has insecure permissions and could not be fixed: %w", err)
		}
	}

	return nil
}

// SecureConfigFile sets secure permissions on the config file
func SecureConfigFile(path string) error {
	// Set file permissions to 0600 (owner read/write only)
	if err := os.Chmod(path, 0600); err != nil {
		return fmt.Errorf("failed to set secure permissions on config file: %w", err)
	}

	return nil
}

// Load reads the configuration from file and environment
func Load(cfgFile string) (*Config, error) {
	// Set defaults
	viper.SetDefault("server", "http://localhost:8080")
	viper.SetDefault("grpc_server", "localhost:9090")
	viper.SetDefault("timeout", "30s")
	viper.SetDefault("output_format", "table")
	viper.SetDefault("log_level", "info")

	// Environment variable support with RECON_ prefix
	viper.SetEnvPrefix("RECON")
	viper.AutomaticEnv()

	// If a config file is specified, use it
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Use default config location
		configDir, err := GetConfigDir()
		if err != nil {
			return nil, err
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Read config file (it's okay if it doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
		// Config file not found, use defaults
	}

	// Parse into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Parse timeout string to duration if needed
	if viper.IsSet("timeout") {
		timeoutStr := viper.GetString("timeout")
		duration, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout format: %w", err)
		}
		cfg.Timeout = duration
	}

	return &cfg, nil
}

// Save writes the current configuration to file
func Save(cfg *Config) error {
	// Ensure config directory exists
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Set values in viper
	viper.Set("server", cfg.Server)
	viper.Set("grpc_server", cfg.GRPCServer)
	viper.Set("api_key", cfg.APIKey)
	viper.Set("timeout", cfg.Timeout.String())
	viper.Set("output_format", cfg.OutputFormat)
	viper.Set("log_level", cfg.LogLevel)

	// Write config file
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// Set secure permissions
	if err := SecureConfigFile(configPath); err != nil {
		return err
	}

	return nil
}

// Set updates a single configuration value
func Set(key, value string) error {
	// Load current config
	cfg, err := Load("")
	if err != nil {
		// If config doesn't exist, start with defaults
		cfg = DefaultConfig()
	}

	// Update the specified key
	switch key {
	case "server":
		cfg.Server = value
	case "grpc-server", "grpc_server":
		cfg.GRPCServer = value
	case "api-key", "api_key":
		cfg.APIKey = value
	case "timeout":
		duration, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("invalid timeout format (use: 30s, 1m, etc.): %w", err)
		}
		cfg.Timeout = duration
	case "output-format", "output_format":
		if value != "table" && value != "json" && value != "yaml" {
			return fmt.Errorf("invalid output format (must be: table, json, or yaml)")
		}
		cfg.OutputFormat = value
	case "log-level", "log_level":
		if value != "debug" && value != "info" && value != "warn" && value != "error" {
			return fmt.Errorf("invalid log level (must be: debug, info, warn, or error)")
		}
		cfg.LogLevel = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	// Save updated config
	return Save(cfg)
}

// Get retrieves a single configuration value
func Get(key string) (string, error) {
	cfg, err := Load("")
	if err != nil {
		return "", err
	}

	switch key {
	case "server":
		return cfg.Server, nil
	case "grpc-server", "grpc_server":
		return cfg.GRPCServer, nil
	case "api-key", "api_key":
		return cfg.APIKey, nil
	case "timeout":
		return cfg.Timeout.String(), nil
	case "output-format", "output_format":
		return cfg.OutputFormat, nil
	case "log-level", "log_level":
		return cfg.LogLevel, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}

// SaveAPIKey saves only the API key to config
func SaveAPIKey(apiKey string) error {
	cfg, err := Load("")
	if err != nil {
		// If config doesn't exist, start with defaults
		cfg = DefaultConfig()
	}

	cfg.APIKey = apiKey
	return Save(cfg)
}

// ValidateAPIKey checks if an API key has the correct format
func ValidateAPIKey(apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API key cannot be empty")
	}
	if len(apiKey) < 20 {
		return fmt.Errorf("API key is too short (minimum 20 characters)")
	}
	if len(apiKey[:4]) < 4 || apiKey[:4] != "rct_" {
		return fmt.Errorf("API key must start with 'rct_'")
	}
	return nil
}

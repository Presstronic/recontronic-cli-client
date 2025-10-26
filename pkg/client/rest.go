package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/presstronic/recontronic-cli-client/pkg/models"
)

// RestClient handles HTTP communication with the Recontronic API
type RestClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	debug      bool
}

// NewRestClient creates a new REST API client
func NewRestClient(baseURL, apiKey string, timeout time.Duration) *RestClient {
	return &RestClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		debug: false,
	}
}

// SetDebug enables or disables debug logging
func (c *RestClient) SetDebug(debug bool) {
	c.debug = debug
}

// SetAPIKey updates the API key for authenticated requests
func (c *RestClient) SetAPIKey(apiKey string) {
	c.apiKey = apiKey
}

// doRequest performs an HTTP request with proper error handling
func (c *RestClient) doRequest(ctx context.Context, method, path string, body interface{}, response interface{}, authenticated bool) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)

		if c.debug {
			fmt.Printf("→ Request Body: %s\n", string(jsonData))
		}
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "recontronic-cli/1.0.0")

	// Add authentication header if required and API key is available
	if authenticated && c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
		if c.debug {
			// Sanitize API key in debug output
			sanitized := c.apiKey
			if len(sanitized) > 12 {
				sanitized = sanitized[:8] + "..." + sanitized[len(sanitized)-4:]
			}
			fmt.Printf("→ Authorization: Bearer %s\n", sanitized)
		}
	}

	if c.debug {
		fmt.Printf("→ %s %s\n", method, url)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if c.debug {
		fmt.Printf("← %d %s\n", resp.StatusCode, resp.Status)
		fmt.Printf("← Response Body: %s\n", string(respBody))
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		var errResp models.ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return &APIError{
				StatusCode: resp.StatusCode,
				Message:    errResp.Error,
			}
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
		}
	}

	// Parse success response
	if response != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, response); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Register creates a new user account
func (c *RestClient) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	req := models.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	var user models.User
	err := c.doRequest(ctx, "POST", "/api/v1/auth/register", req, &user, false)
	if err != nil {
		return nil, fmt.Errorf("registration failed: %w", err)
	}

	return &user, nil
}

// Login authenticates a user and returns an API key
func (c *RestClient) Login(ctx context.Context, username, password string) (*models.LoginResponse, error) {
	req := models.LoginRequest{
		Username: username,
		Password: password,
	}

	var loginResp models.LoginResponse
	err := c.doRequest(ctx, "POST", "/api/v1/auth/login", req, &loginResp, false)
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	return &loginResp, nil
}

// GetCurrentUser retrieves the currently authenticated user
func (c *RestClient) GetCurrentUser(ctx context.Context) (*models.User, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("authentication required: please run 'recon-cli auth login' first")
	}

	var user models.User
	err := c.doRequest(ctx, "GET", "/api/v1/auth/me", nil, &user, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	return &user, nil
}

// CreateAPIKey generates a new API key
func (c *RestClient) CreateAPIKey(ctx context.Context, name string, expiresAt *time.Time) (*models.APIKey, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("authentication required: please run 'recon-cli auth login' first")
	}

	req := models.CreateAPIKeyRequest{
		Name:      name,
		ExpiresAt: expiresAt,
	}

	var apiKey models.APIKey
	err := c.doRequest(ctx, "POST", "/api/v1/auth/keys", req, &apiKey, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}

	return &apiKey, nil
}

// ListAPIKeys retrieves all API keys for the current user
func (c *RestClient) ListAPIKeys(ctx context.Context) (*models.APIKeyListResponse, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("authentication required: please run 'recon-cli auth login' first")
	}

	var response models.APIKeyListResponse
	err := c.doRequest(ctx, "GET", "/api/v1/auth/keys", nil, &response, true)
	if err != nil {
		return nil, fmt.Errorf("failed to list API keys: %w", err)
	}

	return &response, nil
}

// RevokeAPIKey deletes/revokes an API key by ID
func (c *RestClient) RevokeAPIKey(ctx context.Context, keyID int64) error {
	if c.apiKey == "" {
		return fmt.Errorf("authentication required: please run 'recon-cli auth login' first")
	}

	path := fmt.Sprintf("/api/v1/auth/keys/%d", keyID)
	var response struct {
		Message string `json:"message"`
	}

	err := c.doRequest(ctx, "DELETE", path, nil, &response, true)
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	return nil
}

// APIError represents an error returned from the API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// IsAuthError returns true if the error is an authentication error (401)
func IsAuthError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsNotFoundError returns true if the error is a not found error (404)
func IsNotFoundError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsValidationError returns true if the error is a validation error (400)
func IsValidationError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusBadRequest
	}
	return false
}

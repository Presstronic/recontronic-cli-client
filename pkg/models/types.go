package models

import "time"

// User represents a user account in the system
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// APIKey represents an API key for authentication
type APIKey struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Name       string     `json:"name,omitempty"`
	KeyPrefix  string     `json:"key_prefix"`
	PlainKey   string     `json:"plain_key,omitempty"` // Only returned during creation
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
}

// RegisterRequest is the payload for user registration
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the payload for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse contains the login result with API key
type LoginResponse struct {
	User    User   `json:"user"`
	APIKey  string `json:"api_key"`
	KeyID   int64  `json:"key_id"`
	Message string `json:"message"`
}

// APIKeyListResponse contains a list of API keys
type APIKeyListResponse struct {
	APIKeys []APIKey `json:"api_keys"`
	Total   int      `json:"total"`
}

// CreateAPIKeyRequest is the payload for creating a new API key
type CreateAPIKeyRequest struct {
	Name      string     `json:"name,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// Program represents a bug bounty program (future use)
type Program struct {
	ID            int64                  `json:"id"`
	Name          string                 `json:"name"`
	Platform      string                 `json:"platform"`
	Scope         []string               `json:"scope"`
	ScanFrequency string                 `json:"scan_frequency"`
	CreatedAt     time.Time              `json:"created_at"`
	LastScannedAt *time.Time             `json:"last_scanned_at,omitempty"`
	IsActive      bool                   `json:"is_active"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// Scan represents a reconnaissance scan (future use)
type Scan struct {
	ID          int64      `json:"id"`
	ProgramID   int64      `json:"program_id"`
	ScanType    string     `json:"scan_type"`
	Status      string     `json:"status"`
	Progress    int        `json:"progress"`
	AssetsFound int        `json:"assets_found"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Anomaly represents a detected security anomaly (future use)
type Anomaly struct {
	ID            int64                  `json:"id"`
	ProgramID     int64                  `json:"program_id"`
	ProgramName   string                 `json:"program_name"`
	Type          string                 `json:"type"`
	Description   string                 `json:"description"`
	PriorityScore float64                `json:"priority_score"`
	DetectedAt    time.Time              `json:"detected_at"`
	IsReviewed    bool                   `json:"is_reviewed"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

# Server Integration Guide

This document explains how the CLI client integrates with the Recontronic Server.

## Server Overview

**Repository:** https://github.com/Presstronic/recontronic-server

**Current Status:**
- ✅ Authentication system fully implemented
- 🔄 Program management endpoints (planned)
- 🔄 Scan management endpoints (planned)
- 🔄 Anomaly detection system (planned)

## Server Endpoints

### Currently Available (v1)

All authentication endpoints are at base URL `/api/v1/auth/`:

| Method | Endpoint | Auth Required | Purpose |
|--------|----------|---------------|---------|
| POST | `/auth/register` | No | Create user account |
| POST | `/auth/login` | No | Get API key |
| GET | `/auth/me` | Yes | Get current user |
| POST | `/auth/keys` | Yes | Create new API key |
| GET | `/auth/keys` | Yes | List all API keys |
| DELETE | `/auth/keys/{id}` | Yes | Revoke API key |

### Planned (Future Versions)

Based on the vision document, these endpoints are planned:

**Program Management:**
- `POST /api/v1/programs` - Add program
- `GET /api/v1/programs` - List programs
- `GET /api/v1/programs/{id}` - Get program details
- `PATCH /api/v1/programs/{id}` - Update program
- `DELETE /api/v1/programs/{id}` - Delete program

**Scan Management:**
- `POST /api/v1/scans` - Trigger scan
- `GET /api/v1/scans` - List scans
- `GET /api/v1/scans/{id}` - Get scan status

**Anomaly Management:**
- `GET /api/v1/anomalies` - List anomalies
- `GET /api/v1/anomalies/{id}` - Get anomaly details
- `PATCH /api/v1/anomalies/{id}` - Update anomaly (mark reviewed)

## Server Configuration

### Default Settings

```yaml
# Server runs on port 8080 by default
Server Port: 8080
Database: TimescaleDB (PostgreSQL)
Timeouts:
  Read: 15 seconds
  Write: 15 seconds
  Idle: 60 seconds
```

### Environment Variables

Server uses `RECONTRONIC_` prefix for environment variables:

```bash
RECONTRONIC_DATABASE_HOST=localhost
RECONTRONIC_DATABASE_PORT=5432
RECONTRONIC_DATABASE_USER=recontronic
RECONTRONIC_DATABASE_PASSWORD=secure_password
RECONTRONIC_DATABASE_NAME=recontronic_db
RECONTRONIC_SERVER_PORT=8080
```

## API Request/Response Format

### Request Format

All requests use JSON:

```http
POST /api/v1/auth/login HTTP/1.1
Host: api.recontronic.example.com
Content-Type: application/json
Authorization: Bearer rct_YourApiKey123  # For protected endpoints

{
  "field1": "value1",
  "field2": "value2"
}
```

### Response Format

**Success responses:**
```json
{
  "field1": "value1",
  "field2": "value2"
}
```

**Error responses:**
```json
{
  "error": "Descriptive error message"
}
```

### HTTP Status Codes

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful GET, DELETE |
| 201 | Created | Successful POST (resource created) |
| 400 | Bad Request | Validation error, malformed JSON |
| 401 | Unauthorized | Missing/invalid API key |
| 404 | Not Found | Resource doesn't exist |
| 500 | Internal Server Error | Server-side error |

## Authentication

See [AUTHENTICATION.md](./AUTHENTICATION.md) for complete authentication details.

**Quick Summary:**
1. Register once: `POST /api/v1/auth/register`
2. Login to get API key: `POST /api/v1/auth/login`
3. Use API key in all requests: `Authorization: Bearer rct_...`

## Data Models

### User

```go
type User struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    IsActive  bool      `json:"is_active"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### API Key

```go
type APIKey struct {
    ID         int64      `json:"id"`
    UserID     int64      `json:"user_id"`
    Name       string     `json:"name"`
    KeyPrefix  string     `json:"key_prefix"`
    ExpiresAt  *time.Time `json:"expires_at,omitempty"`
    LastUsedAt *time.Time `json:"last_used_at,omitempty"`
    IsActive   bool       `json:"is_active"`
    CreatedAt  time.Time  `json:"created_at"`
    // PlainKey only returned during creation
    PlainKey   string     `json:"plain_key,omitempty"`
}
```

### Program (Planned)

```go
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
```

### Scan (Planned)

```go
type Scan struct {
    ID           int64      `json:"id"`
    ProgramID    int64      `json:"program_id"`
    ScanType     string     `json:"scan_type"`
    Status       string     `json:"status"`
    Progress     int        `json:"progress"`
    AssetsFound  int        `json:"assets_found"`
    StartedAt    *time.Time `json:"started_at,omitempty"`
    CompletedAt  *time.Time `json:"completed_at,omitempty"`
}
```

### Anomaly (Planned)

```go
type Anomaly struct {
    ID                  int64                  `json:"id"`
    ProgramID           int64                  `json:"program_id"`
    ProgramName         string                 `json:"program_name"`
    Type                string                 `json:"type"`
    Description         string                 `json:"description"`
    PriorityScore       float64                `json:"priority_score"`
    DetectedAt          time.Time              `json:"detected_at"`
    IsReviewed          bool                   `json:"is_reviewed"`
    Metadata            map[string]interface{} `json:"metadata,omitempty"`
}
```

## CLI Client Implementation Notes

### REST Client

The CLI should implement a REST client in `pkg/client/rest.go`:

```go
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type RestClient struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}

func NewRestClient(baseURL, apiKey string, timeout time.Duration) *RestClient {
    return &RestClient{
        baseURL: baseURL,
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: timeout,
        },
    }
}

func (c *RestClient) doRequest(ctx context.Context, method, path string, body, response interface{}) error {
    // Implementation details
    // - Add Authorization header if apiKey is set
    // - Marshal body to JSON
    // - Send request
    // - Handle errors
    // - Unmarshal response
}
```

### Error Handling

Implement typed errors:

```go
type APIError struct {
    StatusCode int
    Message    string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// Usage in client:
if resp.StatusCode != http.StatusOK {
    var errResp struct {
        Error string `json:"error"`
    }
    json.NewDecoder(resp.Body).Decode(&errResp)
    return &APIError{
        StatusCode: resp.StatusCode,
        Message:    errResp.Error,
    }
}
```

### Timeout Configuration

Respect user-configured timeouts:

```go
// From config
timeout, _ := time.ParseDuration(config.Timeout) // e.g., "30s"

client := NewRestClient(
    config.Server,
    config.APIKey,
    timeout,
)
```

### Request Logging (Debug Mode)

Log requests/responses in debug mode:

```go
if config.LogLevel == "debug" {
    log.Printf("→ %s %s", method, url)
    log.Printf("  Headers: %v", headers)
    if body != nil {
        log.Printf("  Body: %s", bodyJSON)
    }

    log.Printf("← %d %s", resp.StatusCode, resp.Status)
    log.Printf("  Body: %s", responseBody)
}
```

## Testing Against the Server

### Local Development

1. Clone the server repository:
   ```bash
   git clone https://github.com/Presstronic/recontronic-server.git
   cd recontronic-server
   ```

2. Start the server with Docker Compose:
   ```bash
   docker-compose up -d
   make run
   ```

3. Server will be available at `http://localhost:8080`

4. Configure CLI client:
   ```bash
   recon-cli config set server http://localhost:8080
   ```

### Integration Tests

Tag integration tests to run only when server is available:

```go
// +build integration

func TestLogin(t *testing.T) {
    if os.Getenv("RECON_SERVER") == "" {
        t.Skip("Integration tests require RECON_SERVER")
    }

    // Test implementation
}
```

Run integration tests:
```bash
export RECON_SERVER=http://localhost:8080
go test -tags=integration ./...
```

## Coordinating with Server Development

### Current Status

The server currently has **authentication endpoints only**. The CLI should:

1. **Phase 1 (Immediate):**
   - Implement authentication commands (`auth register`, `auth login`, `auth keys`)
   - Test against live server
   - Provide feedback to server team

2. **Phase 2 (When server adds endpoints):**
   - Implement program/scan/anomaly commands
   - Coordinate data model changes
   - Test integration

### Communication

When server adds new endpoints, they should:
- Update their README/documentation
- Provide request/response examples
- Notify CLI team of any breaking changes

When CLI needs new endpoints, request:
- Clear specification of endpoint behavior
- Expected request/response format
- Error scenarios and status codes

## Version Compatibility

**Current Versions:**
- Server: In development (pre-v1.0)
- CLI: In development (pre-v1.0)

**Planned Versioning:**
- Both will use semantic versioning (v1.0.0, v1.1.0, etc.)
- API version in URL path (`/api/v1/`, `/api/v2/`, etc.)
- CLI should support multiple API versions when possible

**Breaking Changes:**
- Server will increment API version for breaking changes
- CLI can support multiple API versions
- Deprecation notices provided in advance

---

**Related Documentation:**
- [Authentication Guide](./AUTHENTICATION.md)
- [Server Repository](https://github.com/Presstronic/recontronic-server)
- [API Reference](./API-REFERENCE.md) (coming soon)

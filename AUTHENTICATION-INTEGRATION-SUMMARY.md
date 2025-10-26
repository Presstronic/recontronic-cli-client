# Authentication Integration Summary

Based on analysis of the [Recontronic Server repository](https://github.com/Presstronic/recontronic-server), I've added complete authentication support to the CLI client project.

## What Was Added

### 📄 Documentation (2 new files)

1. **docs/AUTHENTICATION.md** (13 KB)
   - Complete authentication guide
   - All auth endpoints documented
   - Request/response examples for every operation
   - Security specifications
   - Best practices for development and production
   - Troubleshooting guide
   - Environment variable reference

2. **docs/SERVER-INTEGRATION.md** (8 KB)
   - Server overview and status
   - All available and planned endpoints
   - Data models (User, APIKey, Program, Scan, Anomaly)
   - REST client implementation guidance
   - Testing strategies
   - Version compatibility notes

### 🎯 New Issues (12 authentication issues)

Added **RECON-051 through RECON-062** to `mvp-issues.csv`:

| Issue | Type | Priority | Title |
|-------|------|----------|-------|
| RECON-051 | User Story | Critical | Auth Command - User Registration |
| RECON-052 | User Story | Critical | Auth Command - User Login |
| RECON-053 | User Story | High | Auth Command - Who Am I |
| RECON-054 | User Story | Medium | Auth Command - Create API Key |
| RECON-055 | User Story | Medium | Auth Command - List API Keys |
| RECON-056 | User Story | Medium | Auth Command - Revoke API Key |
| RECON-057 | User Story | Low | Auth Command - Logout |
| RECON-058 | Tech Story | High | Implement Secure Password Input |
| RECON-059 | Tech Story | High | Implement API Key Storage with Secure Permissions |
| RECON-060 | Tech Story | Medium | Implement Authentication Middleware for REST Client |
| RECON-061 | Tech Story | Low | Implement API Key Validation |
| RECON-062 | User Story | Low | Auth Command - Rotate API Key |

**New Total:** 62 issues (was 50, added 12)

## Server API Details

### Current Server Status

The Recontronic Server has **authentication fully implemented**:

✅ **Available Now:**
- User registration
- User login with API key generation
- API key management (create, list, revoke)
- Current user retrieval

🔄 **Planned (Not Yet Implemented):**
- Program management endpoints
- Scan management endpoints
- Anomaly detection endpoints
- gRPC streaming

### Authentication Endpoints

All endpoints are at base URL `/api/v1/auth/`:

#### Public (No Auth Required)

**Register:**
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",       // 3-50 alphanumeric chars
  "email": "john@example.com", // Valid email format
  "password": "SecureP@ss123"  // 8-72 characters
}

Response (201):
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "is_active": true,
  "created_at": "2025-10-26T10:00:00Z",
  "updated_at": "2025-10-26T10:00:00Z"
}
```

**Login:**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "johndoe",
  "password": "SecureP@ss123"
}

Response (200):
{
  "user": { /* user object */ },
  "api_key": "rct_AbCdEf123456789...",  // SAVE THIS!
  "key_id": 1,
  "message": "Login successful. Save this API key securely - it won't be shown again."
}
```

#### Protected (Requires Bearer Token)

All protected requests need this header:
```http
Authorization: Bearer rct_YourApiKey123...
```

**Get Current User:**
```http
GET /api/v1/auth/me
Authorization: Bearer rct_...

Response (200):
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "is_active": true,
  "created_at": "2025-10-26T10:00:00Z",
  "updated_at": "2025-10-26T10:00:00Z"
}
```

**Create API Key:**
```http
POST /api/v1/auth/keys
Authorization: Bearer rct_...
Content-Type: application/json

{
  "name": "Production Server",         // Optional
  "expires_at": "2026-01-24T10:00:00Z" // Optional
}

Response (201):
{
  "id": 2,
  "user_id": 1,
  "name": "Production Server",
  "key_prefix": "rct_AbCd",
  "plain_key": "rct_NewKey987654321...",  // SAVE THIS!
  "expires_at": "2026-01-24T10:00:00Z",
  "is_active": true,
  "created_at": "2025-10-26T10:30:00Z"
}
```

**List API Keys:**
```http
GET /api/v1/auth/keys
Authorization: Bearer rct_...

Response (200):
{
  "api_keys": [
    {
      "id": 1,
      "user_id": 1,
      "name": "Main Key",
      "key_prefix": "rct_AbCd",
      "expires_at": null,
      "last_used_at": "2025-10-26T12:00:00Z",
      "is_active": true,
      "created_at": "2025-10-26T10:00:00Z"
    }
  ],
  "total": 1
}
```

**Revoke API Key:**
```http
DELETE /api/v1/auth/keys/{id}
Authorization: Bearer rct_...

Response (200):
{
  "message": "API key revoked successfully"
}
```

### Security Features

**Password Hashing:**
- Algorithm: Argon2id
- Memory: 64MB
- Iterations: 3
- Parallelism: 2
- Salt: 16 random bytes per password
- Constant-time comparison (prevents timing attacks)

**API Key Security:**
- Generation: 256-bit cryptographically random values
- Format: `rct_<base64-encoded-bytes>`
- Storage: SHA-256 hashed (plaintext never stored)
- Prefix: First 8 characters stored for identification
- Transport: HTTPS required in production

## Updated Implementation Roadmap

### Phase 1A: Authentication (Week 1) ⭐ NEW

**Priority: Implement authentication FIRST since server is ready**

1. RECON-001 - Initialize Go Module
2. RECON-029 - Create Makefile
3. RECON-003 - Implement Root Command
4. RECON-002 - Configuration System
5. RECON-008 - Define Data Models (add User, APIKey models)
6. RECON-007 - REST API Client (implement auth endpoints)
7. **RECON-051 - Auth Register** ⭐ NEW
8. **RECON-052 - Auth Login** ⭐ NEW
9. **RECON-053 - Auth Whoami** ⭐ NEW
10. **RECON-058 - Secure Password Input** ⭐ NEW
11. **RECON-059 - API Key Storage** ⭐ NEW
12. **RECON-060 - Auth Middleware** ⭐ NEW

**Outcome:** Fully functional authentication, users can register, login, and manage API keys.

### Phase 1B: Additional Auth Features (Week 2)

13. **RECON-054 - Create API Keys** ⭐ NEW
14. **RECON-055 - List API Keys** ⭐ NEW
15. **RECON-056 - Revoke API Keys** ⭐ NEW
16. **RECON-057 - Logout** ⭐ NEW
17. **RECON-061 - API Key Validation** ⭐ NEW
18. **RECON-062 - Rotate API Keys** ⭐ NEW
19. RECON-019 - Output Formatting
20. RECON-031 - Unit Tests for API Client

**Outcome:** Complete auth command suite, ready for integration testing with server.

### Phase 2: Wait for Server Program Endpoints (Week 2-3)

**Note:** Program, scan, and anomaly endpoints are not yet implemented on the server.

Options:
1. **Pause and help server team** implement program endpoints
2. **Continue with CLI infrastructure:**
   - RECON-020 - Error Handling Framework
   - RECON-028 - Logging System
   - RECON-029 - Makefile enhancements
   - RECON-030 - CI/CD Pipeline
   - RECON-034 - Version Command
   - RECON-035 - Shell Completion
   - RECON-046 - Health Check Command

3. **Mock server responses** for program/scan/anomaly commands (for testing)

### Phase 3: Program Management (When Server Ready)

Implement program commands once server has endpoints:
- RECON-009 - Program Add
- RECON-010 - Program List
- RECON-011 - Program Get
- RECON-012 - Program Delete
- RECON-047 - Program Update

### Phase 4: Scan & Anomaly Management (When Server Ready)

Continue with scan and anomaly commands once endpoints are available.

## Testing Strategy

### Integration Testing with Live Server

**Setup local server:**
```bash
# Clone server repo
git clone https://github.com/Presstronic/recontronic-server.git
cd recontronic-server

# Start with Docker Compose
docker-compose up -d
make run
```

**Test CLI against server:**
```bash
# Configure CLI
recon-cli config set server http://localhost:8080

# Test auth flow
recon-cli auth register
recon-cli auth login
recon-cli auth whoami
recon-cli auth keys create --name "Test Key"
recon-cli auth keys list
```

### Example Integration Test

```go
// +build integration

func TestAuthFlow(t *testing.T) {
    if os.Getenv("RECON_SERVER") == "" {
        t.Skip("Integration tests require RECON_SERVER")
    }

    client := NewRestClient(os.Getenv("RECON_SERVER"), "", 30*time.Second)

    // Test registration
    user, err := client.Register("testuser", "test@example.com", "password123")
    require.NoError(t, err)
    assert.Equal(t, "testuser", user.Username)

    // Test login
    loginResp, err := client.Login("testuser", "password123")
    require.NoError(t, err)
    assert.NotEmpty(t, loginResp.APIKey)
    assert.True(t, strings.HasPrefix(loginResp.APIKey, "rct_"))

    // Test authenticated request
    client.SetAPIKey(loginResp.APIKey)
    me, err := client.GetCurrentUser()
    require.NoError(t, err)
    assert.Equal(t, "testuser", me.Username)
}
```

## CLI Usage Examples

### Registration Flow
```bash
$ recon-cli auth register
Username: demian
Email: demian@example.com
Password: [hidden]
Confirm password: [hidden]

✓ Registration successful!
Account created for: demian

Next step: Login to get your API key
$ recon-cli auth login
```

### Login Flow
```bash
$ recon-cli auth login
Username: demian
Password: [hidden]

✓ Login successful!

Your API key: rct_AbCdEf123456789...

⚠️  IMPORTANT: Save this key securely!
   It has been saved to: ~/.recon-cli/config.yaml
   This key will not be shown again.

You're now authenticated and ready to use the CLI.
```

### Verify Authentication
```bash
$ recon-cli auth whoami
Username:     demian
Email:        demian@example.com
Account ID:   1
Status:       Active
Created:      2025-10-26 10:00:00
API Key:      rct_AbCd... (prefix)
```

### Manage API Keys
```bash
$ recon-cli auth keys create --name "CI/CD Pipeline" --expires-in 90d
✓ New API key created!

API Key: rct_NewKey987654321...
Name:    CI/CD Pipeline
Expires: 2026-01-24

⚠️  Save this key! It won't be shown again.

$ recon-cli auth keys list
ID  NAME             PREFIX      LAST USED      EXPIRES    STATUS
1   Main Key         rct_AbCd    2 hours ago    Never      Active
2   CI/CD Pipeline   rct_NewK    Never          in 90d     Active

$ recon-cli auth keys revoke --id 2
Are you sure you want to revoke 'CI/CD Pipeline' (rct_NewK...)? [y/N]: y
✓ API key revoked successfully
```

## Configuration File

After authentication, `~/.recon-cli/config.yaml`:

```yaml
server: http://localhost:8080
api_key: rct_AbCdEf123456789...
timeout: 30s
output_format: table
log_level: info
```

**Security:** File automatically set to `0600` permissions (owner read/write only).

## Next Steps

1. **Review new documentation:**
   - `docs/AUTHENTICATION.md` - Complete auth guide
   - `docs/SERVER-INTEGRATION.md` - Server integration details

2. **Import new issues** (RECON-051 through RECON-062) into your issue tracker

3. **Start with Phase 1A:**
   - Focus on authentication first (server is ready!)
   - Can test immediately against live server
   - Provides foundation for all other commands

4. **Coordinate with server team:**
   - Test auth integration
   - Provide feedback on API
   - Plan for program/scan/anomaly endpoints

## Summary

**Before:** 50 issues, no authentication details
**After:** 62 issues, complete authentication system documented

**Ready to implement:**
- ✅ Full authentication command suite
- ✅ Secure password input
- ✅ API key management
- ✅ Integration testing strategy
- ✅ Server API fully documented

**Waiting on server:**
- ⏳ Program management endpoints
- ⏳ Scan management endpoints
- ⏳ Anomaly detection endpoints
- ⏳ gRPC streaming

**Recommendation:** Start with Phase 1A (authentication) since the server is ready. This provides immediate value and allows real integration testing.

---

**Total Issues:** 62 (50 original + 12 authentication)
**Total Story Points:** ~235 points
**Auth Story Points:** ~35 points (~17-18 developer-days)
**Files Created:** 2 documentation files
**Ready for Development:** Yes! ✅

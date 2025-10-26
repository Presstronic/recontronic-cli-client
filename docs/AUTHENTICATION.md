# Authentication Guide

This document explains how the Recontronic CLI authenticates with the server.

## Authentication Flow

The Recontronic Server uses **API key-based authentication** optimized for CLI clients.

### 1. User Registration (One-time)

```bash
recon-cli auth register
# Interactive prompts for:
# - Username (3-50 alphanumeric characters)
# - Email (valid email format)
# - Password (8-72 characters, hidden input)
```

**API Request:**
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecureP@ssw0rd123"
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "is_active": true,
  "created_at": "2025-10-26T10:00:00Z",
  "updated_at": "2025-10-26T10:00:00Z"
}
```

### 2. User Login (Get API Key)

```bash
recon-cli auth login
# Prompts for username and password
# Returns API key - SAVE THIS SECURELY!
```

**API Request:**
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "johndoe",
  "password": "SecureP@ssw0rd123"
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "is_active": true,
    "created_at": "2025-10-26T10:00:00Z",
    "updated_at": "2025-10-26T10:00:00Z"
  },
  "api_key": "rct_AbCdEf123456789...",
  "key_id": 1,
  "message": "Login successful. Save this API key securely - it won't be shown again."
}
```

**⚠️ IMPORTANT**: The API key is shown **only once**. Save it immediately!

### 3. Configure CLI with API Key

The CLI will automatically save the API key to `~/.recon-cli/config.yaml` after login.

**Manual configuration:**
```bash
recon-cli config set api-key rct_AbCdEf123456789...
recon-cli config set server https://your-server.com
```

**Config file location:** `~/.recon-cli/config.yaml`
```yaml
server: https://api.recontronic.example.com
api_key: rct_AbCdEf123456789...
timeout: 30s
output_format: table
log_level: info
```

**File permissions:** Automatically set to `0600` (read/write for owner only)

### 4. Using Authenticated Requests

All subsequent CLI commands automatically include the API key:

```bash
# The CLI adds this header automatically:
Authorization: Bearer rct_AbCdEf123456789...
```

## API Key Management

### Generate Additional Keys

Create multiple API keys for different machines or environments:

```bash
recon-cli auth keys create --name "Production Server"
recon-cli auth keys create --name "CI/CD Pipeline" --expires-in 90d
```

**API Request:**
```http
POST /api/v1/auth/keys
Authorization: Bearer rct_AbCdEf123456789...
Content-Type: application/json

{
  "name": "Production Server",
  "expires_at": "2026-01-24T10:00:00Z"  // optional
}
```

**Response (201 Created):**
```json
{
  "id": 2,
  "user_id": 1,
  "name": "Production Server",
  "key_prefix": "rct_AbCd",
  "plain_key": "rct_NewKey987654321...",
  "expires_at": "2026-01-24T10:00:00Z",
  "is_active": true,
  "created_at": "2025-10-26T10:30:00Z"
}
```

### List All API Keys

```bash
recon-cli auth keys list
```

**API Request:**
```http
GET /api/v1/auth/keys
Authorization: Bearer rct_AbCdEf123456789...
```

**Response (200 OK):**
```json
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
    },
    {
      "id": 2,
      "user_id": 1,
      "name": "Production Server",
      "key_prefix": "rct_XyZw",
      "expires_at": "2026-01-24T10:00:00Z",
      "last_used_at": null,
      "is_active": true,
      "created_at": "2025-10-26T10:30:00Z"
    }
  ],
  "total": 2
}
```

**Note:** Plain API keys are **never** returned from this endpoint (only shown once during creation).

### Revoke API Keys

```bash
recon-cli auth keys revoke --id 2
# or
recon-cli auth keys revoke --name "Production Server"
```

**API Request:**
```http
DELETE /api/v1/auth/keys/2
Authorization: Bearer rct_AbCdEf123456789...
```

**Response (200 OK):**
```json
{
  "message": "API key revoked successfully"
}
```

### View Current User

```bash
recon-cli auth whoami
```

**API Request:**
```http
GET /api/v1/auth/me
Authorization: Bearer rct_AbCdEf123456789...
```

**Response (200 OK):**
```json
{
  "id": 1,
  "username": "johndoe",
  "email": "john@example.com",
  "is_active": true,
  "created_at": "2025-10-26T10:00:00Z",
  "updated_at": "2025-10-26T10:00:00Z"
}
```

## Security Specifications

### Password Security
- **Algorithm**: Argon2id
- **Memory**: 64MB
- **Iterations**: 3
- **Parallelism**: 2
- **Salt**: 16 random bytes per password
- **Comparison**: Constant-time to prevent timing attacks

### API Key Security
- **Generation**: 256-bit cryptographically random values
- **Format**: `rct_<base64-encoded-bytes>`
- **Storage**: SHA-256 hashed (plaintext never stored on server)
- **Prefix**: First 8 characters stored for identification
- **Transport**: Always use HTTPS in production

## Error Handling

All errors return JSON format:

```json
{
  "error": "Error message describing what went wrong"
}
```

### Common HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| 200 | Success | Login successful |
| 201 | Created | User registered, API key created |
| 400 | Bad Request | Invalid username format, missing fields |
| 401 | Unauthorized | Invalid credentials, missing API key |
| 404 | Not Found | API key ID doesn't exist |
| 500 | Server Error | Database connection failed |

### Validation Errors

**Username validation:**
- Required field
- 3-50 characters
- Alphanumeric only (letters and numbers)

**Email validation:**
- Required field
- Must be valid email format

**Password validation:**
- Required field
- 8-72 characters minimum/maximum
- No specific complexity requirements (but strong passwords recommended)

**API Key Name validation:**
- Optional field
- Maximum 100 characters

## Best Practices

### For Development

1. **Use separate keys per machine:**
   ```bash
   recon-cli auth keys create --name "MacBook Pro"
   recon-cli auth keys create --name "Linux Server"
   ```

2. **Set expiration for CI/CD:**
   ```bash
   recon-cli auth keys create --name "GitHub Actions" --expires-in 90d
   ```

3. **Revoke unused keys:**
   ```bash
   recon-cli auth keys list
   recon-cli auth keys revoke --id 5
   ```

### For Production

1. **Never commit API keys to Git:**
   - Config files are in `.gitignore`
   - Use environment variables in CI/CD

2. **Use HTTPS only:**
   ```yaml
   server: https://api.recontronic.example.com  # ✅ Good
   server: http://api.recontronic.example.com   # ❌ Bad (unencrypted)
   ```

3. **Secure file permissions:**
   - CLI automatically sets `~/.recon-cli/config.yaml` to `0600`
   - Verify: `ls -la ~/.recon-cli/config.yaml`
   - Should show: `-rw-------` (owner read/write only)

4. **Rotate keys periodically:**
   ```bash
   # Create new key
   recon-cli auth keys create --name "New Production Key"

   # Update config
   recon-cli config set api-key rct_NewKey...

   # Revoke old key
   recon-cli auth keys revoke --id 1
   ```

### For CI/CD Environments

Use environment variables instead of config files:

```bash
export RECON_SERVER="https://api.recontronic.example.com"
export RECON_API_KEY="rct_AbCdEf123456789..."

# CLI will use environment variables automatically
recon-cli program list
```

**GitHub Actions example:**
```yaml
env:
  RECON_SERVER: ${{ secrets.RECON_SERVER }}
  RECON_API_KEY: ${{ secrets.RECON_API_KEY }}

steps:
  - name: List programs
    run: recon-cli program list
```

## Troubleshooting

### "Invalid API key" error

**Problem:** Your API key is invalid or has been revoked.

**Solutions:**
1. Verify the key in your config: `recon-cli config get api-key`
2. Check if key is active: `recon-cli auth keys list`
3. Login again to get a new key: `recon-cli auth login`

### "Connection refused" error

**Problem:** Can't connect to the server.

**Solutions:**
1. Check server URL: `recon-cli config get server`
2. Verify server is running: `curl https://your-server.com/health`
3. Check network/firewall settings

### "Unauthorized" (401) error

**Problem:** Missing or invalid authentication.

**Solutions:**
1. Ensure API key is set: `recon-cli config get api-key`
2. Re-login if needed: `recon-cli auth login`
3. Check config file exists: `ls ~/.recon-cli/config.yaml`

### Config file permissions error

**Problem:** Config file has incorrect permissions.

**Solution:**
```bash
chmod 600 ~/.recon-cli/config.yaml
```

## Environment Variables Reference

All configuration can be set via environment variables with `RECON_` prefix:

| Variable | Config Key | Example |
|----------|------------|---------|
| `RECON_SERVER` | server | `https://api.recontronic.example.com` |
| `RECON_API_KEY` | api_key | `rct_AbCdEf123456789...` |
| `RECON_TIMEOUT` | timeout | `30s` |
| `RECON_OUTPUT_FORMAT` | output_format | `table`, `json`, `yaml` |
| `RECON_LOG_LEVEL` | log_level | `debug`, `info`, `warn`, `error` |

**Precedence order:**
1. Command-line flags (highest priority)
2. Environment variables
3. Config file
4. Default values (lowest priority)

---

**Next Steps:**
- [Quick Start Guide](../QUICKSTART.md)
- [Configuration Reference](../README.md#configuration)
- [API Documentation](./API-REFERENCE.md) (coming soon)

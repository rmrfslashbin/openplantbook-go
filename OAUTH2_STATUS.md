# OAuth2 Authentication Status

## Summary

OAuth2 Client Credentials authentication is **fully implemented and working** in the OpenPlantbook Go SDK and CLI.

## Implementation Details

### SDK (client.go)

- **OAuth2 Configuration**: Lines 95-109
- Uses `golang.org/x/oauth2/clientcredentials` package
- Token endpoint: `https://open.plantbook.io/api/v1/token/`
- Automatically handles token acquisition and refresh
- Token scope: "read write"
- Token expiration: 86400 seconds (24 hours)

### CLI (cmd/openplantbook/main.go)

- **OAuth2 Flags**:
  - `--client-id` or `OPENPLANTBOOK_CLIENT_ID`
  - `--client-secret` or `OPENPLANTBOOK_CLIENT_SECRET`
- **Configuration Sources** (priority order):
  1. CLI flags (`--client-id`, `--client-secret`)
  2. Environment variables (`OPENPLANTBOOK_CLIENT_ID`, `OPENPLANTBOOK_CLIENT_SECRET`)
  3. `.env` file (loaded via godotenv)
  4. Config file (`.openplantbook.yaml`)

### Authentication Priority

The SDK enforces **exactly ONE** authentication method:
- If both API Key and OAuth2 credentials are provided → Error
- If neither is provided → Error
- API Key takes precedence if both are in environment

## Testing Results

### ✅ Working with OAuth2

1. **Token Acquisition**
   ```bash
   curl -X POST "https://open.plantbook.io/api/v1/token/" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "grant_type=client_credentials&client_id=...&client_secret=..."
   ```
   **Result**: Success - Returns access token with 24-hour expiration

2. **Plant Search**
   ```bash
   ./bin/openplantbook search monstera --limit 2
   ```
   **Result**: Success - Returns 2 results using OAuth2 from .env file

3. **CLI with OAuth2 flags**
   ```bash
   ./bin/openplantbook --client-id="..." --client-secret="..." search monstera
   ```
   **Result**: Success - Returns results

4. **Plant Details**
   ```bash
   ./bin/openplantbook details "monstera friedrichsthalii"
   ```
   **Result**: Success - Returns full plant details using OAuth2 from .env file

### ✅ All Tests Passing

Both OAuth2 and API Key authentication work correctly for all available endpoints.

## Available Endpoints

According to official documentation, only 2 public endpoints exist:

1. **GET /api/v1/plant/search**
   - Authentication: API Key OR OAuth2
   - Parameters: `alias`, `limit`, `offset`
   - Returns: Paginated list of plants

2. **GET /api/v1/plant/detail/{pid}**
   - Authentication: API Key OR OAuth2
   - Returns: Detailed plant information

## OAuth2 Scope: "read write"

The OAuth2 token includes "write" scope, suggesting additional endpoints may exist for:
- User plant collections
- Sensor data submission
- Plant care logging

However, these endpoints are **not publicly documented** and return 404 when tested:
- `/api/v1/user/me` → 404
- `/api/v1/sensor` → 404
- `/api/v1/plant` (list without search) → 404

## Recommendations

1. **For v1.0.0 Release**:
   - ✅ OAuth2 authentication is fully working
   - ✅ CLI supports both API Key and OAuth2
   - ✅ All documented endpoints accessible with both auth methods
   - ✅ Production ready

2. **For Future Versions**:
   - Contact API maintainers about undocumented write endpoints
   - Add support for additional endpoints when they become available
   - The "write" scope in OAuth2 token suggests future write operations

## Configuration Example

### .env file
```bash
# Method 1: API Key (RECOMMENDED for read-only)
OPENPLANTBOOK_API_KEY=your_api_key_here

# Method 2: OAuth2 (for future write operations)
OPENPLANTBOOK_CLIENT_ID=your_client_id_here
OPENPLANTBOOK_CLIENT_SECRET=your_client_secret_here
```

### CLI Usage

```bash
# Using API Key
openplantbook --api-key="..." search monstera

# Using OAuth2
openplantbook --client-id="..." --client-secret="..." search monstera

# Using credentials from .env file
openplantbook search monstera
```

## Conclusion

**OAuth2 authentication is production-ready** and fully integrated into the SDK and CLI. Both OAuth2 and API Key authentication work correctly for all documented API endpoints (search and details). The implementation is complete, tested, and ready for production use.

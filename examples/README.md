# OpenPlantbook SDK Examples

This directory contains example programs demonstrating how to use the OpenPlantbook Go SDK.

## Prerequisites

1. Sign up at [OpenPlantbook](https://open.plantbook.io/) to get API credentials
2. Set up your environment variables (see below)

## Authentication Setup

The SDK supports two authentication methods:

### API Key (Recommended for Simple Use Cases)

```bash
export OPENPLANTBOOK_API_KEY="your-api-key-here"
```

### OAuth2 Client Credentials (Full API Access)

```bash
export OPENPLANTBOOK_CLIENT_ID="your-client-id"
export OPENPLANTBOOK_CLIENT_SECRET="your-client-secret"
```

## Examples

### 1. Basic Search (`basic_search`)

Demonstrates simple plant search functionality using API Key authentication.

**Run:**
```bash
cd basic_search
go run main.go "monstera"
```

**Features:**
- API Key authentication
- Plant search with query parameter
- Display search results
- Command-line argument support

**Expected Output:**
```
Searching for plants matching 'monstera'...

Found 2 plant(s):

1. Monstera deliciosa
   Alias: Monstera
   PID: monstera-deliciosa
   Category: Houseplant

2. Monstera adansonii
   Alias: Swiss Cheese Vine
   PID: monstera-adansonii
   Category: Houseplant
```

---

### 2. Plant Details (`plant_details`)

Demonstrates retrieving detailed plant care information using OAuth2 authentication.

**Run:**
```bash
cd plant_details
go run main.go "monstera-deliciosa"
```

**Features:**
- OAuth2 authentication
- Detailed plant information retrieval
- Language parameter support
- Comprehensive care requirements display

**Expected Output:**
```
Fetching details for plant 'monstera-deliciosa'...

Plant: Monstera deliciosa
Common Name: Monstera
Category: Houseplant

Care Requirements:
==================
Light (Lux): 2500 - 20000
Temperature (°C): 15.0 - 30.0
Humidity (%): 40 - 80
Soil Moisture (%): 15 - 60
Soil EC (μS/cm): 350 - 2000

Image: https://example.com/monstera.jpg
```

---

### 3. Caching Demonstration (`with_caching`)

Demonstrates the SDK's caching behavior and performance benefits.

**Run:**
```bash
cd with_caching
go run main.go
```

**Features:**
- Custom cache configuration
- Custom rate limiting
- Cache hit/miss demonstration
- Performance comparison
- Cache clearing

**Expected Output:**
```
Demonstrating cache behavior with query 'fern'

=== First Request (API call) ===
Found 5 plants in 245ms

=== Second Request (cached) ===
Found 5 plants in 127μs
Speedup: 1929.13x faster

=== Third Request (still cached) ===
Found 5 plants in 89μs

=== Clearing Cache ===
Cache cleared successfully

=== Fourth Request (API call after clear) ===
Found 5 plants in 238ms

=== Cache Behavior Summary ===
Request 1 (API):    245ms
Request 2 (cache):  127μs
Request 3 (cache):  89μs
Request 4 (API):    238ms

Notice how cached requests are significantly faster!
```

## Common Patterns

### Using Context for Timeout Control

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

results, err := client.SearchPlants(ctx, "fern", nil)
```

### Custom HTTP Client

```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}

client, err := openplantbook.New(
    openplantbook.WithAPIKey(apiKey),
    openplantbook.WithHTTPClient(httpClient),
)
```

### Disable Caching

```go
client, err := openplantbook.New(
    openplantbook.WithAPIKey(apiKey),
    openplantbook.WithCache(openplantbook.NewNoOpCache()),
)
```

### Custom Logger

```go
type customLogger struct {
    logger *slog.Logger
}

func (l *customLogger) Log(msg string, keysAndValues ...interface{}) {
    l.logger.Info(msg, keysAndValues...)
}

logger := &customLogger{logger: slog.Default()}

client, err := openplantbook.New(
    openplantbook.WithAPIKey(apiKey),
    openplantbook.WithLogger(logger),
)
```

## Rate Limiting

The SDK includes client-side rate limiting (200 requests/day by default):

```go
// Custom rate limit (100 requests/day)
client, err := openplantbook.New(
    openplantbook.WithAPIKey(apiKey),
    openplantbook.WithRateLimit(100),
)

// Disable rate limiting (for testing)
client, err := openplantbook.New(
    openplantbook.WithAPIKey(apiKey),
    openplantbook.DisableRateLimit(),
)
```

## Error Handling

```go
results, err := client.SearchPlants(ctx, "unknown", nil)
if err != nil {
    switch {
    case errors.Is(err, openplantbook.ErrUnauthorized):
        log.Fatal("Invalid API credentials")
    case errors.Is(err, openplantbook.ErrRateLimitExceeded):
        log.Fatal("Rate limit exceeded")
    case errors.Is(err, openplantbook.ErrNotFound):
        log.Fatal("Plant not found")
    default:
        log.Fatalf("API error: %v", err)
    }
}
```

## Building Examples

To build all examples:

```bash
# From the examples directory
for dir in */; do
    (cd "$dir" && go build -o "${dir%/}" .)
done
```

## Additional Resources

- [OpenPlantbook API Documentation](https://open.plantbook.io/docs/)
- [SDK Documentation](../README.md)
- [API Reference](https://pkg.go.dev/github.com/rmrfslashbin/openplantbook-go)

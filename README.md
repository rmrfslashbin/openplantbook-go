# OpenPlantbook Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/rmrfslashbin/openplantbook-go.svg)](https://pkg.go.dev/github.com/rmrfslashbin/openplantbook-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/rmrfslashbin/openplantbook-go)](https://goreportcard.com/report/github.com/rmrfslashbin/openplantbook-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A production-ready Go client library for the [OpenPlantbook API](https://open.plantbook.io), providing access to a crowd-sourced database of plant care information.

## Overview

This is the first official Go SDK for OpenPlantbook, enabling developers to:
- Search for plants by common or scientific name
- Retrieve detailed plant care requirements (light, temperature, humidity, soil moisture, etc.)
- Access community-sourced sensor data
- Integrate plant care data into IoT applications, monitoring systems, and AI assistants

## Features

- **Dual Authentication** - Support for both API Key and OAuth2 authentication
- **Intelligent Caching** - Pluggable cache interface with in-memory implementation (1h for search, 24h for details)
- **Rate Limiting** - Built-in client-side rate limiting (200 requests/day default, configurable)
- **Error Handling** - Comprehensive error types with proper wrapping
- **Context Support** - Full `context.Context` support for cancellation and timeouts
- **Idiomatic Go** - Follows Go best practices with functional options pattern
- **Zero Dependencies** - Only uses Go standard library and official `golang.org/x` packages
- **CLI Tool** - Fully-featured command-line interface included
- **Well Tested** - 90%+ test coverage with table-driven tests

## Installation

```bash
go get github.com/rmrfslashbin/openplantbook-go
```

## Quick Start

### API Key Authentication (Recommended)

```go
package main

import (
    "context"
    "fmt"
    "log"

    openplantbook "github.com/rmrfslashbin/openplantbook-go"
)

func main() {
    // Create client with API key
    client, err := openplantbook.New(
        openplantbook.WithAPIKey("your-api-key-here"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Search for plants
    results, err := client.SearchPlants(context.Background(), "monstera", &openplantbook.SearchOptions{
        Limit: 10,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d plants\n", len(results))

    // Get detailed plant information
    if len(results) > 0 {
        details, err := client.GetPlantDetails(context.Background(), results[0].PID, &openplantbook.DetailOptions{
            Language: "en",
        })
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("\nPlant: %s (%s)\n", details.Alias, details.DisplayPID)
        fmt.Printf("Temperature: %.1f-%.1f°C\n", details.MinTemp, details.MaxTemp)
        fmt.Printf("Light: %d-%d lux\n", details.MinLightLux, details.MaxLightLux)
        fmt.Printf("Humidity: %d-%d%%\n", details.MinEnvHumid, details.MaxEnvHumid)
    }
}
```

### OAuth2 Authentication (Full API Access)

```go
// Create client with OAuth2 credentials
client, err := openplantbook.New(
    openplantbook.WithOAuth2("client-id", "client-secret"),
)
```

## Authentication

The SDK supports two authentication methods:

### API Key
- **Use Case**: Read-only operations (search, get details)
- **Pros**: Simple, never expires, no token management
- **Cons**: Limited to read operations
- **Header**: `Authorization: Token <api-key>`

### OAuth2 Client Credentials
- **Use Case**: Full API access (future write operations)
- **Pros**: Access to all endpoints
- **Cons**: Tokens expire, more complex setup
- **Header**: `Authorization: Bearer <access-token>`

Get your credentials at: https://open.plantbook.io/

## Configuration Options

The SDK supports extensive configuration through functional options:

```go
client, err := openplantbook.New(
    // Authentication (required - choose one)
    openplantbook.WithAPIKey("your-key"),
    // OR
    openplantbook.WithOAuth2("client-id", "client-secret"),

    // Optional configuration
    openplantbook.WithBaseURL("https://custom-api.example.com"),
    openplantbook.WithCache(customCache),
    openplantbook.WithRateLimit(100), // requests per day
    openplantbook.WithHTTPClient(customHTTPClient),
    openplantbook.WithLogger(logger),
    openplantbook.DisableRateLimit(), // for testing
)
```

## Examples

See the [examples](./examples/) directory for complete working examples:

- **[basic_search](./examples/basic_search/)** - Simple plant search with API Key
- **[plant_details](./examples/plant_details/)** - Detailed plant information with OAuth2
- **[with_caching](./examples/with_caching/)** - Caching behavior demonstration

## CLI Tool

A full-featured command-line interface is included:

### Installation

```bash
# From source
git clone https://github.com/rmrfslashbin/openplantbook-go.git
cd openplantbook-go
make build-cli
sudo cp bin/openplantbook /usr/local/bin/

# Or using go install
go install github.com/rmrfslashbin/openplantbook-go/cmd/openplantbook@latest
```

### Usage

```bash
# Set up authentication
export OPENPLANTBOOK_API_KEY="your-api-key"

# Search for plants
openplantbook search monstera

# Get plant details
openplantbook details monstera-deliciosa

# JSON output for scripting
openplantbook search fern --json | jq '.[] | .pid'

# Get help
openplantbook help
```

See [CLI Documentation](./cmd/openplantbook/README.md) for complete usage guide.

## API Reference

### Plant Search

```go
results, err := client.SearchPlants(ctx, "query", &openplantbook.SearchOptions{
    Limit:      10,    // Max results to return
    UserPlants: false, // Search user-contributed plants only
})
```

**Returns**: Array of `PlantSearchResult` with:
- `PID` - Plant identifier
- `DisplayPID` - Scientific name
- `Alias` - Common name
- `Category` - Plant category

### Plant Details

```go
details, err := client.GetPlantDetails(ctx, "plant-id", &openplantbook.DetailOptions{
    Language: "en", // ISO 639-1 language code
})
```

**Returns**: `PlantDetails` with comprehensive care information:
- Light requirements (min/max lux)
- Temperature range (°C)
- Humidity levels (%)
- Soil moisture requirements (%)
- Soil electrical conductivity (μS/cm)
- Image URL
- Category and names

## Error Handling

The SDK provides typed errors for common scenarios:

```go
results, err := client.SearchPlants(ctx, "query", nil)
if err != nil {
    switch {
    case errors.Is(err, openplantbook.ErrUnauthorized):
        // Invalid credentials
    case errors.Is(err, openplantbook.ErrRateLimitExceeded):
        // Too many requests
    case errors.Is(err, openplantbook.ErrNotFound):
        // Plant not found
    case errors.Is(err, openplantbook.ErrNoAuthProvided):
        // Missing authentication
    case errors.Is(err, openplantbook.ErrMultipleAuthMethods):
        // Both API key and OAuth2 provided
    default:
        // Other errors
    }
}
```

## Caching

The SDK includes intelligent caching out of the box:

- **Search results**: Cached for 1 hour
- **Plant details**: Cached for 24 hours
- **Cache hits**: Significantly faster (1000x+ speedup)

### Custom Cache

Implement the `Cache` interface for custom caching (Redis, etc.):

```go
type Cache interface {
    Get(key string) ([]byte, bool)
    Set(key string, value []byte, ttl time.Duration)
    Delete(key string)
    Clear()
}

// Use custom cache
client, err := openplantbook.New(
    openplantbook.WithAPIKey("key"),
    openplantbook.WithCache(myRedisCache),
)
```

### Disable Caching

```go
client, err := openplantbook.New(
    openplantbook.WithAPIKey("key"),
    openplantbook.WithCache(openplantbook.NewNoOpCache()),
)
```

## Rate Limiting

Client-side rate limiting prevents exceeding API quotas:

```go
// Default: 200 requests per day
client, err := openplantbook.New(
    openplantbook.WithAPIKey("key"),
)

// Custom rate limit
client, err := openplantbook.New(
    openplantbook.WithAPIKey("key"),
    openplantbook.WithRateLimit(100), // 100 requests/day
)

// Disable for testing
client, err := openplantbook.New(
    openplantbook.WithAPIKey("key"),
    openplantbook.DisableRateLimit(),
)
```

## Logging

Optional logging interface for debugging:

```go
type Logger interface {
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
}

// Enable logging
client, err := openplantbook.New(
    openplantbook.WithAPIKey("key"),
    openplantbook.WithLogger(myLogger),
)
```

## Testing

```bash
# Run all tests
make test

# View coverage report
make coverage

# Run specific tests
go test -v -run TestClient_SearchPlants

# Run with race detector
go test -v -race ./...
```

Current test coverage: **90.5%**

## Building

```bash
# Run tests
make test

# Build CLI
make build-cli

# Build for all platforms
make build-cli-all

# Install CLI locally
make install-cli

# Clean build artifacts
make clean
```

## Project Structure

```
openplantbook-go/
├── cache.go           # Cache interface and implementations
├── client.go          # HTTP client and authentication
├── errors.go          # Error types and handling
├── models.go          # API data structures
├── options.go         # Functional options
├── plants.go          # Plant search and details API
├── cmd/
│   └── openplantbook/ # CLI tool
├── examples/          # Usage examples
└── testdata/          # Test fixtures
```

## Dependencies

Minimal dependencies, all from official sources:

- `golang.org/x/oauth2` - OAuth2 implementation
- `golang.org/x/time` - Rate limiting

## Roadmap

- [ ] Redis cache implementation
- [ ] Write operations support (when API supports it)
- [ ] Batch operations
- [ ] Pagination helpers
- [ ] Webhook support
- [ ] MCP server implementation

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure:
- Tests pass (`make test`)
- Code is formatted (`go fmt ./...`)
- Coverage remains above 80%

## Versioning

This project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html) (SemVer).

## License

MIT License - see [LICENSE](LICENSE) for details.

Copyright (c) 2025 Robert Sigler

## Author

**Robert Sigler**
Email: code@sigler.io
GitHub: [@rmrfslashbin](https://github.com/rmrfslashbin)

## Resources

- [OpenPlantbook Website](https://open.plantbook.io)
- [OpenPlantbook API Documentation](https://documenter.getpostman.com/view/12627470/TVsxBRjD)
- [Get API Credentials](https://open.plantbook.io/apikey/show/)
- [Report Issues](https://github.com/rmrfslashbin/openplantbook-go/issues)

## Acknowledgments

- [OpenPlantbook](https://open.plantbook.io) team for providing the free API
- OpenPlantbook community for maintaining the plant database
- All contributors to this SDK

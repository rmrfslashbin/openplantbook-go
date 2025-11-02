# OpenPlantbook Go SDK

A production-ready Go client library for the [OpenPlantbook API](https://open.plantbook.io), providing access to a crowd-sourced database of plant care information.

## Overview

This is the first official Go SDK for OpenPlantbook, enabling developers to:
- Search for plants by common or scientific name
- Retrieve detailed plant care requirements (moisture, temperature, light, etc.)
- Access community-sourced sensor data
- Integrate plant care data into IoT applications, monitoring systems, and AI assistants

## Features

- **OAuth2 Authentication** - Secure client credentials flow
- **Rate Limiting** - Built-in client-side rate limiting (200 requests/day)
- **Caching** - Intelligent in-memory caching with pluggable interface
- **Error Handling** - Comprehensive error types with context
- **Context Support** - Full `context.Context` support for cancellation/timeout
- **Idiomatic Go** - Follows Go best practices and conventions
- **CLI Tool** - Reference implementation demonstrating SDK usage

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

    "github.com/rmrfslashbin/openplantbook-go"
)

func main() {
    // Create client with API key (simpler)
    client, err := openplantbook.New(
        openplantbook.WithAPIKey("your-api-key-here"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Search for plants
    results, err := client.SearchPlants(context.Background(), "monstera", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Get plant details
    details, err := client.GetPlantDetails(context.Background(), results[0].PID, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Plant: %s\n", details.Alias)
    fmt.Printf("Temperature: %.1f-%.1f°C\n", details.MinTemp, details.MaxTemp)
}
```

### OAuth2 Authentication (Full API Access)

```go
// Create client with OAuth2 (for write operations)
client, err := openplantbook.New(
    openplantbook.WithOAuth2("client-id", "client-secret"),
)
```

## API Credentials

Get your free API credentials at: https://open.plantbook.io/apikey/show/

## Documentation

- [API Documentation](https://documenter.getpostman.com/view/12627470/TVsxBRjD)
- [Examples](./examples/)
- [GoDoc](https://pkg.go.dev/github.com/rmrfslashbin/openplantbook-go)

## CLI Tool

A command-line tool is included for testing and reference:

```bash
# Install
make install-cli

# Search for plants
openplantbook search monstera

# Get plant details
openplantbook details monstera-deliciosa --format json
```

## Contributing

Contributions welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Author

Robert Sigler (code@sigler.io)

## Acknowledgments

- [OpenPlantbook](https://open.plantbook.io) - Free plant care database
- OpenPlantbook community for maintaining the database

# OpenPlantbook Go SDK v0.1.0

First production-ready release of the OpenPlantbook Go SDK!

## 🎉 Highlights

- **Dual Authentication**: Support for both API Key and OAuth2 Client Credentials
- **Production Ready**: 90.5% test coverage, race detector clean
- **CLI Tool Included**: Full-featured command-line interface
- **Well Documented**: Comprehensive godoc, README, and examples

## 📦 Installation

```bash
go get github.com/rmrfslashbin/openplantbook-go@v0.1.0
```

## 🚀 Quick Start

```go
import openplantbook "github.com/rmrfslashbin/openplantbook-go"

// Create client
client, err := openplantbook.New(
    openplantbook.WithAPIKey("your-api-key"),
)

// Search for plants
results, err := client.SearchPlants(ctx, "monstera", nil)

// Get plant details
details, err := client.GetPlantDetails(ctx, results[0].PID, nil)
```

## ✨ Features

### Core SDK
- Plant search API with configurable options
- Plant details API with language support
- Intelligent caching (1h for search, 24h for details)
- Client-side rate limiting (200 req/day default)
- Context support for cancellation and timeouts
- Comprehensive error handling with typed errors
- Functional options pattern
- Optional logging interface

### CLI Tool
- Search plants: `openplantbook search monstera`
- Get details: `openplantbook details monstera-deliciosa`
- JSON output for scripting
- Cross-platform binaries (Linux, macOS, Windows)

### Developer Experience
- Zero dependencies (only Go stdlib + golang.org/x)
- 90.5% test coverage
- Comprehensive documentation
- Three working examples
- Race detector clean

## 📚 Documentation

- [README](https://github.com/rmrfslashbin/openplantbook-go/blob/v0.1.0/README.md)
- [GoDoc](https://pkg.go.dev/github.com/rmrfslashbin/openplantbook-go@v0.1.0)
- [Examples](https://github.com/rmrfslashbin/openplantbook-go/tree/v0.1.0/examples)
- [CLI Tool](https://github.com/rmrfslashbin/openplantbook-go/tree/v0.1.0/cmd/openplantbook)

## 💾 Downloads

Binaries are available for:
- Linux (amd64)
- macOS (amd64, arm64) - **Apple Silicon supported!**
- Windows (amd64)

## 🐛 Known Issues

None! This is a stable initial release.

## 📈 What's Next

Future releases may include:
- Redis cache implementation
- Write operations support (when API supports it)
- Batch operations
- MCP server implementation

## 🙏 Acknowledgments

Thanks to the OpenPlantbook team for providing the free API and the community for maintaining the plant database.

## 📄 License

MIT License - see LICENSE for details.

---

**Full Changelog**: https://github.com/rmrfslashbin/openplantbook-go/blob/v0.1.0/CHANGELOG.md

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Public `Version` constant for SDK consumers to reference the library version
- User-Agent header in HTTP requests for API tracking (`openplantbook-go/{version}`)

## [1.1.0] - 2025-11-03

### Added
- Rate limit error behavior with `RateLimitBehavior` type and `WithRateLimitBehavior()` option
- `RateLimitWait` behavior: blocks until rate limiter allows request (default)
- `RateLimitError` behavior: returns error immediately when rate limited
- `ErrRateLimited` error type with `RetryAfter` field
- CLI PID normalization: accepts both hyphenated and space-separated formats

### Fixed
- CLI now accepts PIDs in both formats (e.g., "monstera-deliciosa" or "monstera deliciosa")
- GitHub Actions workflow permissions added for security best practices
- Field alignment in Client struct

## [1.0.1] - 2025-11-03

### Changed
- Upgraded Go version requirement from 1.23 to 1.24
- Updated `golang.org/x/oauth2` from v0.23.0 to v0.27.0 (fixes CVE-2025-22868)

### Fixed
- Windows CI test execution with bash shell
- CI workflow Go version compatibility

## [0.1.0] - 2025-11-02

### Added
- Initial release of OpenPlantbook Go SDK
- Dual authentication support (API Key and OAuth2 Client Credentials)
- Plant search API with configurable options (limit, user plants filter)
- Plant details API with language support
- Intelligent caching system with pluggable interface
  - In-memory cache with TTL support (1h for search, 24h for details)
  - NoOp cache for disabling caching
- Client-side rate limiting (200 requests/day default, configurable)
- Comprehensive error handling with typed errors
  - `ErrUnauthorized` - Invalid credentials
  - `ErrRateLimitExceeded` - Rate limit exceeded
  - `ErrNotFound` - Plant not found
  - `ErrNoAuthProvided` - Missing authentication
  - `ErrMultipleAuthMethods` - Multiple auth methods provided
- Full `context.Context` support for cancellation and timeouts
- Functional options pattern for client configuration
- Optional logging interface for debugging
- Command-line interface (CLI) tool
  - Search command with limit and JSON output
  - Details command with language support
  - Version information with build metadata
- Comprehensive test suite with 90.5% coverage
  - Table-driven tests
  - Mock HTTP servers
  - Race detector clean
- Complete documentation
  - Package-level godoc
  - Detailed README with examples
  - CLI documentation
  - Three working examples (basic search, plant details, caching demo)

### Dependencies
- `golang.org/x/oauth2` v0.32.0 - OAuth2 implementation
- `golang.org/x/time` v0.14.0 - Rate limiting

### Configuration
- Environment variable support for authentication
- Custom HTTP client support
- Configurable base URL
- Configurable rate limits
- Optional debug logging

[Unreleased]: https://github.com/rmrfslashbin/openplantbook-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/rmrfslashbin/openplantbook-go/releases/tag/v0.1.0

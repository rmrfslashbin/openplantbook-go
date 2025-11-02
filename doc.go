// Package openplantbook provides a Go client library for the OpenPlantbook API.
//
// OpenPlantbook is a crowd-sourced database of plant care information, providing
// sensor data and care requirements for thousands of plants.
//
// # Features
//
//   - Dual Authentication: API Key or OAuth2 Client Credentials
//   - Intelligent Caching: In-memory caching with pluggable interface
//   - Rate Limiting: Client-side rate limiting (200 requests/day default)
//   - Context Support: Full context.Context support for cancellation and timeouts
//   - Error Handling: Comprehensive error types with proper wrapping
//
// # Authentication
//
// The SDK supports two authentication methods:
//
// API Key (recommended for read-only operations):
//
//	client, err := openplantbook.New(
//	    openplantbook.WithAPIKey("your-api-key"),
//	)
//
// OAuth2 Client Credentials (for full API access):
//
//	client, err := openplantbook.New(
//	    openplantbook.WithOAuth2("client-id", "client-secret"),
//	)
//
// Get your credentials at: https://open.plantbook.io/
//
// # Basic Usage
//
// Search for plants:
//
//	results, err := client.SearchPlants(ctx, "monstera", &openplantbook.SearchOptions{
//	    Limit: 10,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Get detailed plant information:
//
//	details, err := client.GetPlantDetails(ctx, "monstera-deliciosa", &openplantbook.DetailOptions{
//	    Language: "en",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Configuration
//
// The SDK supports extensive configuration through functional options:
//
//	client, err := openplantbook.New(
//	    openplantbook.WithAPIKey("key"),
//	    openplantbook.WithCache(customCache),
//	    openplantbook.WithRateLimit(100),
//	    openplantbook.WithLogger(logger),
//	)
//
// # Error Handling
//
// The SDK provides typed errors for common scenarios:
//
//	if errors.Is(err, openplantbook.ErrUnauthorized) {
//	    // Handle authentication error
//	}
//	if errors.Is(err, openplantbook.ErrRateLimitExceeded) {
//	    // Handle rate limit
//	}
//
// # Caching
//
// Results are cached automatically:
//   - Search results: 1 hour
//   - Plant details: 24 hours
//
// Cache can be customized or disabled:
//
//	client, err := openplantbook.New(
//	    openplantbook.WithAPIKey("key"),
//	    openplantbook.WithCache(myRedisCache),
//	)
//
// # Rate Limiting
//
// Client-side rate limiting is enabled by default (200 requests/day).
// This can be customized or disabled:
//
//	client, err := openplantbook.New(
//	    openplantbook.WithAPIKey("key"),
//	    openplantbook.WithRateLimit(100), // 100 requests/day
//	)
//
// For more information, see: https://github.com/rmrfslashbin/openplantbook-go
package openplantbook

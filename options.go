package openplantbook

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Option configures the Client
type Option func(*Client) error

// WithAPIKey sets API Key authentication (simpler, read-only endpoints)
// This is the recommended authentication method for v1.0.0 (search and details).
func WithAPIKey(apiKey string) Option {
	return func(c *Client) error {
		if apiKey == "" {
			return ErrInvalidConfig("API key cannot be empty")
		}
		c.apiKey = apiKey
		return nil
	}
}

// WithOAuth2 sets OAuth2 Client Credentials authentication (full API access)
// Required for write operations (sensor data, user plants).
func WithOAuth2(clientID, clientSecret string) Option {
	return func(c *Client) error {
		if clientID == "" || clientSecret == "" {
			return ErrInvalidConfig("client_id and client_secret cannot be empty")
		}
		c.clientID = clientID
		c.clientSecret = clientSecret
		return nil
	}
}

// WithBaseURL sets a custom base URL (useful for testing)
func WithBaseURL(url string) Option {
	return func(c *Client) error {
		if url == "" {
			return ErrInvalidConfig("base URL cannot be empty")
		}
		c.baseURL = url
		return nil
	}
}

// WithHTTPClient allows providing a custom HTTP client
// NOTE: This bypasses authentication configuration
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		if httpClient == nil {
			return ErrInvalidConfig("HTTP client cannot be nil")
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithCache sets a custom cache implementation
func WithCache(cache Cache) Option {
	return func(c *Client) error {
		if cache == nil {
			return ErrInvalidConfig("cache cannot be nil")
		}
		c.cache = cache
		return nil
	}
}

// WithRateLimit sets a custom rate limiter (requests per day)
func WithRateLimit(requestsPerDay int) Option {
	return func(c *Client) error {
		if requestsPerDay <= 0 {
			return ErrInvalidConfig("rate limit must be positive")
		}
		c.rateLimiter = rate.NewLimiter(rate.Every(24*time.Hour/time.Duration(requestsPerDay)), 1)
		return nil
	}
}

// WithLogger injects a custom logger
func WithLogger(logger Logger) Option {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// DisableRateLimit disables client-side rate limiting (use with caution)
func DisableRateLimit() Option {
	return func(c *Client) error {
		c.rateLimiter = nil
		return nil
	}
}

// WithRateLimitBehavior sets how the client handles rate limiting
//
// RateLimitWait (default): Blocks until the rate limiter allows the request
// RateLimitError: Returns a RateLimitError immediately when rate limited
//
// Example:
//
//	client, _ := openplantbook.New(
//	    openplantbook.WithAPIKey(apiKey),
//	    openplantbook.WithRateLimitBehavior(openplantbook.RateLimitError),
//	)
func WithRateLimitBehavior(behavior RateLimitBehavior) Option {
	return func(c *Client) error {
		c.rateLimitBehavior = behavior
		return nil
	}
}

// Logger is the interface for optional logging injection
// Implemented by slog.Logger, logrus, zap, etc.
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

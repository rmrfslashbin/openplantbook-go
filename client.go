package openplantbook

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"
)

const (
	// DefaultBaseURL is the default OpenPlantbook API base URL
	DefaultBaseURL = "https://open.plantbook.io/api/v1"

	// DefaultRateLimit is the default rate limit (200 requests per day)
	DefaultRateLimit = 200
)

// Client represents an OpenPlantbook API client
type Client struct {
	httpClient  *http.Client
	baseURL     string
	rateLimiter *rate.Limiter
	cache       Cache
	logger      Logger

	// Authentication (only ONE should be set)
	apiKey       string
	clientID     string
	clientSecret string
}

// New creates a new OpenPlantbook client with sensible defaults
// Authentication is auto-detected from provided credentials
func New(opts ...Option) (*Client, error) {
	client := &Client{
		baseURL:     DefaultBaseURL,
		rateLimiter: rate.NewLimiter(rate.Every(24*time.Hour/DefaultRateLimit), 1),
		cache:       NewInMemoryCache(),
		logger:      nil, // No logging by default (library pattern)
	}

	// Apply options (sets authentication credentials and other config)
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	// Validate and configure authentication
	if err := client.configureAuth(); err != nil {
		return nil, err
	}

	// Validate client configuration
	if err := client.validate(); err != nil {
		return nil, err
	}

	return client, nil
}

// configureAuth validates auth credentials and configures HTTP client
func (c *Client) configureAuth() error {
	hasAPIKey := c.apiKey != ""
	hasOAuth2 := c.clientID != "" || c.clientSecret != ""

	// If HTTP client already provided, skip auth configuration
	if c.httpClient != nil {
		c.log("using custom HTTP client")
		return nil
	}

	// Validate: exactly ONE auth method must be provided
	if hasAPIKey && hasOAuth2 {
		return ErrMultipleAuthMethods
	}

	if !hasAPIKey && !hasOAuth2 {
		return ErrNoAuthProvided
	}

	// Configure HTTP client based on auth method
	if hasAPIKey {
		// API Key authentication: simple HTTP client with custom transport
		c.httpClient = &http.Client{
			Transport: &apiKeyTransport{
				apiKey:    c.apiKey,
				transport: http.DefaultTransport,
			},
		}
		c.log("using API Key authentication")
	} else {
		// OAuth2 authentication: use official SDK
		if c.clientID == "" || c.clientSecret == "" {
			return ErrInvalidConfig("both client_id and client_secret required for OAuth2")
		}

		oauthConfig := &clientcredentials.Config{
			ClientID:     c.clientID,
			ClientSecret: c.clientSecret,
			TokenURL:     c.baseURL + "/token/",
		}
		c.httpClient = oauthConfig.Client(context.Background())
		c.log("using OAuth2 Client Credentials authentication")
	}

	return nil
}

// validate ensures the client is properly configured
func (c *Client) validate() error {
	if c.baseURL == "" {
		return ErrInvalidConfig("base URL cannot be empty")
	}
	if c.httpClient == nil {
		return ErrInvalidConfig("HTTP client cannot be nil")
	}
	if c.cache == nil {
		return ErrInvalidConfig("cache cannot be nil")
	}
	return nil
}

// log is a helper that only logs if a logger is configured
func (c *Client) log(msg string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Debug(msg, args...)
	}
}

// apiKeyTransport adds API key authentication to requests
type apiKeyTransport struct {
	apiKey    string
	transport http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface
func (t *apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone request to avoid modifying original
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Token "+t.apiKey)
	return t.transport.RoundTrip(req)
}

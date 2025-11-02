package openplantbook

import (
	"errors"
	"net/http"
	"testing"
)

func TestNew_APIKey(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		wantErr error
	}{
		{
			name:    "valid API key",
			apiKey:  "test-api-key",
			wantErr: nil,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			wantErr: ErrInvalidConfig(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var client *Client
			var err error

			if tt.apiKey != "" {
				client, err = New(WithAPIKey(tt.apiKey))
			} else {
				client, err = New(WithAPIKey(tt.apiKey))
			}

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("New() expected error, got nil")
				}
				// Check error type
				var cfgErr *ConfigError
				if !errors.As(err, &cfgErr) {
					t.Errorf("New() error type = %T, want *ConfigError", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("New() unexpected error: %v", err)
			}

			if client == nil {
				t.Fatal("New() returned nil client")
			}

			// Verify auth configured
			if client.apiKey != tt.apiKey {
				t.Errorf("client.apiKey = %q, want %q", client.apiKey, tt.apiKey)
			}

			if client.httpClient == nil {
				t.Error("client.httpClient is nil")
			}

			// Verify defaults
			if client.baseURL != DefaultBaseURL {
				t.Errorf("client.baseURL = %q, want %q", client.baseURL, DefaultBaseURL)
			}

			if client.cache == nil {
				t.Error("client.cache is nil")
			}

			if client.rateLimiter == nil {
				t.Error("client.rateLimiter is nil")
			}
		})
	}
}

func TestNew_OAuth2(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		wantErr      error
	}{
		{
			name:         "valid OAuth2 credentials",
			clientID:     "test-client-id",
			clientSecret: "test-client-secret",
			wantErr:      nil,
		},
		{
			name:         "missing client secret",
			clientID:     "test-client-id",
			clientSecret: "",
			wantErr:      ErrInvalidConfig(""),
		},
		{
			name:         "missing client ID",
			clientID:     "",
			clientSecret: "test-client-secret",
			wantErr:      ErrInvalidConfig(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(WithOAuth2(tt.clientID, tt.clientSecret))

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("New() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("New() unexpected error: %v", err)
			}

			if client == nil {
				t.Fatal("New() returned nil client")
			}

			// Verify auth configured
			if client.clientID != tt.clientID {
				t.Errorf("client.clientID = %q, want %q", client.clientID, tt.clientID)
			}

			if client.clientSecret != tt.clientSecret {
				t.Errorf("client.clientSecret = %q, want %q", client.clientSecret, tt.clientSecret)
			}

			if client.httpClient == nil {
				t.Error("client.httpClient is nil")
			}
		})
	}
}

func TestNew_NoAuth(t *testing.T) {
	_, err := New()

	if err == nil {
		t.Fatal("New() with no auth expected error, got nil")
	}

	if !errors.Is(err, ErrNoAuthProvided) {
		t.Errorf("New() error = %v, want %v", err, ErrNoAuthProvided)
	}
}

func TestNew_MultipleAuth(t *testing.T) {
	_, err := New(
		WithAPIKey("test-api-key"),
		WithOAuth2("client-id", "client-secret"),
	)

	if err == nil {
		t.Fatal("New() with multiple auth expected error, got nil")
	}

	if !errors.Is(err, ErrMultipleAuthMethods) {
		t.Errorf("New() error = %v, want %v", err, ErrMultipleAuthMethods)
	}
}

func TestNew_WithOptions(t *testing.T) {
	customBaseURL := "https://custom.example.com"
	customRateLimit := 100

	client, err := New(
		WithAPIKey("test-api-key"),
		WithBaseURL(customBaseURL),
		WithRateLimit(customRateLimit),
		WithCache(NewNoOpCache()),
	)

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if client.baseURL != customBaseURL {
		t.Errorf("client.baseURL = %q, want %q", client.baseURL, customBaseURL)
	}

	// Verify NoOpCache
	_, ok := client.cache.(*NoOpCache)
	if !ok {
		t.Errorf("client.cache type = %T, want *NoOpCache", client.cache)
	}
}

func TestNew_WithCustomHTTPClient(t *testing.T) {
	customClient := &http.Client{}

	client, err := New(
		WithHTTPClient(customClient),
	)

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if client.httpClient != customClient {
		t.Error("client.httpClient is not the custom client")
	}
}

func TestNew_DisableRateLimit(t *testing.T) {
	client, err := New(
		WithAPIKey("test-api-key"),
		DisableRateLimit(),
	)

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if client.rateLimiter != nil {
		t.Error("client.rateLimiter should be nil when disabled")
	}
}

func TestAPIKeyTransport_RoundTrip(t *testing.T) {
	apiKey := "test-api-key"
	transport := &apiKeyTransport{
		apiKey:    apiKey,
		transport: http.DefaultTransport,
	}

	// Verify the transport is configured
	if transport.apiKey != apiKey {
		t.Errorf("transport.apiKey = %q, want %q", transport.apiKey, apiKey)
	}

	if transport.transport == nil {
		t.Error("transport.transport is nil")
	}
}

// mockLogger implements the Logger interface for testing
type mockLogger struct {
	debugCalls int
	infoCalls  int
	warnCalls  int
	errorCalls int
}

func (m *mockLogger) Debug(msg string, args ...interface{}) {
	m.debugCalls++
}

func (m *mockLogger) Info(msg string, args ...interface{}) {
	m.infoCalls++
}

func (m *mockLogger) Warn(msg string, args ...interface{}) {
	m.warnCalls++
}

func (m *mockLogger) Error(msg string, args ...interface{}) {
	m.errorCalls++
}

func TestNew_WithLogger(t *testing.T) {
	logger := &mockLogger{}

	client, err := New(
		WithAPIKey("test-api-key"),
		WithLogger(logger),
	)

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if client.logger != logger {
		t.Error("client.logger is not the provided logger")
	}

	// Verify logging was called during auth setup
	if logger.debugCalls == 0 {
		t.Error("logger.Debug was not called during client creation")
	}
}

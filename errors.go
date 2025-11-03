package openplantbook

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Common sentinel errors
var (
	// Authentication errors
	ErrUnauthorized        = errors.New("invalid credentials or token expired")
	ErrMultipleAuthMethods = errors.New("multiple authentication methods provided (use only API key OR OAuth2)")
	ErrNoAuthProvided      = errors.New("no authentication provided (use WithAPIKey or WithOAuth2)")

	// API errors
	ErrRateLimitExceeded = errors.New("rate limit exceeded (200 requests/day)")
	ErrNotFound          = errors.New("plant not found")

	// Input validation
	ErrInvalidInput = func(msg string) error { return &ValidationError{Message: msg} }

	// Configuration errors
	ErrInvalidConfig = func(msg string) error { return &ConfigError{Message: msg} }
)

// APIError represents an error response from the OpenPlantbook API
type APIError struct {
	StatusCode int
	Message    string
	Endpoint   string
	Body       string
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d) at %s: %s", e.StatusCode, e.Endpoint, e.Message)
}

// IsClientError returns true if the error is a 4xx client error
func (e *APIError) IsClientError() bool {
	return e.StatusCode >= 400 && e.StatusCode < 500
}

// IsServerError returns true if the error is a 5xx server error
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500
}

// ValidationError represents invalid input parameters
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation failed for %s='%v': %s", e.Field, e.Value, e.Message)
	}
	return fmt.Sprintf("validation failed: %s", e.Message)
}

// ConfigError represents a configuration error
type ConfigError struct {
	Message string
}

// Error implements the error interface
func (e *ConfigError) Error() string {
	return fmt.Sprintf("configuration error: %s", e.Message)
}

// ErrRateLimited indicates the rate limit has been exceeded
// This error is returned when RateLimitBehavior is set to RateLimitError
// and a request would exceed the configured rate limit.
type ErrRateLimited struct {
	RetryAfter time.Time // When the next request can be made
	Message    string
}

// Error implements the error interface
func (e *ErrRateLimited) Error() string {
	return fmt.Sprintf("rate limit exceeded: %s (retry after %s)",
		e.Message,
		e.RetryAfter.Format(time.RFC3339))
}

// newAPIError creates an APIError from an HTTP response
func newAPIError(resp *http.Response, endpoint string) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Endpoint:   endpoint,
	}

	// Parse common error cases
	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		apiErr.Message = "authentication failed"
		return fmt.Errorf("%w: %s", ErrUnauthorized, apiErr.Message)
	case http.StatusNotFound:
		apiErr.Message = "resource not found"
		return fmt.Errorf("%w: %s", ErrNotFound, apiErr.Message)
	case http.StatusTooManyRequests:
		apiErr.Message = "rate limit exceeded"
		return fmt.Errorf("%w: %s", ErrRateLimitExceeded, apiErr.Message)
	default:
		apiErr.Message = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return apiErr
	}
}

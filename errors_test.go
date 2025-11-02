package openplantbook

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name    string
		apiErr  *APIError
		wantMsg string
	}{
		{
			name: "basic API error",
			apiErr: &APIError{
				StatusCode: 500,
				Message:    "internal server error",
				Endpoint:   "/plant/search/",
			},
			wantMsg: "API error (status 500) at /plant/search/: internal server error",
		},
		{
			name: "404 error",
			apiErr: &APIError{
				StatusCode: 404,
				Message:    "not found",
				Endpoint:   "/plant/detail/unknown/",
			},
			wantMsg: "API error (status 404) at /plant/detail/unknown/: not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.apiErr.Error()
			if got != tt.wantMsg {
				t.Errorf("APIError.Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestAPIError_IsClientError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"400 is client error", 400, true},
		{"404 is client error", 404, true},
		{"499 is client error", 499, true},
		{"500 is not client error", 500, false},
		{"200 is not client error", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := &APIError{StatusCode: tt.statusCode}
			got := apiErr.IsClientError()
			if got != tt.want {
				t.Errorf("APIError.IsClientError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIError_IsServerError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{"500 is server error", 500, true},
		{"502 is server error", 502, true},
		{"599 is server error", 599, true},
		{"400 is not server error", 400, false},
		{"200 is not server error", 200, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := &APIError{StatusCode: tt.statusCode}
			got := apiErr.IsServerError()
			if got != tt.want {
				t.Errorf("APIError.IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name    string
		valErr  *ValidationError
		wantMsg string
	}{
		{
			name: "with field",
			valErr: &ValidationError{
				Field:   "query",
				Value:   "",
				Message: "cannot be empty",
			},
			wantMsg: "validation failed for query='': cannot be empty",
		},
		{
			name: "without field",
			valErr: &ValidationError{
				Message: "invalid input",
			},
			wantMsg: "validation failed: invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.valErr.Error()
			if got != tt.wantMsg {
				t.Errorf("ValidationError.Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestConfigError_Error(t *testing.T) {
	cfgErr := &ConfigError{Message: "invalid configuration"}
	want := "configuration error: invalid configuration"
	got := cfgErr.Error()

	if got != want {
		t.Errorf("ConfigError.Error() = %q, want %q", got, want)
	}
}

func TestNewAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		endpoint   string
		wantErr    error
	}{
		{
			name:       "401 unauthorized",
			statusCode: http.StatusUnauthorized,
			endpoint:   "/plant/search/",
			wantErr:    ErrUnauthorized,
		},
		{
			name:       "403 forbidden",
			statusCode: http.StatusForbidden,
			endpoint:   "/plant/search/",
			wantErr:    ErrUnauthorized,
		},
		{
			name:       "404 not found",
			statusCode: http.StatusNotFound,
			endpoint:   "/plant/detail/test/",
			wantErr:    ErrNotFound,
		},
		{
			name:       "429 rate limit",
			statusCode: http.StatusTooManyRequests,
			endpoint:   "/plant/search/",
			wantErr:    ErrRateLimitExceeded,
		},
		{
			name:       "500 server error",
			statusCode: http.StatusInternalServerError,
			endpoint:   "/plant/search/",
			wantErr:    nil, // Returns *APIError, not sentinel
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock response
			rec := httptest.NewRecorder()
			rec.WriteHeader(tt.statusCode)
			resp := rec.Result()
			defer resp.Body.Close()

			err := newAPIError(resp, tt.endpoint)
			if err == nil {
				t.Fatal("newAPIError() returned nil error")
			}

			if tt.wantErr != nil {
				// Check if error wraps the expected sentinel
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("newAPIError() error = %v, want %v", err, tt.wantErr)
				}
			} else {
				// Should be an *APIError
				var apiErr *APIError
				if !errors.As(err, &apiErr) {
					t.Errorf("newAPIError() error type = %T, want *APIError", err)
				}
			}
		})
	}
}

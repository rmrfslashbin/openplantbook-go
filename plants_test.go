package openplantbook

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestClient_SearchPlants(t *testing.T) {
	// Load test fixture
	searchData, err := os.ReadFile("testdata/search_response.json")
	if err != nil {
		t.Fatalf("failed to load test fixture: %v", err)
	}

	tests := []struct {
		name         string
		query        string
		opts         *SearchOptions
		mockStatus   int
		mockResponse string
		wantResults  int
		wantErr      bool
	}{
		{
			name:         "successful search",
			query:        "monstera",
			opts:         &SearchOptions{Limit: 10},
			mockStatus:   http.StatusOK,
			mockResponse: string(searchData),
			wantResults:  2,
			wantErr:      false,
		},
		{
			name:         "empty query",
			query:        "",
			opts:         nil,
			mockStatus:   0,
			mockResponse: "",
			wantResults:  0,
			wantErr:      true,
		},
		{
			name:         "API 404 error",
			query:        "unknown",
			opts:         nil,
			mockStatus:   http.StatusNotFound,
			mockResponse: `{"error":"not found"}`,
			wantResults:  0,
			wantErr:      true,
		},
		{
			name:         "API 429 rate limit",
			query:        "test",
			opts:         nil,
			mockStatus:   http.StatusTooManyRequests,
			mockResponse: `{"error":"rate limit"}`,
			wantResults:  0,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "GET" {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				if tt.query != "" {
					queryParam := r.URL.Query().Get("alias")
					if queryParam != tt.query {
						t.Errorf("expected query=%s, got %s", tt.query, queryParam)
					}

					if tt.opts != nil && tt.opts.Limit > 0 {
						limitParam := r.URL.Query().Get("limit")
						if limitParam == "" {
							t.Error("expected limit parameter")
						}
					}
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create client
			client, err := New(
				WithAPIKey("test-key"),
				WithBaseURL(server.URL),
				DisableRateLimit(),
			)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			// Execute search
			results, err := client.SearchPlants(context.Background(), tt.query, tt.opts)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchPlants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Check results
			if len(results) != tt.wantResults {
				t.Errorf("SearchPlants() got %d results, want %d", len(results), tt.wantResults)
			}

			// Verify result structure
			if len(results) > 0 {
				first := results[0]
				if first.PID == "" {
					t.Error("result PID is empty")
				}
				if first.Alias == "" {
					t.Error("result Alias is empty")
				}
			}
		})
	}
}

func TestClient_SearchPlants_Caching(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"pid":"test","display_pid":"Test","alias":"Test Plant","category":"Test"}]`))
	}))
	defer server.Close()

	client, err := New(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		DisableRateLimit(),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// First call - should hit API
	_, err = client.SearchPlants(context.Background(), "test", nil)
	if err != nil {
		t.Fatalf("first SearchPlants() failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 API call, got %d", callCount)
	}

	// Second call - should use cache
	_, err = client.SearchPlants(context.Background(), "test", nil)
	if err != nil {
		t.Fatalf("second SearchPlants() failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

func TestClient_GetPlantDetails(t *testing.T) {
	// Load test fixture
	detailData, err := os.ReadFile("testdata/detail_response.json")
	if err != nil {
		t.Fatalf("failed to load test fixture: %v", err)
	}

	tests := []struct {
		name         string
		pid          string
		opts         *DetailOptions
		mockStatus   int
		mockResponse string
		wantErr      bool
	}{
		{
			name:         "successful detail retrieval",
			pid:          "monstera-deliciosa",
			opts:         &DetailOptions{Language: "en"},
			mockStatus:   http.StatusOK,
			mockResponse: string(detailData),
			wantErr:      false,
		},
		{
			name:         "empty pid",
			pid:          "",
			opts:         nil,
			mockStatus:   0,
			mockResponse: "",
			wantErr:      true,
		},
		{
			name:         "API 404 error",
			pid:          "unknown",
			opts:         nil,
			mockStatus:   http.StatusNotFound,
			mockResponse: `{"error":"not found"}`,
			wantErr:      true,
		},
		{
			name:         "API 401 unauthorized",
			pid:          "test",
			opts:         nil,
			mockStatus:   http.StatusUnauthorized,
			mockResponse: `{"error":"unauthorized"}`,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "GET" {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				if tt.pid != "" {
					expectedPath := "/plant/detail/" + tt.pid + "/"
					if r.URL.Path != expectedPath {
						t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
					}

					if tt.opts != nil && tt.opts.Language != "" {
						langParam := r.URL.Query().Get("lang")
						if langParam != tt.opts.Language {
							t.Errorf("expected lang=%s, got %s", tt.opts.Language, langParam)
						}
					}
				}

				w.WriteHeader(tt.mockStatus)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			// Create client
			client, err := New(
				WithAPIKey("test-key"),
				WithBaseURL(server.URL),
				DisableRateLimit(),
			)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			// Execute get details
			details, err := client.GetPlantDetails(context.Background(), tt.pid, tt.opts)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlantDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify details
			if details.PID == "" {
				t.Error("details PID is empty")
			}
			if details.Alias == "" {
				t.Error("details Alias is empty")
			}
			if details.MaxTemp == 0 {
				t.Error("details MaxTemp is zero")
			}
		})
	}
}

func TestClient_GetPlantDetails_Caching(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"pid":"test","display_pid":"Test","alias":"Test Plant","max_temp":25.0,"min_temp":15.0,"category":"Test"}`))
	}))
	defer server.Close()

	client, err := New(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		DisableRateLimit(),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// First call - should hit API
	_, err = client.GetPlantDetails(context.Background(), "test", nil)
	if err != nil {
		t.Fatalf("first GetPlantDetails() failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 API call, got %d", callCount)
	}

	// Second call - should use cache
	_, err = client.GetPlantDetails(context.Background(), "test", nil)
	if err != nil {
		t.Fatalf("second GetPlantDetails() failed: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 API call (cached), got %d", callCount)
	}
}

func TestClient_RateLimiting(t *testing.T) {
	// Skip this test in short mode as it involves timing delays
	if testing.Short() {
		t.Skip("skipping rate limiting test in short mode")
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	// Create client with custom rate limiter for faster testing
	// We'll manually set a very restrictive rate for testing
	client, err := New(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithCache(NewNoOpCache()), // Disable cache to test rate limiting
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Manually override with very restrictive limiter for testing
	// 1 request per 100ms
	client.rateLimiter = rate.NewLimiter(rate.Every(100*time.Millisecond), 1)

	// First request should succeed immediately
	_, err = client.SearchPlants(context.Background(), "test1", nil)
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}

	// Second request should be delayed by ~100ms
	start := time.Now()
	_, err = client.SearchPlants(context.Background(), "test2", nil)
	if err != nil {
		t.Fatalf("second request failed: %v", err)
	}
	elapsed := time.Since(start)

	// Should have waited ~100ms
	if elapsed < 50*time.Millisecond {
		t.Errorf("expected rate limiting delay of ~100ms, got %v", elapsed)
	}

	t.Logf("rate limiting delay: %v", elapsed)
}

func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	client, err := New(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		DisableRateLimit(),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Request should fail due to context cancellation
	_, err = client.SearchPlants(ctx, "test", nil)
	if err == nil {
		t.Error("expected context cancellation error, got nil")
	}
}

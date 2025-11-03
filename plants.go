package openplantbook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// SearchPlants searches for plants by alias/common name
func (c *Client) SearchPlants(ctx context.Context, query string, opts *SearchOptions) ([]PlantSearchResult, error) {
	if query == "" {
		return nil, ErrInvalidInput("query cannot be empty")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("search:%s:%v", query, opts)
	if cached, ok := c.cache.Get(cacheKey); ok {
		var results []PlantSearchResult
		if err := json.Unmarshal(cached, &results); err == nil {
			c.log("cache hit for search", "query", query)
			return results, nil
		}
	}

	// Handle rate limiting based on configured behavior
	if c.rateLimiter != nil {
		if c.rateLimitBehavior == RateLimitError {
			// Check if we can proceed without waiting
			reservation := c.rateLimiter.Reserve()
			if !reservation.OK() {
				return nil, &ErrRateLimited{
					RetryAfter: time.Now().Add(24 * time.Hour),
					Message:    "rate limiter exhausted",
				}
			}

			delay := reservation.Delay()
			if delay > 0 {
				// Cancel the reservation and return error
				reservation.Cancel()
				return nil, &ErrRateLimited{
					RetryAfter: time.Now().Add(delay),
					Message:    "rate limit exceeded, please retry later",
				}
			}
			// If delay is 0, reservation is consumed and we can proceed
		} else {
			// Default behavior: wait for rate limiter
			if err := c.rateLimiter.Wait(ctx); err != nil {
				return nil, fmt.Errorf("rate limit wait: %w", err)
			}
		}
	}

	// Build request
	req, err := c.newRequest(ctx, "GET", "/plant/search", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	q.Set("alias", query)

	if opts != nil {
		if opts.Limit > 0 {
			q.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.UserPlants {
			q.Set("userplant", "user")
		}
	}
	req.URL.RawQuery = q.Encode()

	// Execute request
	var response searchResponse
	if err := c.doRequest(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("search plants: %w", err)
	}

	c.log("search completed", "query", query, "results", len(response.Results))

	// Cache results (1 hour TTL)
	if data, err := json.Marshal(response.Results); err == nil {
		c.cache.Set(cacheKey, data, 1*time.Hour)
	}

	return response.Results, nil
}

// GetPlantDetails retrieves detailed plant care information
func (c *Client) GetPlantDetails(ctx context.Context, pid string, opts *DetailOptions) (*PlantDetails, error) {
	if pid == "" {
		return nil, ErrInvalidInput("pid cannot be empty")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("detail:%s:%v", pid, opts)
	if cached, ok := c.cache.Get(cacheKey); ok {
		var details PlantDetails
		if err := json.Unmarshal(cached, &details); err == nil {
			c.log("cache hit for details", "pid", pid)
			return &details, nil
		}
	}

	// Handle rate limiting based on configured behavior
	if c.rateLimiter != nil {
		if c.rateLimitBehavior == RateLimitError {
			// Check if we can proceed without waiting
			reservation := c.rateLimiter.Reserve()
			if !reservation.OK() {
				return nil, &ErrRateLimited{
					RetryAfter: time.Now().Add(24 * time.Hour),
					Message:    "rate limiter exhausted",
				}
			}

			delay := reservation.Delay()
			if delay > 0 {
				// Cancel the reservation and return error
				reservation.Cancel()
				return nil, &ErrRateLimited{
					RetryAfter: time.Now().Add(delay),
					Message:    "rate limit exceeded, please retry later",
				}
			}
			// If delay is 0, reservation is consumed and we can proceed
		} else {
			// Default behavior: wait for rate limiter
			if err := c.rateLimiter.Wait(ctx); err != nil {
				return nil, fmt.Errorf("rate limit wait: %w", err)
			}
		}
	}

	// Build request
	path := fmt.Sprintf("/plant/detail/%s", pid)
	req, err := c.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Add query parameters
	if opts != nil && opts.Language != "" {
		q := req.URL.Query()
		q.Set("lang", opts.Language)
		req.URL.RawQuery = q.Encode()
	}

	// Execute request
	var details PlantDetails
	if err := c.doRequest(ctx, req, &details); err != nil {
		return nil, fmt.Errorf("get plant details: %w", err)
	}

	c.log("details retrieved", "pid", pid)

	// Cache results (24 hours TTL)
	if data, err := json.Marshal(details); err == nil {
		c.cache.Set(cacheKey, data, 24*time.Hour)
	}

	return &details, nil
}

// newRequest creates a new HTTP request with the base URL
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "openplantbook-go/"+Version)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// doRequest executes an HTTP request and decodes the JSON response
func (c *Client) doRequest(ctx context.Context, req *http.Request, result interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return newAPIError(resp, req.URL.Path)
	}

	// Decode JSON response
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

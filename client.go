package openplantbook

import (
	"net/http"

	"golang.org/x/time/rate"
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

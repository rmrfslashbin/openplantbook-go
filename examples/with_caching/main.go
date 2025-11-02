package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	openplantbook "github.com/rmrfslashbin/openplantbook-go"
)

func main() {
	// This example demonstrates caching behavior with custom rate limiting.
	// It makes multiple identical requests to show cache hits.

	apiKey := os.Getenv("OPENPLANTBOOK_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENPLANTBOOK_API_KEY environment variable is required")
	}

	// Create client with custom cache and rate limiting
	cache := openplantbook.NewInMemoryCache()
	client, err := openplantbook.New(
		openplantbook.WithAPIKey(apiKey),
		openplantbook.WithCache(cache),
		openplantbook.WithRateLimit(100), // 100 requests per day
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	query := "fern"
	fmt.Printf("Demonstrating cache behavior with query '%s'\n\n", query)

	// First request - will hit the API
	fmt.Println("=== First Request (API call) ===")
	start := time.Now()
	results1, err := client.SearchPlants(context.Background(), query, nil)
	if err != nil {
		log.Fatalf("First search failed: %v", err)
	}
	elapsed1 := time.Since(start)
	fmt.Printf("Found %d plants in %v\n\n", len(results1), elapsed1)

	// Second request - should use cache (much faster)
	fmt.Println("=== Second Request (cached) ===")
	start = time.Now()
	results2, err := client.SearchPlants(context.Background(), query, nil)
	if err != nil {
		log.Fatalf("Second search failed: %v", err)
	}
	elapsed2 := time.Since(start)
	fmt.Printf("Found %d plants in %v\n", len(results2), elapsed2)
	fmt.Printf("Speedup: %.2fx faster\n\n", float64(elapsed1)/float64(elapsed2))

	// Third request - still cached
	fmt.Println("=== Third Request (still cached) ===")
	start = time.Now()
	results3, err := client.SearchPlants(context.Background(), query, nil)
	if err != nil {
		log.Fatalf("Third search failed: %v", err)
	}
	elapsed3 := time.Since(start)
	fmt.Printf("Found %d plants in %v\n\n", len(results3), elapsed3)

	// Demonstrate cache clearing
	fmt.Println("=== Clearing Cache ===")
	cache.Clear()
	fmt.Println("Cache cleared successfully\n")

	// Fourth request - will hit the API again
	fmt.Println("=== Fourth Request (API call after clear) ===")
	start = time.Now()
	results4, err := client.SearchPlants(context.Background(), query, nil)
	if err != nil {
		log.Fatalf("Fourth search failed: %v", err)
	}
	elapsed4 := time.Since(start)
	fmt.Printf("Found %d plants in %v\n\n", len(results4), elapsed4)

	// Summary
	fmt.Println("=== Cache Behavior Summary ===")
	fmt.Printf("Request 1 (API):    %v\n", elapsed1)
	fmt.Printf("Request 2 (cache):  %v\n", elapsed2)
	fmt.Printf("Request 3 (cache):  %v\n", elapsed3)
	fmt.Printf("Request 4 (API):    %v\n", elapsed4)
	fmt.Println("\nNotice how cached requests are significantly faster!")
}

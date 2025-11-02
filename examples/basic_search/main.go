package main

import (
	"context"
	"fmt"
	"log"
	"os"

	openplantbook "github.com/rmrfslashbin/openplantbook-go"
)

func main() {
	// This example demonstrates basic plant search functionality
	// using API Key authentication.

	apiKey := os.Getenv("OPENPLANTBOOK_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENPLANTBOOK_API_KEY environment variable is required")
	}

	// Create client with API Key
	client, err := openplantbook.New(
		openplantbook.WithAPIKey(apiKey),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Search for plants
	query := "monstera"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	fmt.Printf("Searching for plants matching '%s'...\n\n", query)

	results, err := client.SearchPlants(context.Background(), query, &openplantbook.SearchOptions{
		Limit: 10,
	})
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	// Display results
	fmt.Printf("Found %d plant(s):\n\n", len(results))
	for i, plant := range results {
		fmt.Printf("%d. %s\n", i+1, plant.DisplayPID)
		fmt.Printf("   Alias: %s\n", plant.Alias)
		fmt.Printf("   PID: %s\n", plant.PID)
		fmt.Printf("   Category: %s\n", plant.Category)
		fmt.Println()
	}
}

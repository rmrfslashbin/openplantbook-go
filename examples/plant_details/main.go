package main

import (
	"context"
	"fmt"
	"log"
	"os"

	openplantbook "github.com/rmrfslashbin/openplantbook-go"
)

func main() {
	// This example demonstrates retrieving detailed plant care information
	// using OAuth2 authentication.

	clientID := os.Getenv("OPENPLANTBOOK_CLIENT_ID")
	clientSecret := os.Getenv("OPENPLANTBOOK_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("OPENPLANTBOOK_CLIENT_ID and OPENPLANTBOOK_CLIENT_SECRET environment variables are required")
	}

	// Create client with OAuth2
	client, err := openplantbook.New(
		openplantbook.WithOAuth2(clientID, clientSecret),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Get plant PID from command line or use default
	pid := "monstera-deliciosa"
	if len(os.Args) > 1 {
		pid = os.Args[1]
	}

	fmt.Printf("Fetching details for plant '%s'...\n\n", pid)

	// Get plant details with English language
	details, err := client.GetPlantDetails(context.Background(), pid, &openplantbook.DetailOptions{
		Language: "en",
	})
	if err != nil {
		log.Fatalf("Failed to get plant details: %v", err)
	}

	// Display detailed care information
	fmt.Printf("Plant: %s\n", details.DisplayPID)
	fmt.Printf("Common Name: %s\n", details.Alias)
	fmt.Printf("Category: %s\n\n", details.Category)

	fmt.Println("Care Requirements:")
	fmt.Println("==================")
	fmt.Printf("Light (Lux): %d - %d\n", details.MinLightLux, details.MaxLightLux)
	fmt.Printf("Temperature (°C): %.1f - %.1f\n", details.MinTemp, details.MaxTemp)
	fmt.Printf("Humidity (%%): %d - %d\n", details.MinEnvHumid, details.MaxEnvHumid)
	fmt.Printf("Soil Moisture (%%): %d - %d\n", details.MinSoilMoist, details.MaxSoilMoist)
	fmt.Printf("Soil EC (μS/cm): %d - %d\n", details.MinSoilEC, details.MaxSoilEC)

	if details.ImageURL != "" {
		fmt.Printf("\nImage: %s\n", details.ImageURL)
	}
}

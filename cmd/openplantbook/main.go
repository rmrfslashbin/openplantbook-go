package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"text/tabwriter"

	openplantbook "github.com/rmrfslashbin/openplantbook-go"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

type Config struct {
	APIKey       string
	ClientID     string
	ClientSecret string
	BaseURL      string
	Debug        bool
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version", "-v", "--version":
		printVersion()
	case "search":
		runSearch()
	case "details":
		runDetails()
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func runSearch() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: openplantbook search <query> [--limit N] [--json]")
		os.Exit(1)
	}

	query := os.Args[2]
	limit := 10
	jsonOutput := false

	// Parse flags
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--limit":
			if i+1 < len(os.Args) {
				fmt.Sscanf(os.Args[i+1], "%d", &limit)
				i++
			}
		case "--json":
			jsonOutput = true
		}
	}

	config := loadConfig()
	client, err := createClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	results, err := client.SearchPlants(context.Background(), query, &openplantbook.SearchOptions{
		Limit: limit,
	})
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	if jsonOutput {
		outputJSON(results)
	} else {
		outputSearchResults(results)
	}
}

func runDetails() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: openplantbook details <pid> [--lang LANG] [--json]")
		os.Exit(1)
	}

	pid := os.Args[2]
	lang := "en"
	jsonOutput := false

	// Parse flags
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--lang":
			if i+1 < len(os.Args) {
				lang = os.Args[i+1]
				i++
			}
		case "--json":
			jsonOutput = true
		}
	}

	config := loadConfig()
	client, err := createClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	details, err := client.GetPlantDetails(context.Background(), pid, &openplantbook.DetailOptions{
		Language: lang,
	})
	if err != nil {
		log.Fatalf("Failed to get details: %v", err)
	}

	if jsonOutput {
		outputJSON(details)
	} else {
		outputPlantDetails(details)
	}
}

func loadConfig() Config {
	return Config{
		APIKey:       os.Getenv("OPENPLANTBOOK_API_KEY"),
		ClientID:     os.Getenv("OPENPLANTBOOK_CLIENT_ID"),
		ClientSecret: os.Getenv("OPENPLANTBOOK_CLIENT_SECRET"),
		BaseURL:      os.Getenv("OPENPLANTBOOK_BASE_URL"),
		Debug:        os.Getenv("OPENPLANTBOOK_DEBUG") == "true",
	}
}

func createClient(config Config) (*openplantbook.Client, error) {
	opts := []openplantbook.Option{}

	// Authentication
	if config.APIKey != "" {
		opts = append(opts, openplantbook.WithAPIKey(config.APIKey))
	} else if config.ClientID != "" && config.ClientSecret != "" {
		opts = append(opts, openplantbook.WithOAuth2(config.ClientID, config.ClientSecret))
	} else {
		return nil, fmt.Errorf("no authentication provided: set OPENPLANTBOOK_API_KEY or OPENPLANTBOOK_CLIENT_ID/CLIENT_SECRET")
	}

	// Optional base URL override
	if config.BaseURL != "" {
		opts = append(opts, openplantbook.WithBaseURL(config.BaseURL))
	}

	// Debug logging
	if config.Debug {
		logger := &cliLogger{slog.New(slog.NewTextHandler(os.Stderr, nil))}
		opts = append(opts, openplantbook.WithLogger(logger))
	}

	return openplantbook.New(opts...)
}

func outputSearchResults(results []openplantbook.PlantSearchResult) {
	if len(results) == 0 {
		fmt.Println("No plants found")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SCIENTIFIC NAME\tCOMMON NAME\tPID\tCATEGORY")
	fmt.Fprintln(w, "---------------\t-----------\t---\t--------")
	for _, plant := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", plant.DisplayPID, plant.Alias, plant.PID, plant.Category)
	}
	w.Flush()
	fmt.Printf("\nFound %d plant(s)\n", len(results))
}

func outputPlantDetails(details *openplantbook.PlantDetails) {
	fmt.Printf("Plant: %s\n", details.DisplayPID)
	fmt.Printf("Common Name: %s\n", details.Alias)
	fmt.Printf("PID: %s\n", details.PID)
	fmt.Printf("Category: %s\n\n", details.Category)

	fmt.Println("Care Requirements:")
	fmt.Println("==================")
	fmt.Printf("Light (Lux):       %d - %d\n", details.MinLightLux, details.MaxLightLux)
	fmt.Printf("Temperature (°C):  %.1f - %.1f\n", details.MinTemp, details.MaxTemp)
	fmt.Printf("Humidity (%%):      %d - %d\n", details.MinEnvHumid, details.MaxEnvHumid)
	fmt.Printf("Soil Moisture (%%): %d - %d\n", details.MinSoilMoist, details.MaxSoilMoist)
	fmt.Printf("Soil EC (μS/cm):   %d - %d\n", details.MinSoilEC, details.MaxSoilEC)

	if details.ImageURL != "" {
		fmt.Printf("\nImage: %s\n", details.ImageURL)
	}
}

func outputJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		log.Fatalf("Failed to encode JSON: %v", err)
	}
}

func printVersion() {
	fmt.Printf("openplantbook version %s\n", version)
	fmt.Printf("  commit: %s\n", commit)
	fmt.Printf("  built:  %s\n", date)
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: openplantbook <command> [options]

Commands:
  search <query>    Search for plants by name
  details <pid>     Get detailed plant information
  version           Show version information
  help              Show this help message

Run 'openplantbook <command> --help' for command-specific options.
`)
}

func printHelp() {
	fmt.Printf(`OpenPlantbook CLI - Plant care information from the command line

USAGE:
  openplantbook <command> [options]

COMMANDS:
  search <query>         Search for plants by name or alias
    --limit N            Limit results to N plants (default: 10)
    --json               Output results as JSON

  details <pid>          Get detailed care information for a plant
    --lang LANG          Specify language (ISO 639-1 code, default: en)
    --json               Output results as JSON

  version                Show version information
  help                   Show this help message

AUTHENTICATION:
  API Key (recommended for simple use):
    export OPENPLANTBOOK_API_KEY="your-api-key"

  OAuth2 (for full API access):
    export OPENPLANTBOOK_CLIENT_ID="your-client-id"
    export OPENPLANTBOOK_CLIENT_SECRET="your-client-secret"

CONFIGURATION:
  OPENPLANTBOOK_BASE_URL    Override API base URL
  OPENPLANTBOOK_DEBUG=true  Enable debug logging

EXAMPLES:
  # Search for plants
  openplantbook search monstera

  # Search with limit
  openplantbook search fern --limit 5

  # Get plant details
  openplantbook details monstera-deliciosa

  # Get details in Spanish
  openplantbook details monstera-deliciosa --lang es

  # Output as JSON for scripting
  openplantbook search monstera --json | jq '.[] | .pid'

For more information, visit: https://github.com/rmrfslashbin/openplantbook-go
`)
}

// cliLogger implements the openplantbook.Logger interface
type cliLogger struct {
	logger *slog.Logger
}

func (l *cliLogger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

func (l *cliLogger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

func (l *cliLogger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

func (l *cliLogger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

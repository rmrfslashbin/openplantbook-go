package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	openplantbook "github.com/rmrfslashbin/openplantbook-go"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"

	cfgFile string
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "openplantbook",
		Short: "OpenPlantbook CLI - Plant care information from the command line",
		Long: `OpenPlantbook CLI provides access to the OpenPlantbook API,
a crowd-sourced database of plant care information.

Get your free API credentials at: https://open.plantbook.io/`,
		SilenceUsage: true,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.openplantbook.yaml)")
	rootCmd.PersistentFlags().String("api-key", "", "OpenPlantbook API key")
	rootCmd.PersistentFlags().String("client-id", "", "OAuth2 client ID")
	rootCmd.PersistentFlags().String("client-secret", "", "OAuth2 client secret")
	rootCmd.PersistentFlags().String("base-url", "", "API base URL (default: https://open.plantbook.io/api/v1)")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging")

	// Bind flags to viper
	viper.BindPFlag("api-key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("client-id", rootCmd.PersistentFlags().Lookup("client-id"))
	viper.BindPFlag("client-secret", rootCmd.PersistentFlags().Lookup("client-secret"))
	viper.BindPFlag("base-url", rootCmd.PersistentFlags().Lookup("base-url"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Add commands
	rootCmd.AddCommand(newSearchCmd())
	rootCmd.AddCommand(newDetailsCmd())
	rootCmd.AddCommand(newVersionCmd())

	cobra.OnInitialize(initConfig)

	return rootCmd
}

func initConfig() {
	// Load .env file if it exists (silently ignore if not found)
	_ = godotenv.Load()

	// Read from environment variables
	viper.SetEnvPrefix("OPENPLANTBOOK")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Read config file if specified or search for it
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		} else if viper.GetBool("debug") {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
		}
	} else {
		// Search for config in home directory
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home)
		}

		// Also search in current directory
		viper.AddConfigPath(".")
		viper.SetConfigName(".openplantbook")
		viper.SetConfigType("yaml")

		// Try to read config file (ignore error if not found)
		if err := viper.ReadInConfig(); err == nil && viper.GetBool("debug") {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
		}
	}
}

func newSearchCmd() *cobra.Command {
	var (
		limit      int
		userPlants bool
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search for plants by name or alias",
		Long: `Search for plants by common name or scientific name.

Examples:
  openplantbook search monstera
  openplantbook search fern --limit 5
  openplantbook search monstera --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			client, err := createClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			results, err := client.SearchPlants(context.Background(), query, &openplantbook.SearchOptions{
				Limit:      limit,
				UserPlants: userPlants,
			})
			if err != nil {
				return fmt.Errorf("search failed: %w", err)
			}

			if jsonOutput {
				return outputJSON(results)
			}

			return outputSearchResults(results)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to return")
	cmd.Flags().BoolVar(&userPlants, "user-plants", false, "Include user-contributed plants")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")

	return cmd
}

func newDetailsCmd() *cobra.Command {
	var (
		language   string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "details <pid>",
		Short: "Get detailed care information for a plant",
		Long: `Retrieve detailed care information for a specific plant by its PID.

Examples:
  openplantbook details monstera-deliciosa
  openplantbook details monstera-deliciosa --lang es
  openplantbook details monstera-deliciosa --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Normalize PID: convert hyphens to spaces (e.g., "monstera-deliciosa" -> "monstera deliciosa")
			// This allows users to use either format for convenience
			pid := strings.ReplaceAll(args[0], "-", " ")

			client, err := createClient()
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			details, err := client.GetPlantDetails(context.Background(), pid, &openplantbook.DetailOptions{
				Language: language,
			})
			if err != nil {
				return fmt.Errorf("failed to get details: %w", err)
			}

			if jsonOutput {
				return outputJSON(details)
			}

			return outputPlantDetails(details)
		},
	}

	cmd.Flags().StringVar(&language, "lang", "en", "Language code (ISO 639-1)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")

	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("openplantbook version %s\n", version)
			fmt.Printf("  commit: %s\n", commit)
			fmt.Printf("  built:  %s\n", date)
		},
	}
}

func createClient() (*openplantbook.Client, error) {
	opts := []openplantbook.Option{}

	// Authentication - check for API key first, then OAuth2
	apiKey := viper.GetString("api-key")
	clientID := viper.GetString("client-id")
	clientSecret := viper.GetString("client-secret")

	if apiKey != "" {
		opts = append(opts, openplantbook.WithAPIKey(apiKey))
	} else if clientID != "" && clientSecret != "" {
		opts = append(opts, openplantbook.WithOAuth2(clientID, clientSecret))
	} else {
		return nil, fmt.Errorf("no authentication provided: set OPENPLANTBOOK_API_KEY or OPENPLANTBOOK_CLIENT_ID/CLIENT_SECRET")
	}

	// Optional base URL override
	if baseURL := viper.GetString("base-url"); baseURL != "" {
		opts = append(opts, openplantbook.WithBaseURL(baseURL))
	}

	// Debug logging
	if viper.GetBool("debug") {
		logger := &cliLogger{slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))}
		opts = append(opts, openplantbook.WithLogger(logger))
	}

	return openplantbook.New(opts...)
}

func outputSearchResults(results []openplantbook.PlantSearchResult) error {
	if len(results) == 0 {
		fmt.Println("No plants found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SCIENTIFIC NAME\tCOMMON NAME\tPID\tCATEGORY")
	fmt.Fprintln(w, "---------------\t-----------\t---\t--------")
	for _, plant := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", plant.DisplayPID, plant.Alias, plant.PID, plant.Category)
	}
	w.Flush()
	fmt.Printf("\nFound %d plant(s)\n", len(results))
	return nil
}

func outputPlantDetails(details *openplantbook.PlantDetails) error {
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
	return nil
}

func outputJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
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

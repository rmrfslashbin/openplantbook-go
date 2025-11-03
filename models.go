package openplantbook

// PlantSearchResult represents a single plant in search results
type PlantSearchResult struct {
	PID        string `json:"pid"`
	DisplayPID string `json:"display_pid"`
	Alias      string `json:"alias"`
	Category   string `json:"category"`
}

// searchResponse wraps the paginated API response
type searchResponse struct {
	Count    int                 `json:"count"`
	Next     *string             `json:"next"`
	Previous *string             `json:"previous"`
	Results  []PlantSearchResult `json:"results"`
}

// PlantDetails represents complete plant care information
type PlantDetails struct {
	PID          string  `json:"pid"`
	DisplayPID   string  `json:"display_pid"`
	Alias        string  `json:"alias"`
	MaxLightLux  int     `json:"max_light_lux"`
	MinLightLux  int     `json:"min_light_lux"`
	MaxTemp      float64 `json:"max_temp"`
	MinTemp      float64 `json:"min_temp"`
	MaxEnvHumid  int     `json:"max_env_humid"`
	MinEnvHumid  int     `json:"min_env_humid"`
	MaxSoilMoist int     `json:"max_soil_moist"`
	MinSoilMoist int     `json:"min_soil_moist"`
	MaxSoilEC    int     `json:"max_soil_ec"`
	MinSoilEC    int     `json:"min_soil_ec"`
	ImageURL     string  `json:"image_url"`
	Category     string  `json:"category"`
}

// SearchOptions configures plant search behavior
type SearchOptions struct {
	// Limit is the maximum number of results to return (0 = API default)
	Limit int

	// UserPlants includes user-contributed plants in results
	UserPlants bool
}

// DetailOptions configures plant detail retrieval
type DetailOptions struct {
	// Language is the ISO 639-1 language code (e.g., "en", "de", "es")
	Language string
}

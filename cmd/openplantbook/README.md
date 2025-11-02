# OpenPlantbook CLI

A command-line interface for the OpenPlantbook API, providing plant care information directly from your terminal.

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/rmrfslashbin/openplantbook-go.git
cd openplantbook-go

# Build and install
make build-cli
sudo cp bin/openplantbook /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/rmrfslashbin/openplantbook-go/cmd/openplantbook@latest
```

## Authentication

Before using the CLI, you need to obtain API credentials from [OpenPlantbook](https://open.plantbook.io/).

### Option 1: API Key (Recommended)

```bash
export OPENPLANTBOOK_API_KEY="your-api-key-here"
```

Add to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) for persistence:

```bash
echo 'export OPENPLANTBOOK_API_KEY="your-api-key-here"' >> ~/.zshrc
```

### Option 2: OAuth2 Client Credentials

```bash
export OPENPLANTBOOK_CLIENT_ID="your-client-id"
export OPENPLANTBOOK_CLIENT_SECRET="your-client-secret"
```

## Usage

### Search for Plants

Search for plants by name or common alias:

```bash
# Basic search
openplantbook search monstera

# Limit results
openplantbook search fern --limit 5

# JSON output for scripting
openplantbook search monstera --json
```

**Output:**
```
SCIENTIFIC NAME       COMMON NAME         PID                    CATEGORY
---------------       -----------         ---                    --------
Monstera deliciosa    Monstera            monstera-deliciosa     Houseplant
Monstera adansonii    Swiss Cheese Vine   monstera-adansonii     Houseplant

Found 2 plant(s)
```

### Get Plant Details

Retrieve detailed care information for a specific plant:

```bash
# Get details by PID
openplantbook details monstera-deliciosa

# Get details in a different language
openplantbook details monstera-deliciosa --lang es

# JSON output
openplantbook details monstera-deliciosa --json
```

**Output:**
```
Plant: Monstera deliciosa
Common Name: Monstera
PID: monstera-deliciosa
Category: Houseplant

Care Requirements:
==================
Light (Lux):       2500 - 20000
Temperature (°C):  15.0 - 30.0
Humidity (%):      40 - 80
Soil Moisture (%): 15 - 60
Soil EC (μS/cm):   350 - 2000

Image: https://example.com/monstera.jpg
```

### Version Information

```bash
openplantbook version
```

**Output:**
```
openplantbook version v0.1.0
  commit: a1b2c3d
  built:  2025-11-02T21:00:00Z
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `OPENPLANTBOOK_API_KEY` | API key for authentication | Yes* |
| `OPENPLANTBOOK_CLIENT_ID` | OAuth2 client ID | Yes* |
| `OPENPLANTBOOK_CLIENT_SECRET` | OAuth2 client secret | Yes* |
| `OPENPLANTBOOK_BASE_URL` | Override API base URL | No |
| `OPENPLANTBOOK_DEBUG` | Enable debug logging (`true`/`false`) | No |

*Either API key OR OAuth2 credentials are required

## Scripting Examples

### Extract PIDs from Search Results

```bash
openplantbook search monstera --json | jq -r '.[].pid'
```

### Get Multiple Plant Details

```bash
for pid in $(openplantbook search fern --json | jq -r '.[].pid'); do
  openplantbook details "$pid"
  echo "---"
done
```

### Check Temperature Requirements

```bash
openplantbook details monstera-deliciosa --json | \
  jq '{plant: .display_pid, min_temp: .min_temp, max_temp: .max_temp}'
```

### CSV Export

```bash
echo "Scientific Name,Common Name,PID,Category" > plants.csv
openplantbook search monstera --json | \
  jq -r '.[] | [.display_pid, .alias, .pid, .category] | @csv' >> plants.csv
```

## Error Handling

The CLI uses standard exit codes:

- `0` - Success
- `1` - General error (invalid input, API error, etc.)

Example error handling in scripts:

```bash
if openplantbook search "unknown-plant" --json > results.json 2>&1; then
  echo "Search successful"
else
  echo "Search failed"
  exit 1
fi
```

## Building

### Build for Current Platform

```bash
make build-cli
```

Binary will be created at `bin/openplantbook`

### Build for All Platforms

```bash
make build-cli-all
```

Creates binaries for:
- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64)

## Development

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage
```

### Debug Mode

Enable debug logging to see API requests and cache hits:

```bash
export OPENPLANTBOOK_DEBUG=true
openplantbook search monstera
```

## Troubleshooting

### "no authentication provided" Error

Make sure you've exported your API credentials:

```bash
# Check if API key is set
echo $OPENPLANTBOOK_API_KEY

# If empty, export it
export OPENPLANTBOOK_API_KEY="your-api-key"
```

### Rate Limit Errors

The CLI respects the API's rate limit (200 requests/day by default). If you exceed this:

- Wait for the rate limit to reset (24 hours)
- Use caching effectively (search results cached for 1 hour, details for 24 hours)
- Consider upgrading your API plan

### "plant not found" Error

- Verify the PID is correct
- Try searching first to find the correct PID
- Check plant availability in the API

## See Also

- [OpenPlantbook SDK Documentation](../../README.md)
- [SDK Examples](../../examples/README.md)
- [OpenPlantbook API Documentation](https://open.plantbook.io/docs/)

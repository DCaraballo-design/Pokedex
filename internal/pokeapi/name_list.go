package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// LocationNames fetches Pokémon encounter data for a specific location area
func (c *Client) LocationNames(locationName string) error {
	// Construct the URL for the specific location area
	url := fmt.Sprintf("%s/location-area/%s", baseURL, locationName)

	// Check the cache first
	if val, ok := c.cache.Get(url); ok {
		fmt.Println("Using cached data...")
		return displayPokemonNames(val)
	}

	// Make the HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch location data: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Cache the response data
	c.cache.Add(url, data)

	// Display the Pokémon names
	return displayPokemonNames(data)
}

// displayPokemonNames parses and displays Pokémon names from the response data
func displayPokemonNames(data []byte) error {
	var encountersResp RespPokemonEncounters
	err := json.Unmarshal(data, &encountersResp)
	if err != nil {
		return fmt.Errorf("failed to parse location data: %v", err)
	}

	// Display the location name
	fmt.Printf("Exploring %s...\n", encountersResp.Name)

	// Display the Pokémon names
	fmt.Println("Found Pokémon:")
	for _, encounter := range encountersResp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

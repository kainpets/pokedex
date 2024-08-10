package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"github.com/kainpets/pokedex/internal/pokecache"
	"strings"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type locationAreaResponse struct {
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var (
	cache  *pokecache.Cache
	offset int = 0
)

func main() {
	for {
		cache = pokecache.NewCache(5 * time.Minute)
		fmt.Println("Welcome to the Pokedex!")
		displayHelp()
		fmt.Print("pokedex > ")
		input := getInput()
		if err := parseInput(input); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	return strings.TrimSpace(input)
}

func parseInput(input string) error {
	commands := []cliCommand{
		{
			name:        "exit",
			description: "Exit the program",
			callback:    func() error { os.Exit(0); return nil },
		},
		{
			name:        "help",
			description: "Display this help message",
			callback:    displayHelp,
		},
		{
			name:        "map",
			description: "Display map (next 20 locations)",
			callback:    displayMap,
		},
		{
			name:        "mapb",
			description: "Display map (previous 20 locations)",
			callback:    displayMapBack,
		},
	}

	input = strings.ToLower(input)
	for _, cmd := range commands {
		if input == cmd.name {
			return cmd.callback()
		}
	}

	return fmt.Errorf("unknown command. Type 'help' for a list of commands")
}

func displayHelp() error {
	fmt.Println("Commands:")
	fmt.Println("map: Display map (next 20 locations)")
	fmt.Println("mapb: Display map (previous 20 locations)")
	fmt.Println("exit: Exit the program")
	fmt.Println("help: Display this help message")
	return nil
}

func displayMap() error {
	locations, err := getLocations(offset)
	if err != nil {
		return err
	}

	if len(locations) == 0 {
		fmt.Println("No more location to display.")
		return nil
	}

	for _, loc := range locations {
		fmt.Println(loc)
	}
	offset += 20
	return nil
}

func displayMapBack() error {
	if offset >= 20 {
		offset -= 20
	} else {
		return fmt.Errorf("no previous locations to display")
	}
	locations, err := getLocations(offset)
	if err != nil {
		return err
	}
	for _, loc := range locations {
		fmt.Println(loc)
	}
	return nil
}

func getLocations(offset int) ([]string, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?offset=%d&limit=20", offset)

	if cachedData, found := cache.Get(url); found {
		var result locationAreaResponse
		err := json.Unmarshal(cachedData, &result)
		if err != nil {
			return nil, fmt.Errorf("error parsing cached response: %w", err)
		}
		return extractLocationNames(result)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching locations: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	cache.Add(url, body)

	var result locationAreaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing response body: %w", err)
	}

	return extractLocationNames(result)
}

func extractLocationNames(result locationAreaResponse) ([]string, error) {
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no locations found")
	}

	locations := make([]string, len(result.Results))
	for i, loc := range result.Results {
		locations[i] = loc.Name
	}
	return locations, nil
}

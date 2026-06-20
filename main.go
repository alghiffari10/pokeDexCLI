package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/alghiffari10/pokeDexCLI/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	Next     *string
	Previous *string
	cache    *pokecache.Cache
}
type ExploreResponse struct {
	PokemmonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type locationResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func mapCommand(cfg *config, cmd []string) error {

	var url = "https://pokeapi.co/api/v2/location-area/"

	if cfg.Next != nil {
		url = *cfg.Next
	}

	dataByte, err := getLocationData(url, cfg.cache)
	if err != nil {
		return err
	}

	var mapPokemon locationResponse
	if err := json.Unmarshal(dataByte, &mapPokemon); err != nil {
		return err
	}

	cfg.Next = mapPokemon.Next
	cfg.Previous = mapPokemon.Previous

	for _, location := range mapPokemon.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func mapCommandBack(cfg *config, cmd []string) error {

	if cfg.Previous == nil {
		fmt.Println("you already on the last page")
		return nil
	}

	dataByte, err := getLocationData(*cfg.Previous, cfg.cache)
	if err != nil {
		return err
	}

	var mapPokemon locationResponse
	if err := json.Unmarshal(dataByte, &mapPokemon); err != nil {
		return err
	}

	cfg.Next = mapPokemon.Next
	cfg.Previous = mapPokemon.Previous

	for _, location := range mapPokemon.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func exitCommand(cfg *config, cmd []string) error {

	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func exploreCommand(cfg *config, cmd []string) error {

	if len(cmd) == 0 {
		return fmt.Errorf("Please provide a location area")
	}

	areaName := cmd[0]
	fmt.Printf("Exploring %v...", areaName)

	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	data, err := getLocationData(url, cfg.cache)
	if err != nil {
		return err
	}

	var explore ExploreResponse

	if err := json.Unmarshal(data, &explore); err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")

	for _, encounter := range explore.PokemmonEncounters {
		fmt.Printf(" - %v\n", encounter.Pokemon.Name)
	}

	return nil
}

func helpCommand(cfg *config, cmd []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp\t: Displays a help message\nmap\t: Displays a map\nmapb\t: Displaying a previous location\nexplore\t: Displaying list of pokemon based on location being choose\nexit\t: Exit the Pokedex")
	return nil
}
func getLocationData(url string, cache *pokecache.Cache) ([]byte, error) {

	if data, ok := cache.Get(url); ok {
		fmt.Println("Using cache")
		return data, nil
	}

	fmt.Println("Making HTTP request")

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	cache.Add(url, data)

	return data, nil
}

func main() {

	cfg := &config{
		Next:     nil,
		Previous: nil,
		cache:    pokecache.NewCache(5 * time.Second),
	}

	supportCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    exitCommand,
		},

		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    helpCommand,
		},

		"map": {
			name:        "map",
			description: "Displaying a next location",
			callback:    mapCommand,
		},

		"mapb": {
			name:        "Map Back",
			description: "Displaying a previous location",
			callback:    mapCommandBack,
		},

		"explore": {
			name:        "Explore Command",
			description: "Displaying list of pokemon based on location being choose",
			callback:    exploreCommand,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()

		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}
		scanner.Err()

		commandName := input[0]
		args := input[1:]

		command, ok := supportCommands[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(cfg, args)
		if err != nil {
			fmt.Println(err)
		}
	}

}

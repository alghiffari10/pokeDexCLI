package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/alghiffari10/pokeDexCLI/internal/pokecache"
)

// TODO: Refactor the code to make it more readable

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type PokemonResponse struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`

	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`

	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
	Pokedex  map[string]PokemonResponse
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

	dataByte, err := getLocationData(url, cfg.Cache)
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

	dataByte, err := getLocationData(*cfg.Previous, cfg.Cache)
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
	fmt.Printf("Exploring %v...\n", areaName)

	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	data, err := getLocationData(url, cfg.Cache)
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

func catchCommand(cfg *config, cmd []string) error {
	if len(cmd) == 0 {
		return fmt.Errorf("Please provide pokemon name")
	}

	PokemonName := cmd[0]
	fmt.Printf("Throwing a Pokeball at %v...\n", PokemonName)

	url := "https://pokeapi.co/api/v2/pokemon/" + PokemonName

	data, err := getLocationData(url, cfg.Cache)
	if err != nil {
		return err
	}

	var pokemon PokemonResponse
	if err := json.Unmarshal(data, &pokemon); err != nil {
		return nil
	}

	prop := rand.Intn(pokemon.BaseExperience)

	if prop < 35 {
		fmt.Printf("%v was caught!\n", pokemon.Name)
		cfg.Pokedex[pokemon.Name] = pokemon
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}

	return nil
}

func inspectCommand(cfg *config, cmd []string) error {

	if len(cmd) == 0 {
		return fmt.Errorf("Please provide a pokemon name")
	}
	name := cmd[0]

	pokemon, ok := cfg.Pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("\t- %v: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("\t- %v\n", pokemonType.Type.Name)
	}

	return nil
}

func pokedexCommand(cfg *config, cmd []string) error {
	if len(cfg.Pokedex) == 0 {
		return fmt.Errorf("You don't have any pokemon")
	}
	fmt.Println("Your Pokedex:")
	for name := range cfg.Pokedex {
		fmt.Printf("\t- %v\n", name)
	}
	return nil
}

func helpCommand(cfg *config, cmd []string) error {
	fmt.Println(`Welcome to the Pokedex!
Usage:

help       : Displays a help message
map        : Displays the next locations
mapb       : Displays the previous locations
explore    : Displays Pokémon in a location area
catch      : Attempt to catch a Pokémon
inspect    : Inspect a caught pokemon
pokedex    : Check pokemon you have
exit       : Exit the Pokedex`)
	return nil
}
func getLocationData(url string, cache *pokecache.Cache) ([]byte, error) {

	if data, ok := cache.Get(url); ok {
		fmt.Println("Using cache")
		return data, nil
	}

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
		Cache:    pokecache.NewCache(5 * time.Second),
		Pokedex:  make(map[string]PokemonResponse),
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

		"catch": {
			name:        "Catch Command",
			description: "Catch your Pokemon",
			callback:    catchCommand,
		},

		"inspect": {
			name:        "Inspect Command",
			description: "Inspect a caught pokemon",
			callback:    inspectCommand,
		},

		"pokedex": {
			name:        "Pokedex Command",
			description: "Check pokemon you have",
			callback:    pokedexCommand,
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

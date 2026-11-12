# Pokedex CLI

A small interactive command-line Pokedex written in Go. It talks to the public
[PokéAPI](https://pokeapi.co/), lets you browse location areas, explore the
Pokémon that live there, catch them, and inspect the ones you have caught.

## Features

- Interactive REPL with a `Pokedex >` prompt
- Paginated browsing of location areas with `map` and `mapb`
- Explore a location area to see which Pokémon appear there
- Catch Pokémon with a chance-based throw
- Inspect caught Pokémon (height, weight, stats, and types)
- View everything currently in your Pokedex
- In-memory HTTP response cache with a 5-second TTL

## Requirements

- [Go](https://go.dev/dl/) 1.26 or newer
- An internet connection (data comes from `pokeapi.co`)

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/alghiffari10/pokeDexCLI.git
cd pokeDexCLI
go build -o pokeDexCLI
```

Or run it directly without building:

```bash
go run .
```

## Usage

Start the program and type commands at the prompt:

```text
$ ./pokeDexCLI
Pokedex > help
Welcome to the Pokedex!
Usage:

help       : Displays a help message
map        : Displays the next locations
mapb       : Displays the previous locations
explore    : Displays Pokémon in a location area
catch      : Attempt to catch a Pokémon
inspect    : Inspect a caught pokemon
pokedex    : Check pokemon you have
exit       : Exit the Pokedex
```

### Commands

| Command              | Description                                              |
| -------------------- | -------------------------------------------------------- |
| `help`               | Show the list of available commands                      |
| `map`                | Show the next page of location areas                     |
| `mapb`               | Show the previous page of location areas                 |
| `explore <area>`     | List the Pokémon found in a location area                |
| `catch <pokemon>`    | Attempt to catch a Pokémon                               |
| `inspect <pokemon>`  | Show details for a Pokémon you have caught               |
| `pokedex`            | List all Pokémon you have caught                         |
| `exit`               | Quit the program                                         |

### Example session

```text
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
...

Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - wingull
 - starly

Pokedex > catch starly
Throwing a Pokeball at starly...
starly was caught!
You may now inspect it with the inspect command.

Pokedex > inspect starly
Name: starly
Height: 3
Weight: 20
Stats:
	- hp: 40
	- attack: 55
	...
Types:
	- normal
	- flying

Pokedex > pokedex
Your Pokedex:
	- starly

Pokedex > exit
Closing the Pokedex... Goodbye!
```

## How it works

The program is a simple read-eval-print loop:

1. Read a line from standard input and split it into words.
2. Look up the first word in the command table in `main.go`.
3. Call the matching handler with the remaining words as arguments.
4. Print any error and prompt again.

All Pokémon data is fetched from PokéAPI over HTTP. Raw response bodies are
stored in a small in-memory cache (`internal/pokecache`) so repeated requests
within the TTL do not hit the network again. Caught Pokémon are kept in an
in-memory `map[string]PokemonResponse` for the duration of the session.

Catching uses the Pokémon's base experience as a difficulty value: a random
number below 35 results in a successful catch, otherwise the Pokémon escapes.

## Project structure

```text
.
├── main.go                        # REPL loop, commands, and PokéAPI types
├── repl.go                        # Input parsing (cleanInput)
├── main_test.go                   # Placeholder test for main
├── repl_test.go                   # Tests for input parsing
├── internal/
│   └── pokecache/
│       ├── pokecache.go           # Thread-safe cache with periodic reaping
│       └── pokecache_test.go      # Cache tests
└── go.mod
```

## Testing

Run the full test suite:

```bash
go test ./...
```

Run with verbose output:

```bash
go test -v ./...
```

## Notes

- The Pokédex is not persisted; caught Pokémon are lost when you exit.
- Cache entries expire 5 seconds after they are added.
- `main_test.go` is a placeholder and currently contains no test cases.

## License

No license has been specified for this project.

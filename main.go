package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/DreamEcho100/boot-dot-dev-build-a-pokedex-in-go/internal/cache"
	"github.com/DreamEcho100/boot-dot-dev-build-a-pokedex-in-go/internal/fetch"
	"github.com/DreamEcho100/boot-dot-dev-build-a-pokedex-in-go/internal/poke_api"
)

type cliCommandConfig struct {
	previous *string
	next     *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(...string) error
	config      *cliCommandConfig
}

func stringPtr(s string) *string {
	return &s
}

var (
	commands       map[string]cliCommand
	mapCliConfig   cliCommandConfig
	caughtPokemons map[string]poke_api.PokeAPIPokemonResult
	pokemonCache      cache.Cache
)

func init() {
	caughtPokemons = make(map[string]poke_api.PokeAPIPokemonResult)
	pokemonCache = *cache.NewCache(10 * time.Second)

	mapCliConfig = cliCommandConfig{
		previous: nil,
		next:     stringPtr(poke_api.PokeAPILocationAreaListBaseURL),
	}

	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
			config:      nil,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
			config:      nil,
		},
		"map": {
			name:        "map",
			description: "Displays a map of the area",
			callback:    commandMap,
			config:      &mapCliConfig,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays a map of the area",
			callback:    commandMapB,
			config:      &mapCliConfig,
		},
		"explore": {
			name:        "explore",
			description: "Find a location area and list of all the Pokémon located there",
			callback:    commandExplore,
			config:      nil,
		},
		"catch": {
			name:        "catch",
			description: "Try catching the Pokemon",
			callback:    commandCatch,
			config:      nil,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect one of your caught pokemons",
			callback:    commandInspect,
			config:      nil,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See all of the names of the Pokemon that you caught!",
			callback:    commandPokedex,
			config:      nil,
		},
	}
}

func commandExit(args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	// Iterate over the commands map to print descriptions
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(args ...string) error {
	cliCommandConfig := commands["map"].config

	nextURL := cliCommandConfig.next
	if nextURL == nil {
		fmt.Println("No more location areas to display.")
		return nil
	}

	locationAreas, err := fetch.Get[poke_api.PokeAPILocationAreaResultsList](*nextURL, &pokemonCache)
	if err != nil {
		return err
	}
	cliCommandConfig.previous = locationAreas.Previous
	cliCommandConfig.next = locationAreas.Next

	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}

func commandMapB(args ...string) error {
	cliCommandConfig := commands["mapb"].config

	previousURL := cliCommandConfig.previous
	if previousURL == nil {
		fmt.Println("No previous location areas to display.")
		return nil
	}

	previousLocationAreas, err := fetch.Get[poke_api.PokeAPILocationAreaResultsList](*previousURL, &pokemonCache)
	if err != nil {
		return err
	}
	cliCommandConfig.previous = previousLocationAreas.Previous
	cliCommandConfig.next = previousLocationAreas.Next

	for _, locationArea := range previousLocationAreas.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}

func commandExplore(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no location area name provided")
	}

	areaName := args[0]

	fmt.Printf("Exploring %s...\n", areaName)
	data, err := fetch.Get[poke_api.PokeAPILocationAreaResult](
		poke_api.MakePokeAPILocationAreaBaseURL(areaName), &pokemonCache,
	)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, item := range data.PokemonEncounters {
		fmt.Println(" - " + item.Pokemon.Name)
	}

	return nil
}

func commandCatch(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no pokemon name provided")
	}

	pokemonIDOrName := args[0]

	data, err := fetch.Get[poke_api.PokeAPIPokemonResult](
		poke_api.MakePokeAPIPokemonURL(pokemonIDOrName), &pokemonCache,
	)
	if err != nil {
		return err
	}

	// TODO: add a gauard a gainst already caught ones

	fmt.Printf("Throwing a Pokeball at %s...\n", data.Name)

	catchChance := rand.Intn(data.BaseExperience)

	if catchChance < 50 {
		caughtPokemons[data.Name] = data
		fmt.Printf("%s was caught!\n", caughtPokemons[data.Name].Name)
	} else {
		fmt.Printf("%s escaped!\n", data.Name)
	}

	return nil
}

func commandInspect(args ...string) error {
	pokemonName := args[0]

	if len(args) == 0 {
		return fmt.Errorf("no pokemon name provided")
	}

	pokemon, ok := caughtPokemons[pokemonName]

	if !ok {
		fmt.Printf("Can't inspect %s, it's not caught!\n", pokemonName)
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf(" - %s\n", pokemonType.Type.Name)
	}

	return nil
}

func commandPokedex(args ...string) error {
	if len(caughtPokemons) == 0 {
		fmt.Println("No caught pokemon yet!")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, value := range caughtPokemons {
		fmt.Printf(" - %s\n", value.Name)
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex >")

		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		cleanedInput := strings.ToLower(strings.TrimSpace(input))

		words := strings.Fields(cleanedInput)

		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		cmd, ok := commands[commandName]

		if ok {
			if err := cmd.callback(words[1:]...); err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}

		// firstWord := words[0]
		// fmt.Println("Your command was:", firstWord)

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "reading input: %v\n", err)
			os.Exit(1)
		}

	}
}

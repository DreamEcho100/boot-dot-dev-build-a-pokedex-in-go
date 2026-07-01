package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	callback    func() error
	config      *cliCommandConfig
}

func stringPtr(s string) *string {
	return &s
}

var (
	commands     map[string]cliCommand
	mapCliConfig cliCommandConfig
)

func init() {
	mapCliConfig = cliCommandConfig{
		previous: nil,
		next:     stringPtr(poke_api.PokeAPILocationAreaBaseURL),
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
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	// Iterate over the commands map to print descriptions
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap() error {
	cliCommandConfig := commands["map"].config

	nextURL := cliCommandConfig.next
	if nextURL == nil {
		fmt.Println("No more location areas to display.")
		return nil
	}

	locationAreas, err := fetch.Get[poke_api.PokeAPILocationAreaResultsList](*nextURL)
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

func commandMapB() error {
	cliCommandConfig := commands["mapb"].config

	previousURL := cliCommandConfig.previous
	if previousURL == nil {
		fmt.Println("No previous location areas to display.")
		return nil
	}

	previousLocationAreas, err := fetch.Get[poke_api.PokeAPILocationAreaResultsList](*previousURL)
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
			if err := cmd.callback(); err != nil {
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

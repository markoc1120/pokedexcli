package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/markoc1120/pokedexcli/internal"
)

type commandsArg map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	Callback    func(*commandsArg, *internal.Config) error
}

func commandExit(commands *commandsArg, config *internal.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands *commandsArg, config *internal.Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cmd := range *commands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(commands *commandsArg, config *internal.Config) error {
	locationArea, err := internal.GetLocation(config)
	if err != nil {
		return err
	}
	config.Next = locationArea.Next
	config.Previous = locationArea.Previous
	for _, loc := range locationArea.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(commands *commandsArg, config *internal.Config) error {
	if config.Previous == nil {
		return fmt.Errorf("there is no previous locations")
	}
	config.Next = config.Previous
	commandMap(commands, config)
	return nil
}

func commandExplore(commands *commandsArg, config *internal.Config) error {
	fmt.Println("Exploring pastoria-city-area...")

	locationAreaDetail, err := internal.GetPokemon(config, config.Arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationAreaDetail.PokemonEncounters {
		fmt.Printf(" - %v\n", pokemon.Pokemon.Name)
	}
	return nil
}

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func GetCommands() commandsArg {
	commands := commandsArg{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 LocationAreas",
			Callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 LocationAreas",
			Callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays pokemons in a specific LocationArea",
			Callback:    commandExplore,
		},
	}
	return commands
}

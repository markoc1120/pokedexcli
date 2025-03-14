package app

import (
	"fmt"
	"math/rand"
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

	locationAreaDetail, err := internal.GetLocationDetail(config, config.Arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range locationAreaDetail.PokemonEncounters {
		fmt.Printf(" - %v\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(commands *commandsArg, config *internal.Config) error {
	pokemonName := config.Arguments[0]
	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)
	pokemon, err := internal.GetPokemon(config, pokemonName)
	if err != nil {
		return fmt.Errorf("there is no pokemon called: %v", pokemonName)
	}
	probabilityOfCatch := rand.Float32() * float32(pokemon.BaseExperience)
	if probabilityOfCatch < 30.0 {
		fmt.Printf("%v escaped!\n", pokemonName)
	} else {
		fmt.Printf("%v was caught!\n", pokemonName)
		config.Pokemons[pokemonName] = pokemon
	}
	return nil
}

func commandList(commands *commandsArg, config *internal.Config) error {
	fmt.Print("Listing all your pokemons...\n")
	if len(config.Pokemons) == 0 {
		fmt.Println("You haven't caught any pokemons yet")
		return nil
	}
	for _, pokemon := range config.Pokemons {
		fmt.Printf(" - %v\n", pokemon.Name)
	}
	return nil
}

func commandInspect(commands *commandsArg, config *internal.Config) error {
	pokemonName := config.Arguments[0]
	if pokemon, ok := config.Pokemons[pokemonName]; !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %v\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Printf("Stats: \n")
		for _, stat := range pokemon.Stats {
			fmt.Printf(" -%v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		for _, tp := range pokemon.Types {
			fmt.Printf(" - %v\n", tp.Type.Name)
		}
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
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			Callback:    commandCatch,
		},
		"list": {
			name:        "list",
			description: "List all the pokemons you caught previously",
			Callback:    commandList,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon you caught previously",
			Callback:    commandInspect,
		},
	}
	return commands
}

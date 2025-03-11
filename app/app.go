package app

import (
	"fmt"
	"os"
	"strings"
)

type commandsArg map[string]cliCommand

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(commands *commandsArg) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands *commandsArg) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, cmd := range *commands {
		fmt.Printf("%v: %v\n", cmd.name, cmd.description)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*commandsArg) error
}

func GetCommands() commandsArg {
	commands := commandsArg{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
	return commands
}

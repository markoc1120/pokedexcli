package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/markoc1120/pokedexcli/app"
	"github.com/markoc1120/pokedexcli/internal"
)

func main() {
	commands := app.GetCommands()
	config := internal.Config{
		Cache: internal.NewCache(5 * time.Minute),
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := app.CleanInput(scanner.Text())
		if command, ok := commands[input[0]]; ok {
			err := command.Callback(&commands, &config)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unkown command")
		}
	}
}

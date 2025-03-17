package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/markoc1120/pokedexcli/app"
	"github.com/markoc1120/pokedexcli/internal"
)

func main() {
	commands := app.GetCommands()
	config := internal.Config{
		Cache:    internal.NewCache(5 * time.Minute),
		Pokemons: make(map[string]internal.Pokemon),
	}

	l, err := readline.NewEx(&readline.Config{
		Prompt:            "Pokedex > ",
		HistoryFile:       app.HistoryFilePath,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()
	l.CaptureExitSignal()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		input := strings.Fields(strings.TrimSpace(line))
		config.Arguments = input[1:]
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

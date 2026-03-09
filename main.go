package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedex/helpers"
	"pokedex/internal/pokecache"
	"time"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	cache := pokecache.NewCache(time.Minute)

	pokedex := make(map[string]Pokemon)

	commands := map[string]CliCommands{
		"exit": {
			name:        "exit",
			description: "terminates the program",
			callback:    Exit,
		},
		"help": {
			name:        "help",
			description: "lists all the commands",
			callback:    Help,
		},
		"map": {
			name:        "map",
			description: "each subsequent call displays 20 new locations",
			callback:    Map,
		},
		"mapb": {
			name:        "mapb",
			description: "prev area's page (limit=20)",
			callback:    Mapb,
		},
		"explore": {
			name:        "explore",
			description: "get all pokemons that are located in selected area",
			callback:    Explore,
		},
		"catch": {
			name:        "catch",
			description: "catch a pokemon with a random chance",
			callback:    Catch,
		},
		"inspect": {
			name:        "inspect",
			description: "shows stats of chooses pokemon",
			callback:    Inspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "shows you're pokemon collection",
			callback:    Pokedex,
		},
	}

	conf := &config{
		prevUrl:  nil,
		nextUrl:  nil,
		cache:    cache,
		commands: commands,
	}

	fmt.Println("Welcome to Pokedex! Type help to start")

	// game loop
	for {
		fmt.Print("\033[1;33mPokedex > \033[0m")

		if !scanner.Scan() {
			break
		}

		input := helpers.CleanInput(scanner.Text())
		if len(input) == 0 {

			fmt.Print("\tPlease provide a command. example: exit (to terminate) or help (to start)\n")
			continue
		}

		cmd := input[0]
		command, ok := commands[cmd]
		if !ok {
			fmt.Printf("\tUnknown command: %s, type help for help\n", cmd)
			continue
		}

		if len(input) > 1 {
			conf.option = &input[1]
		}

		command.callback(conf, pokedex)
	}

}

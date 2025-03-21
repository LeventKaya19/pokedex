package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/LeventKaya19/pokedex/mapper"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     *string
	Previous *string
}

var commands = map[string]CliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Help",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Map",
		callback:    commandMap,
	},
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	c := Config{}
	*c.Next = "https://pokeapi.co/api/v2/location-area"
	for {
		fmt.Print("Pokedex > ")
		s.Scan()
		res := cleanInput(strings.ToLower(s.Text()))
		if res[0] == "exit" || res[0] == "help" || res[0] == "map" {
			commands[res[0]].callback(&c)
		}
	}
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n\nhelp: Displays a help message\nexit: Exit the Pokedex")
	return nil
}

func commandMap(config *Config) error {
	if config.Next == nil {
		return errors.New("no next area to show")
	}
	results, err := mapper.GetMap(*config.Next)
	if err != nil {
		return err
	}

	for _, location := range results.Results {
		fmt.Println(location.Name)
	}
	*config.Next = results.Next

	return nil
}

func cleanInput(text string) []string {
	res := []string{}
	slice := strings.Split(text, " ")

	for _, s := range slice {
		if s != "" {
			res = append(res, s)
		}
	}

	return res
}

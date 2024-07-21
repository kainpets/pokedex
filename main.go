package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

func main() {
	for {
		fmt.Println("Welcome to the Pokedex!")
		fmt.Println("pokedex >")
		input := getInput()
		if err := parseInput(input); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)	
		return ""
	}
	return strings.TrimSpace(input)
}


func parseInput(input string) error {
	commands := []cliCommand{
		{
			name:        "exit",
			description: "Exit the program",
			callback:    func() error { os.Exit(0); return nil },
		},
		{
			name:        "help",
			description: "Display this help message",
			callback:    displayHelp,
		},
	}

	input = strings.ToLower(input)
	for _, cmd := range commands {
		if input == cmd.name {
			return cmd.callback()
		}
	}

	return fmt.Errorf("unknown command. Type 'help' for a list of commands")
}

func displayHelp() error {
	fmt.Println("Commands:")
	fmt.Println("exit: Exit the program")
	fmt.Println("help: Display this help message")
	return nil
}

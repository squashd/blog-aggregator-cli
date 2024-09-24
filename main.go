package main

import (
	"fmt"
	"log"
	"os"

	"github.com/squashd/blog-aggregator-cli/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	appState := state{cfg: &cfg}

	commands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := []string{}
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}

	err = commands.run(&appState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}

}

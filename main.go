package main

import (
	"fmt"
	"log"
	"os"

	"github/ansht2000/BlogAggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	programState := &state{cfg: &cfg}
	commands := commands{command_map: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	if err = commands.run(programState, command{Name: cmdName, Args: cmdArgs}); err != nil {
		log.Fatal(err)
	}
}
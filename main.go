package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ArashPoorazam/Gator/internal/config"
)

func main() {
	fmt.Println("Dont forget to build the app first then test it!!!")

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// creating state and command handler
	State := state{&cfg}
	commands := commands{handler: make(map[string]func(s *state, cmd command) error)}

	// commands
	commands.register("login", handlerLogin)

	// get args
	argsSlice := os.Args[:]
	if len(argsSlice) < 2 {
		log.Fatal("Peogram need at least one argument")
	}

	// input command
	newCommand := command{Name: argsSlice[1], Args: argsSlice[2:]}
	err = commands.run(&State, newCommand)
	if err != nil {
		log.Fatal(err)
	}
}

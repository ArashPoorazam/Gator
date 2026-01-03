package main

import (
	"fmt"

	"github.com/ArashPoorazam/Gator/internal/database"
)

func handlerHelp(s *state, cmd command, user database.User) error {
	fmt.Println("I am too lazy rn maybe in the future updates I will write a help handler.")

	return nil
}

package main

import (
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username.")
	}

	username := ""
	for _, arg := range cmd.Args {
		username += " " + arg
		if len(username) > 30 {
			return fmt.Errorf("username is too long")
		}
	}
	username = username[1:]

	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set\n", username)

	return nil
}

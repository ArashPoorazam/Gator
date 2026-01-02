package main

import (
	"context"
	"fmt"
)

func handlerDeleteUser(s *state, cmd command) error {
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

	user, err := s.Queries.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user does not exist: %w", err)
	}

	err = s.Queries.DeleteUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}

	return nil
}

func handlerResetTable(s *state, cmd command) error {
	err := s.Queries.ResetTable(context.Background())
	if err != nil {
		return fmt.Errorf("could not reset table: %w", err)
	}

	fmt.Println("Users Table Reset!")

	return nil
}

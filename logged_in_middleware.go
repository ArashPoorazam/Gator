package main

import (
	"context"
	"fmt"

	"github.com/ArashPoorazam/Gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		if cmd.Name == "register" {
			return handler(s, cmd, database.User{})
		}

		user, err := s.Queries.GetUser(context.Background(), s.Config.Current_user_name)
		if err != nil {
			return fmt.Errorf("please register first: %w", err)
		}
		return handler(s, cmd, user)
	}
}

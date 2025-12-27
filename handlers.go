package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ArashPoorazam/Gator/internal/database"
	"github.com/google/uuid"
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

func handlerRegister(s *state, cmd command) error {
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

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	_, err := s.Queries.GetUser(context.Background(), username)
	if err == nil {
		os.Exit(1)
	}

	user, err := s.Queries.CreateUser(context.Background(), newUser)
	if err != nil {
		return fmt.Errorf("could not create new user in database: %w", err)
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("could not set new user: %w", err)
	}

	return nil
}

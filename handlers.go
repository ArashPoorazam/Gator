package main

import (
	"context"
	"fmt"
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

	user, err := s.Queries.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user does not exist: %w", err)
	}

	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set\n", user.Name)

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
		return fmt.Errorf("user already exist: %w", err)
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

func handlerGetUsers(s *state, cmd command) error {
	names, err := s.Queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not get users names: %w", err)
	}

	for i, name := range names {
		if i >= 50 {
			return nil
		}
		if s.Config.Current_user_name == name {
			fmt.Printf("* %s (current)\n", name)
			continue
		}
		fmt.Printf("* %s\n", name)
	}

	return nil
}

package main

import (
	"fmt"

	"github.com/ArashPoorazam/Gator/internal/config"
)

type state struct {
	Config *config.Config
}

type command struct {
	Name string
	Args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username.")
	}

	username := ""
	for _, arg := range cmd.Args {
		username += arg + " "
		if len(username) > 30 {
			return fmt.Errorf("username is too long")
		}
	}

	err := s.Config.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User %s has been set\n", username)

	return nil
}

type commands struct {
	handler map[string]func(s *state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {
	f := c.handler[cmd.Name]

	// run command
	err := f(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.handler[name] = f
}

package main

import (
	"errors"

	"github.com/ArashPoorazam/Gator/internal/config"
)

type state struct {
	Config *config.Config
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(s *state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.registeredCommands[name] = f
}

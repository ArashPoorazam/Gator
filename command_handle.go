package main

import (
	"errors"

	"github.com/ArashPoorazam/Gator/internal/config"
	"github.com/ArashPoorazam/Gator/internal/database"
)

type state struct {
	Config  *config.Config
	Queries *database.Queries
}

type command struct {
	Name string
	Args []string
}

type commands struct {
	addedCommands map[string]func(s *state, cmd command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.addedCommands[cmd.Name]
	if !ok {
		return errors.New("command not found! use 'Gator help'")
	}
	return f(s, cmd)
}

func (c *commands) add(name string, f func(s *state, cmd command) error) {
	c.addedCommands[name] = f
}

package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	command_map map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.command_map[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	callback, ok := c.command_map[cmd.Name]
	if !ok {
		return errors.New("given command has not been registered")
	}
	return callback(s, cmd)
}
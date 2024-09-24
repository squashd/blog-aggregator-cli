package main

import (
	"fmt"
)

type (
	command struct {
		Name string
		Args []string
	}

	commands struct {
		registeredCommands map[string]func(*state, command) error
	}
)

func (c *commands) register(name string, handler func(*state, command) error) {
	c.registeredCommands[name] = handler
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

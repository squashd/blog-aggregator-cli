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
		registeredCmds map[string]func(*state, command) error
	}
)

func (c *commands) register(name string, handler func(*state, command) error) {
	c.registeredCmds[name] = handler
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registeredCmds[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

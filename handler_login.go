package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login requires a username")
	}
	s.cfg.SetUser(cmd.Args[0])

	fmt.Printf("Logged in as %s\n", cmd.Args[0])

	return nil
}

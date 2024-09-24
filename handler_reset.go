package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to reset users table %w", err)
	}
	fmt.Println("Users table reset")
	return nil
}

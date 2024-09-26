package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/squashd/blog-aggregator-cli/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}
	userName := cmd.Args[0]
	user, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("user with name %s found: %w", userName, err)
	}
	s.cfg.SetUser(user.Name)

	fmt.Printf("Logged in as %s\n", user.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <username>", cmd.Name)
	}
	userName := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      userName,
	})
	if err != nil {
		return err
	}
	err = s.cfg.SetUser(user.Name)

	fmt.Printf("Logged in as %s\n", cmd.Args[0])

	return err
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to get users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name)
			continue
		}
		fmt.Printf("%s\n", user.Name)
	}

	return nil
}

package main

import (
	"context"
	"errors"
	"fmt"
	"github/ansht2000/BlogAggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("login command needs a username")
	}

	name := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return errors.New("user does not exist")
	}

	if err = s.cfg.SetUser(name); err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("register command needs a name")
	}
	name := cmd.Args[0]
	
	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return errors.New("user already exists")
	}

	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: name,
	}
	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	if err = s.cfg.SetUser(name); err != nil {
		return err
	}

	fmt.Println("User has been created")
	fmt.Println(user.ID, user.CreatedAt, user.UpdatedAt, user.Name)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list users: %w", err)
	}
	var userString string
	for _, user := range users {
		userString = fmt.Sprintf("* %s", user)
		if s.cfg.CurrentUserName == user {
			userString += " (current)"
		}
		fmt.Println(userString)
	}
	return nil
}
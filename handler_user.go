package main

import (
	"fmt"
	"errors"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return errors.New("login command needs a username")
	}
	err := s.cfg.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}
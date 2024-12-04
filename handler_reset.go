package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("users not deleted: %v", err)
	}
	fmt.Println("Database reset successful")
	return nil
}
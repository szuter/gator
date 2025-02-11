package main

import (
	"context"
	"fmt"
)

func resetUsers(s *state, _ command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset users: %w", err)
	}
	fmt.Println("Users reset")
	return nil
}

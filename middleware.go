package main

import (
	"context"

	"github.com/gminton07/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	
	return func(s *state, cmd command) error {
		// Get User from database
		user, err := s.db.GetUser(
			context.Background(),
			s.cfg.CurrentUserName,
		)
		if err != nil {
			return err
		}

		// Call original handler func
		return handler(s, cmd, user)
	}
}

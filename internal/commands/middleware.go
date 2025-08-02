package commands

import (
	"context"
	"fmt"

	"github.com/agialias-dev/gator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Database.GetUserByName(context.Background(), s.Current_config.User)
		if err != nil {
			return fmt.Errorf("error retrieving user: %v", err)
		}
		err = handler(s, cmd, user)
		return err
	}
}

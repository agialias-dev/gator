package commands

import (
	"fmt"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login command requires a username argument")
	}
	username := cmd.Args[0]
	err := s.Current_config.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}
	fmt.Printf("User set to %s\n", username)
	return nil
}

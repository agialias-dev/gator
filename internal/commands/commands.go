package commands

import (
	"fmt"

	"github.com/agialias-dev/gator/internal/config"
	"github.com/agialias-dev/gator/internal/database"
)

type State struct {
	Database       *database.Queries
	Current_config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Command map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, exists := c.Command[cmd.Name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("command %s not found", cmd.Name)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Command[name] = f
}

package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/agialias-dev/gator/internal/database"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login command requires a username argument")
	}
	username := cmd.Args[0]

	if _, err := s.Database.GetUser(context.Background(), sql.NullString{String: username, Valid: true}); err != nil {
		if err == sql.ErrNoRows {
			log.Fatalf("User '%s' does not exist", username)
		}
	}

	err := s.Current_config.SetUser(username)
	if err != nil {
		return fmt.Errorf("failed to set user: %v", err)
	}
	fmt.Printf("User set to %s\n", username)
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("register command requires a username argument")
	}
	context := context.Background()
	id := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()
	name := cmd.Args[0]
	if _, err := s.Database.GetUser(context, sql.NullString{String: name, Valid: true}); err == nil {
		log.Fatalf("User '%s' already exists", name)
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	s.Database.CreateUser(context, database.CreateUserParams{
		ID:        id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      sql.NullString{String: name, Valid: true},
	})

	s.Current_config.SetUser(name)

	current_user, _ := s.Database.GetUser(context, sql.NullString{String: name, Valid: true})

	fmt.Printf("User '%s' created\n", name)
	printUser(current_user)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
	fmt.Printf(" * Created: %v\n", user.CreatedAt)
	fmt.Printf(" * Updated: %v\n", user.UpdatedAt)
}

func HandlerReset(s *State, cmd Command) error {
	/*fmt.Println("Are you sure you want to reset the database? This will delete all users and their data.")
	fmt.Println("Type 'yes' to confirm:")
	var scan string
	fmt.Scanln(&scan)
	if scan != "yes" {
		fmt.Println("Reset cancelled.")
		return nil
	}*/
	s.Database.ResetUsers(context.Background())
	fmt.Println("Database reset successfully.")
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	users, err := s.Database.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving users: %v", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found.")
		return nil
	}

	fmt.Println("Users:")
	for _, user := range users {
		if user.Name.String == s.Current_config.User {
			fmt.Printf(" * %s (current)\n", user.Name.String)
		} else {
			fmt.Printf(" * %s\n", user.Name.String)
		}
	}
	return nil
}

package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/agialias-dev/gator/internal/database"
	"github.com/agialias-dev/gator/internal/rss"
	"github.com/google/uuid"
)

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("login command requires a username argument")
	}
	username := cmd.Args[0]

	if _, err := s.Database.GetUserByName(context.Background(), username); err != nil {
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
	if _, err := s.Database.GetUserByName(context, name); err == nil {
		log.Fatalf("User '%s' already exists", name)
	} else if err != sql.ErrNoRows {
		return fmt.Errorf("error checking user existence: %v", err)
	}

	s.Database.CreateUser(context, database.CreateUserParams{
		ID:        id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name:      name,
	})

	s.Current_config.SetUser(name)

	current_user, _ := s.Database.GetUserByName(context, name)

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
		if user.Name == s.Current_config.User {
			fmt.Printf(" * %s (current)\n", user.Name)
		} else {
			fmt.Printf(" * %s\n", user.Name)
		}
	}
	return nil
}

func HandlerAggregate(s *State, cmd Command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching RSS feed: %v", err)
	}

	fmt.Printf("Title: %s\n", feed.Channel.Title)
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Printf("Items: %d\n", len(feed.Channel.Item))
	for _, item := range feed.Channel.Item {
		fmt.Printf(" * Title: %s\n", item.Title)
		fmt.Printf(" * Link: %s\n", item.Link)
		fmt.Printf(" * Description: %s\n", item.Description)
	}
	return nil
}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		log.Fatal("addfeed command requires a name and url argument")
	}

	feed, err := s.Database.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}

	_, err = s.Database.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}

	fmt.Printf("Feed '%s' successfully followed\n", cmd.Args[0])
	return nil
}

func HandlerListFeeds(s *State, cmd Command) error {
	feeds, err := s.Database.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Println("Feeds:")
	for _, feed := range feeds {
		name, err := s.Database.FindUserName(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error retrieving user name: %v", err)
		}
		fmt.Printf(" * Name: %s\n", feed.Name)
		fmt.Printf(" * URL: %s\n", feed.Url)
		fmt.Printf(" * User: %s\n", name)
	}
	return nil
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		log.Fatal("follow command requires a feed URL argument")
	}

	feed, err := s.Database.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error retrieving feed: %v", err)
	}

	_, err = s.Database.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %v", err)
	}

	fmt.Printf("%s successfully followed %s\n", s.Current_config.User, cmd.Args[0])
	return nil
}

func HandlerFollowing(s *State, cmd Command, user database.User) error {

	follows, err := s.Database.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error retrieving feed follows: %v", err)
	}

	if len(follows) == 0 {
		fmt.Println("No feeds followed.")
		return nil
	}

	fmt.Printf("%s is following:\n", user.Name)
	for _, follow := range follows {
		feed, err := s.Database.GetFeedById(context.Background(), follow.FeedID)
		if err != nil {
			return fmt.Errorf("error retrieving feed: %v", err)
		}
		fmt.Printf(" * %s\n", feed.Name)
	}
	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		log.Fatal("unfollow command requires a feed URL argument")
	}

	feed, err := s.Database.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error retrieving feed: %v", err)
	}

	err = s.Database.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %v", err)
	}

	fmt.Printf("%s successfully unfollowed %s\n", s.Current_config.User, cmd.Args[0])
	return nil
}

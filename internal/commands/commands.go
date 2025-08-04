package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/agialias-dev/gator/internal/config"
	"github.com/agialias-dev/gator/internal/database"
	"github.com/agialias-dev/gator/internal/rss"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func (s *State) ScrapeFeeds() error {
	nextFeed, err := s.Database.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	if nextFeed == (database.Feed{}) {
		return fmt.Errorf("no feeds available to fetch")
	}

	if err := s.Database.MarkFeedFetched(context.Background(), nextFeed.ID); err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	for _, item := range feed.Channel.Item {
		dateForm := "Mon, 2 Jan 2006 15:04:05 MST"
		pTime, err := time.Parse(dateForm, item.PubDate)
		if err != nil {
			log.Printf("error while parsing date: %v", err)
		}
		if _, err := s.Database.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: pTime,
			FeedID:      nextFeed.ID,
		}); err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Message == "duplicate key value violates unique constraint \"posts_url_key\"" {
				continue
			}
			log.Printf("error while creating post: %v", err)
		}
	}
	return nil
}

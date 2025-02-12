package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerAgg(s *state, _ command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("add-feed requires two arguments")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	rssFeed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to add feed: %w", err)
	}
	fmt.Printf("Feed added:\n ID: %v\n Name: %s\n URL: %s\n UserID: %v\n CreatedAt: %s\n UpdatedAt: %s\n", rssFeed.ID, rssFeed.Name, rssFeed.Url, rssFeed.UserID.UUID, rssFeed.CreatedAt, rssFeed.UpdatedAt)
	return nil
}

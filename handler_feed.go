package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

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
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to add feed: %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    rssFeed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}
	fmt.Printf("Feed added:\n ID: %v\n Name: %s\n URL: %s\n UserID: %v\n CreatedAt: %s\n UpdatedAt: %s\n", rssFeed.ID, rssFeed.Name, rssFeed.Url, rssFeed.UserID, rssFeed.CreatedAt, rssFeed.UpdatedAt)
	return nil
}

func handlerListFeeds(s *state, _ command) error {
	rssFeeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}
	for _, rssFeed := range rssFeeds {

		user, err := s.db.GetUserByID(context.Background(), rssFeed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		printFeed(rssFeed, user.Name)
	}
	return nil
}

func printFeed(feed database.Feed, userName string) {
	fmt.Printf("%s | %s | %s\n", userName, feed.Name, feed.Url)
}

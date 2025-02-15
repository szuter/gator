package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/database"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("agg requires exactly one argument")
	}
	time_between_reqs := cmd.args[0]
	duration, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	fmt.Printf("Collecting feeds every %s\n", time_between_reqs)
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			log.Printf("failed to scrape feeds: %v", err)
		}
	}
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int
	var err error
	if len(cmd.args) == 1 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("failed to parse limit: %w", err)
		}
		if limit < 1 {
			return fmt.Errorf("limit must be greater than 0")
		}
	}
	if len(cmd.args) > 1 {
		return fmt.Errorf("browse requires at most one argument")
	}
	if len(cmd.args) == 0 {
		limit = 2
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed to get posts: %w", err)
	}
	if len(posts) == 0 {
		fmt.Println("No posts found!")
		return nil
	}
	for _, post := range posts {
		fmt.Printf("%s\n%s\n\n", post.Title, post.Url)
		fmt.Printf("Published at: %v\n", post.PublishedAt)
		if post.Description.Valid {
			fmt.Printf("%s\n", post.Description.String)
		}
		fmt.Println("==============================")
	}
	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed: %w", err)
	}
	_, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("failed to mark feed fetched: %w", err)
	}
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}
	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		parsedLink := parseURL(feed.Url, item.Link)
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			FeedID:      feed.ID,
			Title:       item.Title,
			Url:         parsedLink,
			Description: sql.NullString{String: item.Description, Valid: item.Description != ""},
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				continue
			}
			log.Printf("failed to create post: %v", err)
		}
	}
	return nil
}

func parseURL(baseURL, linkURL string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return linkURL
	}
	rel, err := url.Parse(linkURL)
	if err != nil {
		return linkURL
	}
	return base.ResolveReference(rel).String()
}

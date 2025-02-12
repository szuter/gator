package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var feed RSSFeed
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	title := html.UnescapeString(feed.Channel.Title)
	description := html.UnescapeString(feed.Channel.Description)
	feed.Channel.Title = title
	feed.Channel.Description = description
	for i := range feed.Channel.Item {
		title := html.UnescapeString(feed.Channel.Item[i].Title)
		description := html.UnescapeString(feed.Channel.Item[i].Description)
		feed.Channel.Item[i].Title = title
		feed.Channel.Item[i].Description = description
	}
	return &feed, nil

}

package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, _ command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", rssFeed)
	return nil
}

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ansht2000/BlogAggregator/internal/rss"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}
	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}
	fmt.Println(feed.Channel.Title)
	// for _, item := range feed.Channel.Item {
	// 	fmt.Println(item.Title)
	// }
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}
  
	duration, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", duration)
	ticker := time.NewTicker(duration)
	for ;; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}
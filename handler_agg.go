package main

import (
	"context"
	"fmt"
	"github.com/ansht2000/BlogAggregator/internal/rss"
)

func handlerAgg(s *state, cmd command) error {
	feedUrl := "https://www.wagslane.dev/index.xml"
	var feed *rss.RSSFeed
	feed, err := rss.FetchFeed(context.Background(), feedUrl)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}
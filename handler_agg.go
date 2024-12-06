package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ansht2000/BlogAggregator/internal/database"
	"github.com/ansht2000/BlogAggregator/internal/rss"
	"github.com/google/uuid"
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
	for _, item := range feed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time: t,
				Valid: true,
			}
		}
		createPostParams := database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title: item.Title,
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID: nextFeed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), createPostParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", nextFeed.Name, len(feed.Channel.Item))
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
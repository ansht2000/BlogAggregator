package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ansht2000/BlogAggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]
	feed, err := s.db.GetFeedFromURL(context.Background(), url)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:  user.ID,
		FeedID: feed.ID,
	}
	feedRow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println("Feed follow created:")
	printFeedFollow(feedRow.UserName, feedRow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(followedFeeds) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}


	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, followedFeed := range followedFeeds {
		fmt.Printf("* %s\n", followedFeed.FeedName)
	}
	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]
	deleteFeedFollowParams := database.DeleteFeedFollowByUserFeedURLParams{
		UserID: user.ID,
		Url: url,
	}
	feedName, err := s.db.DeleteFeedFollowByUserFeedURL(context.Background(), deleteFeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println("Feed unfollowed successfully:")
	printFeedFollow(user.Name, feedName)
	return nil
}
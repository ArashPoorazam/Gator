package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ArashPoorazam/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollowFeed(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("command needs at least one input. the URL.")
	}

	feed, err := s.Queries.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("the feed does not exist: %w", err)
	}

	user, err := s.Queries.GetUser(context.Background(), s.Config.Current_user_name)
	if err != nil {
		return fmt.Errorf("you have not registered: %w", err)
	}

	followFeed := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	followed, err := s.Queries.CreateFeedFollow(context.Background(), followFeed)
	if err != nil {
		return fmt.Errorf("could not follow this feed: %w", err)
	}

	fmt.Println("------------------------------------------")
	fmt.Printf("User %s followed %s\n", followed.UserName, followed.FeedName)
	fmt.Println("------------------------------------------")

	return nil
}

func handlerUserFollows(s *state, cmd command) error {
	user, err := s.Queries.GetUser(context.Background(), s.Config.Current_user_name)
	if err != nil {
		return fmt.Errorf("please register first: %w", err)
	}

	allFeeds, err := s.Queries.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not catch all feeds: %w", err)
	}

	for i, feed := range allFeeds {
		fmt.Println(feed.FeedName)
		if i >= 50 {
			return nil
		}
	}
	return nil
}

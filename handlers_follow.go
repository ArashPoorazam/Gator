package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ArashPoorazam/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("command needs at least one input. the URL.")
	}

	feed, err := s.Queries.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("the feed does not exist: %w", err)
	}

	err = funcFollowFeed(s, user.ID, feed.ID)
	if err != nil {
		return err
	}

	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("command needs at least one input. the URL.")
	}

	feed, err := s.Queries.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("the feed does not exist: %w", err)
	}

	unfollowParams := database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.Queries.UnfollowFeed(context.Background(), unfollowParams)
	if err != nil {
		return fmt.Errorf("could not unfollow this feed: %w", err)
	}

	fmt.Println("------------------------------------------")
	fmt.Printf("User %s unfollowed %s\n", user.Name, feed.Name)
	fmt.Println("------------------------------------------")

	return nil
}

func handlerUserFollowings(s *state, cmd command, user database.User) error {
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

func funcFollowFeed(s *state, userID uuid.UUID, feedID uuid.UUID) error {
	followFeedParam := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    userID,
		FeedID:    feedID,
	}

	followFeed, err := s.Queries.CreateFeedFollow(context.Background(), followFeedParam)
	if err != nil {
		return fmt.Errorf("could not follow this feed: %w", err)
	}

	fmt.Println("------------------------------------------")
	fmt.Printf("User %s followed %s\n", followFeed.UserName, followFeed.FeedName)
	fmt.Println("------------------------------------------")

	return nil
}

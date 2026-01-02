package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/ArashPoorazam/Gator/internal/database"
	"github.com/ArashPoorazam/Gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command, user database.User) error {
	url := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("feed does not exist, %w", err)
	}

	fmt.Printf("%+v\n", feed)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("the login handler expects a -Two- argument, the -name- of feed and -url-")
	}

	_, err := url.ParseRequestURI(cmd.Args[1])

	feedPram := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	feed, err := s.Queries.CreateFeed(context.Background(), feedPram)
	if err != nil {
		return fmt.Errorf("error on creating the feed: %w", err)
	}

	fmt.Println("New Feed Created!")
	println("=======================================================")
	printFeed(feed, user.Name)

	err = funcFollowFeed(s, user.ID, feed.ID)
	if err != nil {
		return err
	}

	return nil
}

func handlerGetAllFeeds(s *state, cmd command, user database.User) error {
	feeds, err := s.Queries.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not catch all feeds: %w", err)
	}

	for _, feed := range feeds {
		name, err := s.Queries.GetUserName(context.Background(), user.ID)
		if err != nil {
			return err
		}
		println("=======================================================")
		printFeed(feed, name)
	}

	return nil
}

func printFeed(feed database.Feed, name string) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:        	 %s\n", name)
}

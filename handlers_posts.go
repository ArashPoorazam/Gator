package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ArashPoorazam/Gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("you should use it like this -> Gator browse 'URL' '(optional) number of posts'")
	}

	limit := 2
	var err error

	if len(cmd.Args) > 1 {
		limit, err = strconv.Atoi(cmd.Args[1])
		if err != nil {
			return fmt.Errorf("make sure the second argument is a number: %w", err)
		}
	}

	feed, err := s.Queries.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("bad url or the feed does not exist: %w", err)
	}

	getPostParams := database.GetPostsForUserParams{
		FeedID: feed.ID,
		Limit:  int32(limit),
	}

	posts, err := s.Queries.GetPostsForUser(context.Background(), getPostParams)
	if err != nil {
		return fmt.Errorf("error while retriving posts: %w", err)
	}

	for _, post := range posts {
		printpost(post)
	}

	return nil
}

func handlerClearPosts(s *state, cmd command, user database.User) error {
	err := s.Queries.ClearPosts(context.Background())
	if err != nil {
		return fmt.Errorf("could not clear all posts: %w", err)
	}

	return nil
}

func printpost(post database.Post) {
	fmt.Printf("* URL:           %s\n", post.Url)
	fmt.Printf("* Published:     %v\n", post.PublishedAt)
	fmt.Printf("* Title:         %s\n", post.Title)
	fmt.Println("------------------------------------------")
	if post.Description.Valid && post.Description.String != "" {
		fmt.Printf("Description:\n%s\n", post.Description.String)
	} else {
		fmt.Println("Description:\n(No description provided)")
	}
	fmt.Println("------------------------------------------")
}

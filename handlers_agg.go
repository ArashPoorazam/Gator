package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ArashPoorazam/Gator/internal/database"
	"github.com/ArashPoorazam/Gator/internal/rss"
)

var availableReqRate = []string{"1s", "2s", "3s", "5s", "10s", "30s", "1m", "2m", "5m", "10m"}

func handlerAgg(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("need at least -one- argument. the time between requests!")
	}

	timeRate := cmd.Args[0]
	if !isValidRate(timeRate) {
		times := strings.Join(availableReqRate, ", ")

		return fmt.Errorf("try again! input must be one of: \n%s\n", times)
	}

	parsedTime, err := time.ParseDuration(timeRate)
	if err != nil {
		return fmt.Errorf("invalid time input: %w", err)
	}

	ticker := time.NewTicker(parsedTime)
	fmt.Printf("Collecting feeds every %s.\n", timeRate)
	println("=======================================================")
	for range ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping: %v\n", err)
			continue
		}
	}

	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.Queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not find any feed in database: %w", err)
	}

	feedFetchedParams := database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: time.Now().UTC(),
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true, // This tells the DB the value is NOT NULL
		},
	}

	err = s.Queries.MarkFeedFetched(context.Background(), feedFetchedParams)
	if err != nil {
		return fmt.Errorf("could not mark the feed as fetched: %w", err)
	}

	RssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch the new feeds: %w", err)
	}

	for _, item := range RssFeed.Channel.Items {
		fmt.Println(item.Title)
	}

	// for i := range RssFeed.Channel.Items {
	// 	fmt.Println(RssFeed.Channel.Items[i].Title)
	// }

	return nil
}

func isValidRate(rate string) bool {
	for _, v := range availableReqRate {
		if v == rate {
			return true
		}
	}
	return false
}

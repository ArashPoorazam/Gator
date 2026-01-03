package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ArashPoorazam/Gator/internal/database"
	"github.com/ArashPoorazam/Gator/internal/rss"
	"github.com/google/uuid"
)

var availableReqRate = []string{"1s", "2s", "3s", "5s", "10s", "30s", "1m", "2m", "5m", "10m"}

func handlerAgg(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		times := strings.Join(availableReqRate, ", ")

		return fmt.Errorf("you should use it like this -> Gator agg 'tickrate'\ntickrates: %s\n", times)
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
	defer ticker.Stop()

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

	fmt.Printf("scraping: %s\n", feed.Name)
	println("=======================================================")

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

	for _, rssItem := range RssFeed.Channel.Items {

		// saving the posts
		err := savePost(s, &rssItem, feed)

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") ||
				strings.Contains(err.Error(), "unique constraint") {
				continue
			}
			fmt.Printf("Couldn't save post '%s': %v\n", rssItem.Title, err)
		}
	}

	return nil
}

func savePost(s *state, rssItem *rss.RSSItem, feed database.Feed) error {
	// description for NullString in sql database
	description := sql.NullString{
		String: rssItem.Description,
		Valid:  rssItem.Description != "",
	}

	// turn string time to a valid type of time
	pubAt := parsePublishedAt(rssItem.PubDate)

	postParams := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		PublishedAt: pubAt,
		Title:       rssItem.Title,
		Url:         rssItem.Link,
		Description: description,
		FeedID:      feed.ID,
	}

	err := s.Queries.CreatePost(context.Background(), postParams)
	if err != nil {
		return err
	}

	return nil
}

func parsePublishedAt(dateStr string) time.Time {
	// Try the two most common RSS date formats
	formats := []string{time.RFC1123Z, time.RFC1123, time.RFC822, time.RFC822Z}

	for _, layout := range formats {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t
		}
	}

	return time.Now().UTC()
}

func isValidRate(rate string) bool {
	for _, v := range availableReqRate {
		if v == rate {
			return true
		}
	}
	return false
}

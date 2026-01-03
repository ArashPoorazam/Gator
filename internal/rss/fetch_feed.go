package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not make new request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bad request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading request body: %w", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, fmt.Errorf("feed unmarshal error: %w", err)
	}

	cleanFeed(&feed)

	return &feed, nil
}

func cleanFeed(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(stripHTML(feed.Channel.Title))
	feed.Channel.Description = html.UnescapeString(stripHTML(feed.Channel.Description))

	for i := range feed.Channel.Items {
		item := &feed.Channel.Items[i]
		item.Title = html.UnescapeString(stripHTML(item.Title))
		item.Description = html.UnescapeString(stripHTML(item.Description))
	}
}

func stripHTML(src string) string {
	// This regex looks for anything between < and >
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(src, "")
}

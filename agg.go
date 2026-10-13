package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gminton07/gator/internal/database"
)

func scrapeFeeds(s *state) error {
	// Get next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	// Mark as fetched
	err = s.db.MarkFeedFetched(
		context.Background(),
		database.MarkFeedFetchedParams{
			ID:        feed.ID,
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		return err
	}

	// Fetch the feed data
	rssFeed, err := fetchFeed(
		context.Background(),
		feed.Url,
	)
	if err != nil {
		return err
	}

	// Iterate and print out items
	fmt.Printf("Channel title: %s\n", rssFeed.Channel.Title)

	for _, v := range rssFeed.Channel.Item {
		fmt.Printf("\tItem title: %s\n", v.Title)
	}

	fmt.Println()

	return nil
}

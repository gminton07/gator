package main

import (
	"context"
	"database/sql"
	//"errors"
	"fmt"
	"strings"
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
		return fmt.Errorf("error with %s: %w", feed.Url, err)
	}

	// Iterate and print out items
	fmt.Printf("%s\n", rssFeed.Channel.Title)

	for _, v := range rssFeed.Channel.Item {
		fmt.Printf("\t%s\n", v.Title)

		// Save post to db
		currTime := time.Now()

		pubTime, err := parsePublishTime(v.PubDate)
		if err != nil {
			return err
		}

		_, err = s.db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				CreatedAt: currTime,
				UpdatedAt: currTime,
				Title: v.Title,
				Url: v.Link,
				Description: sql.NullString{
					String: v.Description,
					Valid: true,
				},
				PublishedAt: pubTime,
				FeedID: feed.ID,
			},
		)
		if err != nil {
			if dupe := strings.Contains(
				fmt.Sprintf("%w", err), 
				"violates unique",
			); dupe {
				continue
			} else {
				return fmt.Errorf("error with feed.rssFeed.Channel.Title: %w", err)
			}
		}
	}
	fmt.Println()

	return nil
}

func parsePublishTime(value string) (sql.NullTime, error) {
	nullTime := sql.NullTime{}

	//fmt.Printf("Time: %s\n", value)

	// Parse time string into time.Time
	v, err := time.Parse(time.RFC1123Z, value)
	if err != nil {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}, err
	}

	// Set nullTime fields
	nullTime.Time = v
	nullTime.Valid = true

	return nullTime, nil 

}

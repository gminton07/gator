package main

import (
	"context"
	"encoding/xml"
	"html"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create http request
	req, err := http.NewRequestWithContext(ctx,"GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")

	// Create client & make request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	// Decode xml data
	var rssFeed RSSFeed
	decoder := xml.NewDecoder(res.Body)
	err = decoder.Decode(&rssFeed)
	if err != nil {
		return &RSSFeed{}, err
	}

	// Unescape Title/Description strings
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, v := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(v.Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(v.Description)
	}

	return &rssFeed, nil
}

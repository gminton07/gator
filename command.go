package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/gminton07/gator/internal/config"
	"github.com/gminton07/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handle map[string]func(*state, command) error
}

// commands struct methods
func (c *commands) run(s *state, cmd command) error {
	// Run the given command
	err := c.handle[cmd.name](s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	// Register the handler function
	c.handle[name] = f
}

// Basic handler function signature
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("error: Login handler must be given 'user' argument")
	}

	// Poll from database
	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Printf("error in GetUser: %w\n", err)
		os.Exit(1)
	}

	// Set new username
	s.cfg.SetUser(name)

	fmt.Printf("User %s logged in.\n", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("error: Register handler must be given 'user' argument")
	}

	// Add user to database
	name := cmd.args[0]
	currTime := time.Now()
	usr, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      name,
	})
	if err != nil {
		fmt.Printf("error in CreateUser: %w\n", err)
		os.Exit(1)
	}

	// Set current user
	s.cfg.SetUser(name)

	fmt.Printf("User %s created\n", name)
	fmt.Printf("Data: %+v\n", usr)
	return nil
}

func handlerReset(s *state, cmd command) error {
	// Ignore extra args

	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Printf("error in Reset: %w\n", err)
		os.Exit(1)
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Printf("error in GetUsers: %w\n, err")
	}

	currUser := s.cfg.CurrentUserName
	for _, v := range users {
		fmt.Printf("* %s", v.Name)
		if currUser == v.Name {
			fmt.Printf(" (current)\n")
		} else {
			fmt.Println()
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("RSSFeed: %+v\n", rssFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	// parse cli args
	if len(cmd.args) < 2 {
		return errors.New("error: addfeed handler must be given 'name' and 'url' arguments")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
		
	// Create feed
	currTime := time.Now()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %+v\n", feed)

	// Create feed_follow entry
	err = handlerFollow(
		s,
		command{args: []string{feedURL}},
		user,
	)
	if err != nil {
		return err
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	// Ignore CLI args

	// Call SQL Library
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Printf("error in GetFeeds: %w\n", err)
		os.Exit(1)
	}

	for _, v := range feeds {
		fmt.Printf("%s\n\tURL: %s, By: %s\n", v.Name, v.Url, v.Name_2)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	// parse cli args
	if len(cmd.args) < 1 {
		return errors.New("error: follow handler must be given 'url' argument")
	}
	url := cmd.args[0]

	// Get feed data from database
	feed, err := s.db.GetFeed(
		context.Background(),
		url,
	)
	if err != nil {
		fmt.Printf("error in GetFeed: %s\n", err)
		os.Exit(1)
	}
	fdID := feed.ID
	fmt.Printf("feed_id: %v\n", fdID)

	// Call SQL Library
	currTime := time.Now()
	feedFollow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			CreatedAt: currTime,
			UpdatedAt: currTime,
			UserID:    user.ID,
			FeedID:    fdID,
		},
	)

	// Print output
	//fmt.Printf("Feed: %+v\nFeed name: %s\nUser name: %s\n", feedFollow, feed.Name, usrName)
	fmt.Printf("FeedFollow: %+v\n", feedFollow)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	// Ignore args

	feedFollows, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		user.Name)
	if err != nil {
		return err
	}

	// Print out feedFollows
	for _, v := range feedFollows {
		fmt.Printf("%s (URL: %s)\n", v.FeedName, v.FeedUrl)
	}

	return nil
}

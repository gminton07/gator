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

func handlerAddFeed(s *state, cmd command) error {
	// Parse CLI args
	if len(cmd.args) < 2 {
		return errors.New("error: AddFeed handler must be given 'name' and 'url' arguments")
	}
	feedName := cmd.args[0]
	feedURL := cmd.args[1]
		
	// Get user data from database
	// Get username from config
	name := s.cfg.CurrentUserName

	usr, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Printf("error in GetUser: %w", err)
		os.Exit(1)
	}
	usrID := usr.ID
	
	// Create feed
	currTime := time.Now()
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		CreatedAt: currTime,
		UpdatedAt: currTime,
		Name:      feedName,
		Url:       feedURL,
		UserID:    usrID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

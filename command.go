package main

import (
	"errors"
	"fmt"

	"github.com/gminton07/gator/internal/config"
)

type state struct {
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
	args := cmd.args

	if len(args) == 0 {
		return errors.New("Login func must be given 'user' argument")
	}

	// Set new username
	s.cfg.SetUser(args[0])

	fmt.Printf("User %s added.\n", args)
	return nil
}

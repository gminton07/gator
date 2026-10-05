package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gminton07/gator/internal/config"
	"github.com/gminton07/gator/internal/database"
	_ "github.com/lib/pq"
)

// convert all context.Background() instances to single global instance

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.DbURL)

	// Get queries
	dbQueries := database.New(db)

	sta := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// Command registry
	cmds := commands{
		handle: make(map[string]func(*state, command) error),
	}

	cmds.register("login",    handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset",    handlerReset)
	cmds.register("users",    handlerUsers)
	cmds.register("agg",      handlerAgg)

	// Parse cli args
	if len(os.Args) < 2 {
		fmt.Println("Program requires CLI command")
		os.Exit(1)
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(&sta, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

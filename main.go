package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gminton07/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Old config: %+v\n", cfg)

	sta := state{
		cfg: &cfg,
	}

	cmds := commands{
		handle: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)

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

	fmt.Printf("New config: %+v\n", cfg)
}

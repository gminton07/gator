package main

import (
	"fmt"
	"log"

	"github.com/gminton07/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	cfg.SetUser("Gabe")

	cfg2, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Config: %v\n", cfg2)



}

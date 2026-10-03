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
	fmt.Printf("Old config: %+v\n", cfg)

	cfg.SetUser("Gabe")

	cfg2, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New config: %+v\n", cfg2)



}

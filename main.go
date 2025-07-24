package main

import (
	"fmt"
	"log"

	"github.com/agialias-dev/gator/internal/config"
)

func main() {
	curr, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Current URL:\t%v\nCurrent User:\t%v\n", curr.URL, curr.User)

	err = curr.SetUser("agialias")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

	new, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("New URL:\t%v\nNew User:\t%v\n", new.URL, new.User)
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/agialias-dev/gator/internal/commands"
	"github.com/agialias-dev/gator/internal/config"
)

func main() {
	var err error
	var state commands.State
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	state.Current_config = &c
	fmt.Printf("Current URL:\t%v\nCurrent User:\t%v\n", state.Current_config.URL, state.Current_config.User)

	clicmds := &commands.Commands{
		Command: map[string]func(*commands.State, commands.Command) error{},
	}

	clicmds.Register("login", commands.HandlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("no arguments provided")
	} else {
		cmd := commands.Command{
			Name: os.Args[1],
			Args: os.Args[2:],
		}
		err = clicmds.Run(&state, cmd)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	fmt.Printf("New URL:\t%v\nNew User:\t%v\n", state.Current_config.URL, state.Current_config.User)
}

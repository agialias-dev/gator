package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/agialias-dev/gator/internal/commands"
	"github.com/agialias-dev/gator/internal/config"
	"github.com/agialias-dev/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	var err error
	var state commands.State
	c, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	state.Current_config = &c

	db, err := sql.Open("postgres", state.Current_config.URL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	state.Database = database.New(db)

	clicmds := &commands.Commands{
		Command: map[string]func(*commands.State, commands.Command) error{},
	}

	clicmds.Register("login", commands.HandlerLogin)
	clicmds.Register("register", commands.HandlerRegister)
	clicmds.Register("users", commands.HandlerUsers)
	clicmds.Register("reset", commands.HandlerReset)
	clicmds.Register("agg", commands.HandlerAggregate)
	clicmds.Register("addfeed", commands.MiddlewareLoggedIn(commands.HandlerAddFeed))
	clicmds.Register("feeds", commands.HandlerListFeeds)
	clicmds.Register("follow", commands.MiddlewareLoggedIn(commands.HandlerFollow))
	clicmds.Register("following", commands.MiddlewareLoggedIn(commands.HandlerFollowing))
	clicmds.Register("unfollow", commands.MiddlewareLoggedIn(commands.HandlerUnfollow))

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
}

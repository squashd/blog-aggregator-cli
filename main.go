package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/squashd/blog-aggregator-cli/internal/config"
	"github.com/squashd/blog-aggregator-cli/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	conn, err := sql.Open("postgres", cfg.DBURL)
	queries := database.New(conn)

	appState := state{cfg: &cfg, db: queries}

	cmds := commands{
		registeredCmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerGetFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if len(os.Args) < 2 {
		fmt.Println("usage: cli <command> [args...]")
		os.Exit(1)
		return
	}

	cmdName := os.Args[1]
	cmdArgs := []string{}
	if len(os.Args) > 2 {
		cmdArgs = os.Args[2:]
	}

	err = cmds.run(&appState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
		os.Exit(1)
	}

}

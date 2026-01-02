package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/ArashPoorazam/Gator/internal/config"
	"github.com/ArashPoorazam/Gator/internal/database"
)

func main() {
	// fmt.Println("Dont forget to build the app first then test it!!!")

	// read config
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// creating state and command handler
	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		log.Fatal("connection error: %w", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("ping error: %w", err)
	}

	dbQueries := database.New(db)

	State := state{
		Config:  &cfg,
		Queries: dbQueries,
	}
	commands := commands{make(map[string]func(*state, command) error)}

	// add commands
	commands.add("login", middlewareLoggedIn(handlerLogin))
	commands.add("register", middlewareLoggedIn(handlerRegister))
	commands.add("delete", middlewareLoggedIn(handlerDeleteUser))
	commands.add("reset", middlewareLoggedIn(handlerResetTable))
	commands.add("users", middlewareLoggedIn(handlerGetUsers))
	commands.add("agg", middlewareLoggedIn(handlerAgg))
	commands.add("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.add("feeds", middlewareLoggedIn(handlerGetAllFeeds))
	commands.add("follow", middlewareLoggedIn(handlerFollowFeed))
	commands.add("unfollow", middlewareLoggedIn(handlerUnfollowFeed))
	commands.add("following", middlewareLoggedIn(handlerUserFollowings))

	// get args
	argsSlice := os.Args[:]
	if len(argsSlice) < 2 {
		log.Fatal("Peogram need at least one argument")
	}

	// run command
	newCommand := command{Name: argsSlice[1], Args: argsSlice[2:]}
	err = commands.run(&State, newCommand)
	if err != nil {
		log.Fatal(err)
	}
}

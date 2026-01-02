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
	commands.add("login", handlerLogin)
	commands.add("register", handlerRegister)
	commands.add("delete", handlerDeleteUser)
	commands.add("reset", handlerResetTable)
	commands.add("users", handlerGetUsers)
	commands.add("agg", handlerAgg)
	commands.add("addfeed", handlerAddFeed)
	commands.add("feeds", handlerGetAllFeeds)
	commands.add("follow", handlerFollowFeed)
	commands.add("following", handlerUserFollows)

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

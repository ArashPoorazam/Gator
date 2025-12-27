package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/ArashPoorazam/Gator/internal/config"
	"github.com/ArashPoorazam/Gator/internal/database"
)

func main() {
	fmt.Println("Dont forget to build the app first then test it!!!")

	// read config
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// creating state and command handler
	query, err := newQuery(cfg.Db_url)
	if err != nil {
		log.Fatal(err)
	}
	State := state{
		Config:  &cfg,
		Queries: query,
	}
	commands := commands{make(map[string]func(*state, command) error)}

	// commands
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)

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

func newQuery(dbLink string) (dbQueries *database.Queries, err error) {
	db, err := sql.Open("postgres", dbLink)
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	dbQueries = database.New(db)

	return dbQueries, nil
}

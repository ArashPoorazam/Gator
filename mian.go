package main

import (
	"fmt"
	"log"

	"github.com/ArashPoorazam/Gator/internal/config"
)

func main() {
	fmt.Println("Hello world!")

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.SetUser("Arash Poorazam")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("config: %+v\n", cfg)
}

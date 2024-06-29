package main

import (
	"cmdiet/commands"
	_ "cmdiet/commands" // This imports the commands package to ensure the init functions are called
	"cmdiet/database"
	"flag"
	"log"
)

func main() {
	fake := flag.Bool("fake", false, "Use this flag to initialize a fake database")
	flag.Parse()

	if *fake {
		database.InitDB("fake.db")
	} else {
		database.InitDB("diet.db")
	}
	defer database.CloseDB()

	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}

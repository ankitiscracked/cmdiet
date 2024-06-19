package main

import (
	"cmdiet/commands"
	_ "cmdiet/commands" // This imports the commands package to ensure the init functions are called
	"cmdiet/database"
	"log"
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}

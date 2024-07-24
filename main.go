package main

import (
	"cmdiet/commands"
	_ "cmdiet/commands" // This imports the commands package to ensure the init functions are called
	"cmdiet/diet"
	"cmdiet/meals"
	"fmt"
	"log"

	flag "github.com/spf13/pflag"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openDatabse(isFake bool) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	if isFake {
		db, err = gorm.Open(sqlite.Open("fake.db"), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open("diet.db"), &gorm.Config{})
	}
	if err != nil {
		return db, fmt.Errorf("couldn't open the database: %w", err)
	}

	err = db.AutoMigrate(&meals.Meal{}, &diet.Diet{})
	if err != nil {
		return db, fmt.Errorf("couldn't migrate the database: %w", err)
	}

	return db, nil
}

func main() {
	fake := flag.Bool("fake", false, "Use this flag to initialize a fake database")
	flag.IntP("offset", "o", 7, "Use this flag to set the offset for viewing meals")
	flag.Parse()

	database, err := openDatabse(*fake)
	if err != nil {
		log.Fatal(err)
	}

	diet.DS = &diet.DietServiceImpl{DB: database}
	meals.MS = &meals.MealServiceImpl{DB: database}

	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}

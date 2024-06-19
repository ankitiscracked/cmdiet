package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./diet.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable(CreateDietTableSql)
	createTable(CreateMealTableSql)
}

func createTable(ddl string) {
	_, err := DB.Exec(ddl)
	if err != nil {
		log.Fatal(err)
	}
}

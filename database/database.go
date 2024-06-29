package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbName string) {
	var err error
	DB, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
		DB.Close()
	}

	CreateTable(CreateDietTableSql)
	CreateTable(CreateMealTableSql)
}

func CreateTable(ddl string) {
	_, err := DB.Exec(ddl)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	DB.Close()
}

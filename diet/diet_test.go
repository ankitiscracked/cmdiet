package diet

import (
	"cmdiet/database"
	"database/sql"
	"testing"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	execDdl(db, database.CreateDietTableSql)
	execDdl(db, database.CreateMealTableSql)

	return db, nil
}

func execDdl(db *sql.DB, ddl string) error {
	_, err := db.Exec(ddl)
	return err
}

func TestErrorLoggingTheSameMeal(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	dietService := NewDietService(db)
	err = dietService.LogDiet("breakfast", "eggs", 200, "home")
	if err != nil {
		t.Fatal(err)
	}
	err = dietService.LogDiet("breakfast", "eggs", 200, "home")
	if err == nil {
		t.Fatal("Expected an error but got nil")
	}
}

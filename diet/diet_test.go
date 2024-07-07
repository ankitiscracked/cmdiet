package diet

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("couldn't open the database: %w", err)
	}
	db.AutoMigrate(&Diet{})
	t.Cleanup(func() {
		db.Migrator().DropTable(&Diet{})
	})
	return db, nil
}

func TestErrorLoggingTheSameMeal(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	dietService := &DietServiceImpl{DB: db}
	err = dietService.LogDiet("breakfast", "eggs", 200, "home")
	if err != nil {
		t.Fatal(err)
	}
	err = dietService.LogDiet("breakfast", "eggs", 200, "home")
	if err == nil {
		t.Fatal("Expected an error but got nil")
	}
}

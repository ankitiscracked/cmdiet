package diet

import (
	"cmdiet/meals"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("couldn't open the database: %w", err)
	}

	db.AutoMigrate(&Diet{})
	db.AutoMigrate(&meals.Meal{})

	t.Cleanup(func() {
		db.Migrator().DropTable(&Diet{})
		db.Migrator().DropTable(&meals.Meal{})
	})
	return db, nil
}

func TestLogRepeatedMeal(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	meals.MS = &meals.MealServiceImpl{DB: db}
	dietService := &DietServiceImpl{DB: db}

	err = dietService.LogDietWithNewMeal(Breakfast, "eggs", 200, "home")
	assert.Nil(t, err)

	err = dietService.LogDietWithNewMeal(Breakfast, "eggs", 200, "home")
	assert.NotNil(t, err)
}

func TestLogWithInvalidMeal(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	meals.MS = &meals.MealServiceImpl{DB: db}
	dietService := &DietServiceImpl{DB: db}

	err = dietService.LogDietWithExistingMeal(1, Breakfast, "home")
	assert.NotNil(t, err)
}

func TestLogInputs(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	dietService := &DietServiceImpl{}
	meals.MS = &meals.MealServiceImpl{DB: db}

	t.Run("should throw for unsupported meal type", func(t *testing.T) {
		err := dietService.LogDietWithNewMeal(0, "eggs", 200, "home")
		assert.NotNil(t, err)
	})

	// t.Run("should throw for unsupported source type", func(t *testing.T) {
	// 	err := dietService.LogDietWithNewMeal("breakfast", "eggs", 200, "home")
	// 	if err != nil {
	// 		t.Fatal("Expected an error but got nil")
	// 	}
	// })
}

func TestFechingDiets(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	dietService := &DietServiceImpl{DB: db}
	diets, err := dietService.GetDietsForDay(time.Now())

	assert.Nil(t, err)
	assert.Equal(t, 0, len(diets))

	dietService.LogDietWithNewMeal(Breakfast, "eggs", 200, "home")
	dietService.LogDietWithNewMeal(Lunch, "eggs", 200, "home")
	dietService.LogDietWithNewMeal(Dinner, "eggs", 200, "home")
	diets, err = dietService.GetDietsForDay(time.Now())

	assert.Nil(t, err)
	assert.Equal(t, 3, len(diets))
}

package meals

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("couldn't open the database: %w", err)
	}
	db.AutoMigrate(&Meal{})
	t.Cleanup(func() {
		db.Migrator().DropTable(&Meal{})
	})
	return db, nil
}

func TestFetchingMeals(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	mealService := &MealServiceImpl{DB: db}

	t.Run("should throw for non existent meal", func(t *testing.T) {
		err, _ := mealService.GetMeal(1)
		assert.NotNil(t, err)
	})

	t.Run("should return no meals", func(t *testing.T) {
		meals, _ := mealService.GetAllMeals()
		assert.Equal(t, len(meals), 0)
	})

	mealService.AddMeal("test", 100)

	t.Run("should return the added meal's name", func(t *testing.T) {
		meal, err := mealService.GetMeal(1)
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, meal.Name, "test")
	})

	t.Run("should return the count", func(t *testing.T) {
		meals, _ := mealService.GetAllMeals()
		assert.Equal(t, len(meals), 1)
	})
}

func TestUpdateMeal(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	mealService := &MealServiceImpl{DB: db}

	t.Run("should throw meal with id is invalid", func(t *testing.T) {
		err = mealService.UpdateMeal(Meal{Id: 1})
		assert.NotNil(t, err)
	})

	t.Run("should update the meal if id is valid", func(t *testing.T) {
		added, _ := mealService.AddMeal("test", 100)
		err := mealService.UpdateMeal(Meal{Id: added.Id, Name: "test1"})
		assert.Nil(t, err)

		meal, _ := mealService.GetMeal(added.Id)
		assert.Equal(t, "test1", meal.Name)
	})
}

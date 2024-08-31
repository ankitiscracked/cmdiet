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

	mealService.AddMeal("test")

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

	t.Run("should return all added meals", func(t *testing.T) {
		mealService.AddMeal("test2")
		mealService.AddMeal("test3")
		meals, err := mealService.GetAllMeals()
		assert.Nil(t, err)
		assert.Equal(t, 3, len(meals))
		assert.Equal(t, "test", meals[0].Name)
		assert.Equal(t, "test2", meals[1].Name)
		assert.Equal(t, "test3", meals[2].Name)
	})

	t.Run("should return correct meal by ID", func(t *testing.T) {
		meal, err := mealService.GetMeal(2)
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, "test2", meal.Name)
	})

	t.Run("should return error for non-existent meal ID", func(t *testing.T) {
		_, err := mealService.GetMeal(999)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "couldn't get meal")
	})

	t.Run("should return empty slice when no meals exist", func(t *testing.T) {
		// Clear existing meals
		db.Exec("DELETE FROM meals")

		meals, err := mealService.GetAllMeals()
		assert.Nil(t, err)
		assert.Equal(t, 0, len(meals))
	})

	t.Run("should handle large number of meals", func(t *testing.T) {
		// Add 100 meals
		for i := 0; i < 100; i++ {
			mealService.AddMeal(fmt.Sprintf("meal%d", i))
		}

		meals, err := mealService.GetAllMeals()
		assert.Nil(t, err)
		assert.Equal(t, 100, len(meals))
	})

}

func TestAddingMeal(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	mealService := &MealServiceImpl{DB: db}

	t.Run("should add meal with fixed components", func(t *testing.T) {
		addedComponent, err := mealService.AddFixedMealComponent(FixedMealComponent{Name: "test", Protein: 30, Carbs: 40, Fat: 50})
		assert.Nil(t, err)
		components := []MealComponent{
			{Type: FixedMealComponentType, Count: 1, Id: int(addedComponent.ID)},
		}
		meal, err := mealService.AddMealWithComponents(MealPayload{Name: "test", Components: components})
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, "test", meal.Name)
		assert.Equal(t, 1, len(meal.FixedComponentData))
		assert.Equal(t, 30+40+50, meal.GetTotalCalories()) // Protein + Carbs + Fat * 4 + 9
	})

	t.Run("should add meal with variable components", func(t *testing.T) {
		addedComponent, err := mealService.AddVariableMealComponent(VariableMealComponent{Name: "test", ProteinPerHundredGram: 30, CarbsPerHundredGram: 40, FatPerHundredGram: 50})
		assert.Nil(t, err)
		components := []MealComponent{
			{Type: VariableMealComponentType, AmountInGrams: 200, Id: int(addedComponent.ID)},
		}
		meal, err := mealService.AddMealWithComponents(MealPayload{Name: "test", Components: components})
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, "test", meal.Name)
		assert.Equal(t, 1, len(meal.VariableComponentData))
		assert.Equal(t, 2*(30+40+50), meal.GetTotalCalories()) // (Protein + Carbs + Fat) * 2 * 4 + 9
	})

	t.Run("should add meal with fixed and variable components", func(t *testing.T) {
		fixedComponent, err := mealService.AddFixedMealComponent(FixedMealComponent{Name: "fixed", Protein: 10, Carbs: 20, Fat: 30})
		assert.Nil(t, err)
		variableComponent, err := mealService.AddVariableMealComponent(VariableMealComponent{Name: "variable", ProteinPerHundredGram: 5, CarbsPerHundredGram: 10, FatPerHundredGram: 15})
		assert.Nil(t, err)
		components := []MealComponent{
			{Type: FixedMealComponentType, Count: 1, Id: int(fixedComponent.ID)},
			{Type: VariableMealComponentType, AmountInGrams: 200, Id: int(variableComponent.ID)},
		}
		meal, err := mealService.AddMealWithComponents(MealPayload{Name: "test", Components: components})
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, "test", meal.Name)
		assert.Equal(t, 1, len(meal.FixedComponentData))
		assert.Equal(t, 1, len(meal.VariableComponentData))
		expectedCalories := (10 + 20 + 30) + 2*(5+10+15)
		assert.Equal(t, expectedCalories, meal.GetTotalCalories())
	})

	t.Run("should return error for invalid meal payload", func(t *testing.T) {
		invalidPayload := MealPayload{
			Name:       "",
			Components: []MealComponent{},
		}
		_, err := mealService.AddMealWithComponents(invalidPayload)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid meal payload")
	})

	t.Run("should return error for non-existent component", func(t *testing.T) {
		nonExistentComponentPayload := MealPayload{
			Name: "Test Meal",
			Components: []MealComponent{
				{Type: FixedMealComponentType, Count: 1, Id: 9999},
			},
		}
		_, err := mealService.AddMealWithComponents(nonExistentComponentPayload)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "couldn't find fixed component")
	})

	t.Run("should return error for invalid component type", func(t *testing.T) {
		invalidComponentTypePayload := MealPayload{
			Name: "Test Meal",
			Components: []MealComponent{
				{Type: 999, Count: 1, Id: 1},
			},
		}
		_, err := mealService.AddMealWithComponents(invalidComponentTypePayload)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unknown component type")
	})

	t.Run("should add meal with zero components", func(t *testing.T) {
		emptyComponentsPayload := MealPayload{
			Name:       "Empty Meal",
			Components: []MealComponent{},
		}
		meal, err := mealService.AddMealWithComponents(emptyComponentsPayload)
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, "Empty Meal", meal.Name)
		assert.Equal(t, 0, len(meal.FixedComponentData))
		assert.Equal(t, 0, len(meal.VariableComponentData))
		assert.Equal(t, 0, meal.GetTotalCalories())
	})
}

func TestUpdateMeal(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	mealService := &MealServiceImpl{DB: db}

	t.Run("should throw meal with id is invalid", func(t *testing.T) {
		_, err = mealService.UpdateMeal(1, MealPayload{})
		assert.NotNil(t, err)
	})

	t.Run("should update the meal if id is valid", func(t *testing.T) {
		added, _ := mealService.AddMeal("test")
		_, err := mealService.UpdateMeal(added.Id, MealPayload{Name: "test1"})
		assert.Nil(t, err)

		meal, _ := mealService.GetMeal(added.Id)
		assert.Equal(t, "test1", meal.Name)
	})

	t.Run("should update total calories", func(t *testing.T) {
		added, _ := mealService.AddMeal("test")

		fixedComponent, _ := mealService.AddFixedMealComponent(FixedMealComponent{
			Name:    "Fixed Test Component",
			Protein: 10,
			Carbs:   20,
			Fat:     5,
		})
		payload := MealPayload{
			Name: "Updated Test Meal",
			Components: []MealComponent{
				{
					Type:  FixedMealComponentType,
					Count: 2,
					Id:    int(fixedComponent.ID),
				},
			},
		}
		_, err := mealService.UpdateMeal(added.Id, payload)
		assert.Nil(t, err)
		meal, _ := mealService.GetMeal(added.Id)
		assert.Equal(t, 330, meal.GetTotalCalories())
	})

	t.Run("should update meal with both fixed and variable components", func(t *testing.T) {
		added, _ := mealService.AddMeal("test")

		fixedComponent, _ := mealService.AddFixedMealComponent(FixedMealComponent{
			Name:    "Fixed Test Component",
			Protein: 10,
			Carbs:   20,
			Fat:     5,
		})

		variableComponent, _ := mealService.AddVariableMealComponent(VariableMealComponent{
			Name:                  "Variable Test Component",
			ProteinPerHundredGram: 15,
			CarbsPerHundredGram:   25,
			FatPerHundredGram:     8,
		})

		payload := MealPayload{
			Name: "Updated Test Meal",
			Components: []MealComponent{
				{
					Type:  FixedMealComponentType,
					Count: 2,
					Id:    int(fixedComponent.ID),
				},
				{
					Type:          VariableMealComponentType,
					AmountInGrams: 150,
					Id:            int(variableComponent.ID),
				},
			},
		}
		_, err := mealService.UpdateMeal(added.Id, payload)
		assert.Nil(t, err)

		meal, _ := mealService.GetMeal(added.Id)
		assert.Equal(t, "Updated Test Meal", meal.Name)
		assert.Equal(t, 2, len(meal.FixedComponentData))
		assert.Equal(t, 1, len(meal.VariableComponentData))
		assert.Equal(t, 561, meal.GetTotalCalories()) // 330 (from fixed) + 231 (from variable)
	})

	t.Run("should return error when updating with non-existent component", func(t *testing.T) {
		added, _ := mealService.AddMeal("test")

		payload := MealPayload{
			Name: "Updated Test Meal",
			Components: []MealComponent{
				{
					Type:  FixedMealComponentType,
					Count: 2,
					Id:    9999, // Non-existent component ID
				},
			},
		}
		_, err := mealService.UpdateMeal(added.Id, payload)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "couldn't find fixed component")
	})

	t.Run("should return error when updating with invalid component type", func(t *testing.T) {
		added, _ := mealService.AddMeal("test")

		payload := MealPayload{
			Name: "Updated Test Meal",
			Components: []MealComponent{
				{
					Type:  999, // Invalid component type
					Count: 2,
					Id:    1,
				},
			},
		}
		_, err := mealService.UpdateMeal(added.Id, payload)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unknown component type")
	})

	t.Run("should update meal macros correctly", func(t *testing.T) {
		added, _ := mealService.AddMeal("test")

		fixedComponent, _ := mealService.AddFixedMealComponent(FixedMealComponent{
			Name:    "Fixed Test Component",
			Protein: 10,
			Carbs:   20,
			Fat:     5,
		})

		payload := MealPayload{
			Name: "Updated Test Meal",
			Components: []MealComponent{
				{
					Type:  FixedMealComponentType,
					Count: 3,
					Id:    int(fixedComponent.ID),
				},
			},
		}
		_, err := mealService.UpdateMeal(added.Id, payload)
		assert.Nil(t, err)

		meal, _ := mealService.GetMeal(added.Id)
		assert.Equal(t, 30, meal.GetMealMacro(Protein))
		assert.Equal(t, 60, meal.GetMealMacro(Carbs))
		assert.Equal(t, 15, meal.GetMealMacro(Fat))
	})
}

func TestMealComponentOperations(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	mealService := &MealServiceImpl{DB: db}

	t.Run("should add a fixed meal component", func(t *testing.T) {
		component := FixedMealComponent{
			Name:    "Test Fixed Component",
			Protein: 10,
			Carbs:   20,
			Fat:     5,
		}
		added, err := mealService.AddFixedMealComponent(component)
		assert.Nil(t, err)
		assert.NotNil(t, added)
		assert.Equal(t, component.Name, added.Name)
		assert.Equal(t, component.Protein, added.Protein)
		assert.Equal(t, component.Carbs, added.Carbs)
		assert.Equal(t, component.Fat, added.Fat)
	})

	t.Run("should add a variable meal component", func(t *testing.T) {
		component := VariableMealComponent{
			Name:                  "Test Variable Component",
			ProteinPerHundredGram: 15,
			CarbsPerHundredGram:   25,
			FatPerHundredGram:     8,
		}
		added, err := mealService.AddVariableMealComponent(component)
		assert.Nil(t, err)
		assert.NotNil(t, added)
		assert.Equal(t, component.Name, added.Name)
		assert.Equal(t, component.ProteinPerHundredGram, added.ProteinPerHundredGram)
		assert.Equal(t, component.CarbsPerHundredGram, added.CarbsPerHundredGram)
		assert.Equal(t, component.FatPerHundredGram, added.FatPerHundredGram)
	})

	t.Run("should not add a fixed meal component with invalid data", func(t *testing.T) {
		component := FixedMealComponent{
			Name:    "",
			Protein: -10,
			Carbs:   20,
			Fat:     5,
		}
		_, err := mealService.AddFixedMealComponent(component)
		assert.NotNil(t, err)
	})

	t.Run("should not add a variable meal component with invalid data", func(t *testing.T) {
		component := VariableMealComponent{
			Name:                  "",
			ProteinPerHundredGram: -15,
			CarbsPerHundredGram:   25,
			FatPerHundredGram:     8,
		}
		_, err := mealService.AddVariableMealComponent(component)
		assert.NotNil(t, err)
	})

	t.Run("should retrieve a fixed meal component", func(t *testing.T) {
		component := FixedMealComponent{
			Name:    "Retrievable Fixed Component",
			Protein: 12,
			Carbs:   22,
			Fat:     7,
		}
		added, _ := mealService.AddFixedMealComponent(component)
		retrieved, err := mealService.GetFixedMealComponent(int(added.ID))
		assert.Nil(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, added.ID, retrieved.ID)
		assert.Equal(t, component.Name, retrieved.Name)
	})

	t.Run("should retrieve a variable meal component", func(t *testing.T) {
		component := VariableMealComponent{
			Name:                  "Retrievable Variable Component",
			ProteinPerHundredGram: 17,
			CarbsPerHundredGram:   27,
			FatPerHundredGram:     9,
		}
		added, _ := mealService.AddVariableMealComponent(component)
		retrieved, err := mealService.GetVariableMealComponent(int(added.ID))
		assert.Nil(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, added.ID, retrieved.ID)
		assert.Equal(t, component.Name, retrieved.Name)
	})

	t.Run("should return error when retrieving non-existent fixed component", func(t *testing.T) {
		_, err := mealService.GetFixedMealComponent(9999)
		assert.NotNil(t, err)
	})

	t.Run("should return error when retrieving non-existent variable component", func(t *testing.T) {
		_, err := mealService.GetVariableMealComponent(9999)
		assert.NotNil(t, err)
	})
}

package meals

import (
	"cmdiet/constants"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: nil})
	if err != nil {
		return db, fmt.Errorf("couldn't open the database: %w", err)
	}
	db.AutoMigrate(&FixedMealComponent{})
	db.AutoMigrate(&VariableMealComponent{})
	db.AutoMigrate(&FixedMealComponentData{})
	db.AutoMigrate(&VariableMealComponentData{})
	db.AutoMigrate(&Meal{})
	t.Cleanup(func() {
		db.Migrator().DropTable(&Meal{})
	})
	return db, nil
}

func setupValidator() {
	constants.Validator = validator.New(validator.WithRequiredStructEnabled())
	constants.Validator.RegisterValidation("mealComponentType", ValidateMealComponentType)
}

func TestFetchingMeals(t *testing.T) {
	setupValidator()
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
	setupValidator()
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	ms := &MealServiceImpl{DB: db}
	mcs := &MealComponentService{DB: db}

	t.Run("should add meal with fixed components", func(t *testing.T) {
		addedComponent, err := mcs.AddFixedMealComponent(FixedMealComponent{Name: "test", Protein: 30, Carbs: 40, Fat: 50})
		assert.Nil(t, err)
		components := []MealComponent{
			{Name: "test", Type: FixedMealComponentType, Count: 1, Id: int(addedComponent.ID)},
		}
		meal, err := ms.AddMealWithComponents(MealPayload{Name: "test", Components: components})
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, "test", meal.Name)
		assert.Equal(t, 1, len(meal.FixedComponentData))
		assert.Equal(t, 30*4+40*4+50*9, meal.GetTotalCalories()) // Protein + Carbs + Fat * 4 + 9
	})

	t.Run("should add meal with variable components", func(t *testing.T) {
		addedComponent, err := mcs.AddVariableMealComponent(VariableMealComponent{Name: "test", ProteinPerHundredGram: 30, CarbsPerHundredGram: 40, FatPerHundredGram: 50})
		assert.Nil(t, err)
		components := []MealComponent{
			{Name: "test", Type: VariableMealComponentType, AmountInGrams: 200, Id: int(addedComponent.ID)},
		}
		meal, err := ms.AddMealWithComponents(MealPayload{Name: "test", Components: components})
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, "test", meal.Name)
		assert.Equal(t, 1, len(meal.VariableComponentData))
		assert.Equal(t, 2*(30*4+40*4+50*9), meal.GetTotalCalories()) // (Protein + Carbs + Fat) * 2 * 4 + 9
	})

	t.Run("should add meal with fixed and variable components", func(t *testing.T) {
		fixedComponent, err := mcs.AddFixedMealComponent(FixedMealComponent{Name: "fixed", Protein: 10, Carbs: 20, Fat: 30})
		assert.Nil(t, err)
		variableComponent, err := mcs.AddVariableMealComponent(VariableMealComponent{Name: "variable", ProteinPerHundredGram: 5, CarbsPerHundredGram: 10, FatPerHundredGram: 15})
		assert.Nil(t, err)
		components := []MealComponent{
			{Name: "test", Type: FixedMealComponentType, Count: 1, Id: int(fixedComponent.ID)},
			{Name: "test", Type: VariableMealComponentType, AmountInGrams: 200, Id: int(variableComponent.ID)},
		}
		meal, err := ms.AddMealWithComponents(MealPayload{Name: "test", Components: components})
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, "test", meal.Name)
		assert.Equal(t, 1, len(meal.FixedComponentData))
		assert.Equal(t, 1, len(meal.VariableComponentData))
		expectedCalories := (10*4 + 20*4 + 30*9) + 2*(5*4+10*4+15*9)
		assert.Equal(t, expectedCalories, meal.GetTotalCalories())
	})

	t.Run("should return error for invalid meal payload", func(t *testing.T) {
		invalidPayload := MealPayload{
			Name:       "",
			Components: []MealComponent{},
		}
		_, err := ms.AddMealWithComponents(invalidPayload)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid meal payload")
	})

	t.Run("should return error for non-existent component", func(t *testing.T) {
		nonExistentComponentPayload := MealPayload{
			Name: "Test Meal",
			Components: []MealComponent{
				{Name: "test", Type: FixedMealComponentType, Count: 1, Id: 9999},
			},
		}
		_, err := ms.AddMealWithComponents(nonExistentComponentPayload)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "couldn't find fixed component")
	})

	t.Run("should return error for invalid component type", func(t *testing.T) {
		invalidComponentTypePayload := MealPayload{
			Name: "Test Meal",
			Components: []MealComponent{
				{Name: "test", Type: "", Count: 1, Id: 1},
			},
		}
		_, err := ms.AddMealWithComponents(invalidComponentTypePayload)
		assert.NotNil(t, err)
	})

	t.Run("should add meal with zero components", func(t *testing.T) {
		emptyComponentsPayload := MealPayload{
			Name:       "Empty Meal",
			Components: []MealComponent{},
		}
		meal, err := ms.AddMealWithComponents(emptyComponentsPayload)
		assert.Nil(t, err)
		assert.NotNil(t, meal)
		assert.Equal(t, "Empty Meal", meal.Name)
		assert.Equal(t, 0, len(meal.FixedComponentData))
		assert.Equal(t, 0, len(meal.VariableComponentData))
		assert.Equal(t, 0, meal.GetTotalCalories())
	})
}

func TestUpdateMeal(t *testing.T) {
	setupValidator()
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	ms := &MealServiceImpl{DB: db}
	mcs := &MealComponentService{DB: db}

	t.Run("should throw meal with id is invalid", func(t *testing.T) {
		_, err = ms.UpdateMeal(1, UpdateMealPayload{})
		assert.NotNil(t, err)
	})

	t.Run("should update the meal if id is valid", func(t *testing.T) {
		added, _ := ms.AddMeal("test")
		_, err := ms.UpdateMeal(int(added.ID), UpdateMealPayload{Name: "test1"})
		assert.Nil(t, err)

		meal, _ := ms.GetMeal(int(added.ID))
		assert.Equal(t, "test1", meal.Name)
	})

	t.Run("should update total calories", func(t *testing.T) {
		added, _ := ms.AddMeal("test")

		fixedComponent, _ := mcs.AddFixedMealComponent(FixedMealComponent{
			Name:    "Fixed Test Component",
			Protein: 10,
			Carbs:   20,
			Fat:     5,
		})
		payload := UpdateMealPayload{
			Name: "Updated Test Meal",
			Components: []MealComponent{
				{
					Type:  FixedMealComponentType,
					Count: 2,
					Id:    int(fixedComponent.ID),
				},
			},
		}
		_, err := ms.UpdateMeal(int(added.ID), payload)
		assert.Nil(t, err)
		meal, _ := ms.GetMeal(int(added.ID))
		assert.Equal(t, 330, meal.GetTotalCalories())
	})

	t.Run("should update meal with both fixed and variable components", func(t *testing.T) {
		added, _ := ms.AddMeal("test")

		fixedComponent, _ := mcs.AddFixedMealComponent(FixedMealComponent{
			Name:    "Fixed Test Component",
			Protein: 10,
			Carbs:   20,
			Fat:     5,
		})

		variableComponent, _ := mcs.AddVariableMealComponent(VariableMealComponent{
			Name:                  "Variable Test Component",
			ProteinPerHundredGram: 15,
			CarbsPerHundredGram:   25,
			FatPerHundredGram:     8,
		})

		payload := UpdateMealPayload{
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
		_, err := ms.UpdateMeal(int(added.ID), payload)
		assert.Nil(t, err)

		meal, _ := ms.GetMeal(int(added.ID))
		assert.Equal(t, "Updated Test Meal", meal.Name)
		assert.Equal(t, 1, len(meal.FixedComponentData))
		assert.Equal(t, 1, len(meal.VariableComponentData))
		assert.Equal(t, fixedComponent.GetTotalCalories(2)+variableComponent.GetTotalCalories(150), meal.GetTotalCalories()) // 330 (from fixed) + 231 (from variable)
	})

	t.Run("should return error when updating with non-existent component", func(t *testing.T) {
		added, _ := ms.AddMeal("test")

		payload := UpdateMealPayload{
			Name: "Updated Test Meal",
			Components: []MealComponent{
				{
					Type:  FixedMealComponentType,
					Count: 2,
					Id:    9999, // Non-existent component ID
				},
			},
		}
		_, err := ms.UpdateMeal(int(added.ID), payload)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "couldn't find fixed component")
	})

	t.Run("should return error when updating with invalid component type", func(t *testing.T) {
		added, _ := ms.AddMeal("test")

		payload := UpdateMealPayload{
			Name: "Updated Test Meal",
			Components: []MealComponent{
				{
					Type:  "", // Invalid component type
					Count: 2,
					Id:    1,
				},
			},
		}
		_, err := ms.UpdateMeal(int(added.ID), payload)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unknown component type")
	})

	t.Run("should update meal macros correctly", func(t *testing.T) {
		added, _ := ms.AddMeal("test")

		fixedComponent, _ := mcs.AddFixedMealComponent(FixedMealComponent{
			Name:    "Fixed Test Component",
			Protein: 10,
			Carbs:   20,
			Fat:     5,
		})

		payload := UpdateMealPayload{
			Name: "Updated Test Meal",
			Components: []MealComponent{
				{
					Type:  FixedMealComponentType,
					Count: 3,
					Id:    int(fixedComponent.ID),
				},
			},
		}
		_, err := ms.UpdateMeal(int(added.ID), payload)
		assert.Nil(t, err)

		meal, _ := ms.GetMeal(int(added.ID))
		assert.Equal(t, 30, meal.GetMealMacro(Protein))
		assert.Equal(t, 60, meal.GetMealMacro(Carbs))
		assert.Equal(t, 15, meal.GetMealMacro(Fat))
	})
}

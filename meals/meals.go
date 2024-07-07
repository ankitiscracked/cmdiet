package meals

import (
	"fmt"

	"gorm.io/gorm"
)

var MS *MealServiceImpl

type MealService interface {
	AddMeal(name string, calories int) (int64, error)
	GetMeal(mealId int) Meal
	GetAllMeals() []Meal
	UpdateMeal(mealId int, protien, carbs, fat int)
}

type MealServiceImpl struct {
	DB *gorm.DB
}

func (m *MealServiceImpl) AddMeal(name string, calories int) (Meal, error) {
	meal := Meal{Name: name, Calories: calories}
	if err := m.DB.Create(&meal).Error; err != nil {
		return meal, fmt.Errorf("couldn't add meal %v", err)
	}
	return meal, nil
}

func (m *MealServiceImpl) UpdateMeal(mealId int, protien, carbs, fat int) {
	// now := time.Now().Unix()
	// updateSql := `update meals set protein = ?, carbs = ?, fat = ? where id = ?`
	// statement, err := database.DB.Prepare(updateSql)
	//
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("Meal macros has been updated")
}

func (m *MealServiceImpl) GetAllMeals() ([]Meal, error) {
	var meals []Meal
	if err := m.DB.Find(&meals).Error; err != nil {
		return meals, fmt.Errorf("couldn't get meals: %v", err)
	}
	return meals, nil
}

func (m *MealServiceImpl) GetMeal(mealId int) (Meal, error) {
	var meal Meal
	if err := m.DB.Find(&meal, mealId).Error; err != nil {
		return meal, fmt.Errorf("couldn't get meal: %v", err)
	}
	return meal, nil
}

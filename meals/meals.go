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

func (m *MealServiceImpl) UpdateMeal(updatedMeal Meal) error {
	var meal Meal
	if err := m.DB.Where("id", updatedMeal.Id).First(&meal).Error; err != nil {
		return fmt.Errorf("couldn't find meal with id %d, %v", updatedMeal.Id, err)
	}
	if err := m.DB.Save(updatedMeal).Error; err != nil {
		return fmt.Errorf("couldn't update meal %v", err)
	}

	return nil
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
	if err := m.DB.First(&meal, mealId).Error; err != nil {
		return meal, fmt.Errorf("couldn't get meal: %v", err)
	}
	return meal, nil
}

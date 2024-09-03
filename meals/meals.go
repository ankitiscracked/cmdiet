package meals

import (
	"cmdiet/constants"
	"fmt"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var MS *MealServiceImpl

type MealService interface {
	AddMeal(name string, calories int) (int64, error)
	AddMealWithMacros(meal Meal)
	GetMeal(mealId int) Meal
	GetAllMeals() []Meal
	UpdateMeal(mealId int, meal Meal)
}
type MealServiceImpl struct {
	DB *gorm.DB
}

type MacroType int

const (
	Protein MacroType = iota
	Carbs
	Fat
)

func (m MacroType) String() string {
	switch m {
	case Protein:
		return "Protein"
	case Carbs:
		return "Carbs"
	case Fat:
		return "Fat"
	default:
		return "Unknown"
	}
}

func (meal *Meal) GetTotalCalories() int {
	var totalCalories int
	if meal.FixedComponentData != nil {
		for _, d := range meal.FixedComponentData {
			totalCalories += d.Amount * (d.Component.Protein*4 + d.Component.Carbs*4 + d.Component.Fat*9)
		}
	}

	if meal.VariableComponentData != nil {
		for _, d := range meal.VariableComponentData {
			totalCalories += (d.Component.ProteinPerHundredGram*d.AmountInGrams/100)*4 +
				(d.Component.CarbsPerHundredGram*d.AmountInGrams/100)*4 +
				(d.Component.FatPerHundredGram*d.AmountInGrams/100)*9
		}
	}
	return totalCalories
}

func (meal *Meal) GetMealMacro(macroType MacroType) int {
	var totalMacro int
	if meal.FixedComponentData != nil {
		for _, d := range meal.FixedComponentData {
			switch macroType {
			case Protein:
				totalMacro += d.Amount * d.Component.Protein
			case Carbs:
				totalMacro += d.Amount * d.Component.Carbs
			case Fat:
				totalMacro += d.Amount * d.Component.Fat
			}
		}
	}

	if meal.VariableComponentData != nil {
		for _, d := range meal.VariableComponentData {
			switch macroType {
			case Protein:
				totalMacro += (d.Component.ProteinPerHundredGram * d.AmountInGrams) / 100
			case Carbs:
				totalMacro += (d.Component.CarbsPerHundredGram * d.AmountInGrams) / 100
			case Fat:
				totalMacro += (d.Component.FatPerHundredGram * d.AmountInGrams) / 100
			}
		}
	}
	return totalMacro
}

func (m *MealServiceImpl) AddMeal(name string) (Meal, error) {
	meal := Meal{Name: name}
	if err := m.DB.Create(&meal).Error; err != nil {
		return meal, fmt.Errorf("couldn't add meal %v", err)
	}
	return meal, nil
}

func (m *MealServiceImpl) AddMealWithComponents(payload MealPayload) (Meal, error) {
	if err := constants.Validator.Struct(payload); err != nil {
		return Meal{}, fmt.Errorf("invalid meal payload: %v", err)
	}

	var meal Meal
	if err := updateMealFromPayload(&meal, UpdateMealPayload(payload), m.DB); err != nil {
		return meal, err
	}

	if err := m.DB.Create(&meal).Error; err != nil {
		return Meal{}, fmt.Errorf("couldn't add meal: %v", err)
	}

	return meal, nil
}

func (m *MealServiceImpl) UpdateMeal(mealId int, payload UpdateMealPayload) (Meal, error) {
	if err := constants.Validator.Struct(payload); err != nil {
		return Meal{}, fmt.Errorf("invalid meal payload: %v", err)
	}

	meal, err := m.GetMeal(mealId)
	if err != nil {
		return meal, fmt.Errorf("couldn't find meal with id %d, %v", mealId, err)
	}
	if err := m.DB.Unscoped().Model(&meal).Association("FixedComponentData").Unscoped().Clear(); err != nil {
		return Meal{}, fmt.Errorf("couldn't delete fixed components: %v", err)
	}
	if err := m.DB.Unscoped().Model(&meal).Association("VariableComponentData").Unscoped().Clear(); err != nil {
		return Meal{}, fmt.Errorf("couldn't delete variable components: %v", err)
	}

	if err := updateMealFromPayload(&meal, payload, m.DB); err != nil {
		return meal, err
	}

	err = m.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&meal).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Model(&meal).Association("FixedComponentData").Unscoped().Replace(meal.FixedComponentData); err != nil {
			return err
		}
		if err := tx.Unscoped().Model(&meal).Association("VariableComponentData").Unscoped().Replace(meal.VariableComponentData); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return Meal{}, fmt.Errorf("couldn't update meal: %v", err)
	}

	return meal, nil
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
	if err := m.DB.Preload("FixedComponentData.Component").
		Preload("VariableComponentData.Component").
		Preload(clause.Associations).
		First(&meal, mealId).Error; err != nil {
		return meal, fmt.Errorf("couldn't get meal: %v", err)
	}
	return meal, nil
}

func updateMealFromPayload(meal *Meal, payload UpdateMealPayload, DB *gorm.DB) error {
	meal.Name = payload.Name

	for _, comp := range payload.Components {
		switch comp.Type {
		case FixedMealComponentType:
			var fixedComp FixedMealComponent
			if err := DB.First(&fixedComp, comp.Id).Error; err != nil {
				return fmt.Errorf("couldn't find fixed component with id %d: %v", comp.Id, err)
			}
			meal.FixedComponentData = append(meal.FixedComponentData, FixedMealComponentData{
				Component: fixedComp,
				Amount:    comp.Count,
			})
		case VariableMealComponentType:
			var varComp VariableMealComponent
			if err := DB.First(&varComp, comp.Id).Error; err != nil {
				return fmt.Errorf("couldn't find variable component with id %d: %v", comp.Id, err)
			}
			meal.VariableComponentData = append(meal.VariableComponentData, VariableMealComponentData{
				Component:     varComp,
				AmountInGrams: comp.AmountInGrams,
			})
		default:
			return fmt.Errorf("unknown component type")
		}
	}
	return nil
}

func AtoiIgnoreError(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

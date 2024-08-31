package meals

import (
	"cmdiet/constants"
	"fmt"
	"strconv"

	"gorm.io/gorm"
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
	_, err := updateMealFromPayload(&meal, payload, m.DB)
	if err != nil {
		return meal, err
	}

	if err := m.DB.Create(&meal).Error; err != nil {
		return Meal{}, fmt.Errorf("couldn't add meal: %v", err)
	}

	return meal, nil
}

func (m *MealServiceImpl) UpdateMeal(mealId int, payload MealPayload) (Meal, error) {
	if err := constants.Validator.Struct(payload); err != nil {
		return Meal{}, fmt.Errorf("invalid meal payload: %v", err)
	}

	var meal Meal
	if err := m.DB.Where("id", mealId).First(&meal).Error; err != nil {
		return meal, fmt.Errorf("couldn't find meal with id %d, %v", mealId, err)
	}

	// Delete existing components from the meal
	if err := m.DB.Model(&meal).Association("FixedComponentData").Clear(); err != nil {
		return Meal{}, fmt.Errorf("couldn't delete fixed components: %v", err)
	}
	if err := m.DB.Model(&meal).Association("VariableComponentData").Clear(); err != nil {
		return Meal{}, fmt.Errorf("couldn't delete variable components: %v", err)
	}

	_, err := updateMealFromPayload(&meal, payload, m.DB)
	if err != nil {
		return meal, err
	}

	if err := m.DB.Save(&meal).Error; err != nil {
		return meal, fmt.Errorf("couldn't update meal %v", err)
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
	if err := m.DB.First(&meal, mealId).Error; err != nil {
		return meal, fmt.Errorf("couldn't get meal: %v", err)
	}
	return meal, nil
}

func (m *MealServiceImpl) AddFixedMealComponent(mealComponent FixedMealComponent) (FixedMealComponent, error) {
	if err := constants.Validator.Struct(mealComponent); err != nil {
		return mealComponent, fmt.Errorf("invalid meal component: %v", err)
	}
	if err := m.DB.Create(&mealComponent).Error; err != nil {
		return mealComponent, fmt.Errorf("couldn't add meal %v", err)
	}
	return mealComponent, nil
}

func (m *MealServiceImpl) AddVariableMealComponent(mealComponent VariableMealComponent) (VariableMealComponent, error) {
	if err := constants.Validator.Struct(mealComponent); err != nil {
		return mealComponent, fmt.Errorf("invalid meal component: %v", err)
	}
	if err := m.DB.Create(&mealComponent).Error; err != nil {
		return mealComponent, fmt.Errorf("couldn't add meal %v", err)
	}
	return mealComponent, nil
}

func (m *MealServiceImpl) GetFixedMealComponent(componentId int) (FixedMealComponent, error) {
	var component FixedMealComponent
	if err := m.DB.First(&component, componentId).Error; err != nil {
		return component, fmt.Errorf("couldn't get fixed meal component: %v", err)
	}
	return component, nil
}

func (m *MealServiceImpl) GetVariableMealComponent(componentId int) (VariableMealComponent, error) {
	var component VariableMealComponent
	if err := m.DB.First(&component, componentId).Error; err != nil {
		return component, fmt.Errorf("couldn't get variable meal component: %v", err)
	}
	return component, nil
}

func (m *MealServiceImpl) GetAllMealComponents() ([]MealComponentInput, error) {
	var components []MealComponentInput

	// Fetch fixed meal components
	var fixedComponents []FixedMealComponent
	if err := m.DB.Find(&fixedComponents).Error; err != nil {
		return nil, fmt.Errorf("couldn't get fixed meal components: %v", err)
	}

	// Fetch variable meal components
	var variableComponents []VariableMealComponent
	if err := m.DB.Find(&variableComponents).Error; err != nil {
		return nil, fmt.Errorf("couldn't get variable meal components: %v", err)
	}

	// Convert fixed components to MealComponentInput type
	for _, fc := range fixedComponents {
		components = append(components, MealComponentInput{
			Name: fc.Name,
			Type: FixedMealComponentType,
			Id:   int(fc.ID),
		})
	}

	// Convert variable components to MealComponentInput type
	for _, vc := range variableComponents {
		components = append(components, MealComponentInput{
			Name: vc.Name,
			Type: VariableMealComponentType,
			Id:   int(vc.ID),
		})
	}

	return components, nil
}

func updateMealFromPayload(meal *Meal, payload MealPayload, DB *gorm.DB) (*Meal, error) {
	meal.Name = payload.Name

	for _, comp := range payload.Components {
		switch comp.Type {
		case FixedMealComponentType:
			var fixedComp FixedMealComponent
			if err := DB.First(&fixedComp, comp.Id).Error; err != nil {
				return meal, fmt.Errorf("couldn't find fixed component with id %d: %v", comp.Id, err)
			}
			meal.FixedComponentData = append(meal.FixedComponentData, FixedMealComponentData{
				Component: fixedComp,
				Amount:    comp.Count,
			})
		case VariableMealComponentType:
			var varComp VariableMealComponent
			if err := DB.First(&varComp, comp.Id).Error; err != nil {
				return meal, fmt.Errorf("couldn't find variable component with id %d: %v", comp.Id, err)
			}
			meal.VariableComponentData = append(meal.VariableComponentData, VariableMealComponentData{
				Component:     varComp,
				AmountInGrams: comp.AmountInGrams,
			})
		default:
			return meal, fmt.Errorf("unknown component type")
		}
	}
	return meal, nil
}

func AtoiIgnoreError(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

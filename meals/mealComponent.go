package meals

import (
	"cmdiet/constants"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var MCS *MealComponentService

type MealComponentService struct {
	DB *gorm.DB
}

func (m *MealComponentService) AddFixedMealComponent(mealComponent FixedMealComponent) (FixedMealComponent, error) {
	if err := constants.Validator.Struct(mealComponent); err != nil {
		return mealComponent, fmt.Errorf("invalid meal component: %v", err)
	}
	if err := m.DB.Create(&mealComponent).Error; err != nil {
		return mealComponent, fmt.Errorf("couldn't add meal %v", err)
	}
	return mealComponent, nil
}

func (m *MealComponentService) AddVariableMealComponent(mealComponent VariableMealComponent) (VariableMealComponent, error) {
	if err := constants.Validator.Struct(mealComponent); err != nil {
		return mealComponent, fmt.Errorf("invalid meal component: %v", err)
	}
	if err := m.DB.Create(&mealComponent).Error; err != nil {
		return mealComponent, fmt.Errorf("couldn't add meal %v", err)
	}
	return mealComponent, nil
}

func (m *MealComponentService) UpdateMealComponent(componentId int, payload UpdateMealComponentPayload) error {
	if err := constants.Validator.Struct(payload); err != nil {
		return fmt.Errorf("invalid update payload: %v", err)
	}

	switch payload.Type {
	case FixedMealComponentType:
		var component FixedMealComponent
		if err := m.DB.First(&component, componentId).Error; err != nil {
			return fmt.Errorf("couldn't find fixed meal component: %v", err)
		}
		component.Name = payload.Name
		component.Protein = payload.Protein
		component.Carbs = payload.Carbs
		component.Fat = payload.Fat
		if err := m.DB.Save(&component).Error; err != nil {
			return fmt.Errorf("couldn't update fixed meal component: %v", err)
		}
	case VariableMealComponentType:
		var component VariableMealComponent
		if err := m.DB.First(&component, componentId).Error; err != nil {
			return fmt.Errorf("couldn't find variable meal component: %v", err)
		}
		component.Name = payload.Name
		component.ProteinPerHundredGram = payload.Protein
		component.CarbsPerHundredGram = payload.Carbs
		component.FatPerHundredGram = payload.Fat
		if err := m.DB.Save(&component).Error; err != nil {
			return fmt.Errorf("couldn't update variable meal component: %v", err)
		}
	default:
		return fmt.Errorf("invalid meal component type")
	}

	return nil
}

func (m *MealComponentService) GetFixedMealComponent(componentId int) (FixedMealComponent, error) {
	var component FixedMealComponent
	if err := m.DB.First(&component, componentId).Error; err != nil {
		return component, fmt.Errorf("couldn't get fixed meal component: %v", err)
	}
	return component, nil
}

func (m *MealComponentService) GetVariableMealComponent(componentId int) (VariableMealComponent, error) {
	var component VariableMealComponent
	if err := m.DB.First(&component, componentId).Error; err != nil {
		return component, fmt.Errorf("couldn't get variable meal component: %v", err)
	}
	return component, nil
}

func (m *MealComponentService) GetAllMealComponents() ([]MealComponentInput, error) {
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
			Name:          fc.Name,
			Type:          FixedMealComponentType,
			Id:            int(fc.ID),
			TotalCalories: fc.GetTotalCalories(1),
			TotalProtein:  fc.Protein,
			TotalCarbs:    fc.Carbs,
			TotalFat:      fc.Fat,
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

func (mct MealComponentType) IsValid() bool {
	return mct == FixedMealComponentType || mct == VariableMealComponentType
}

func ValidateMealComponentType(fl validator.FieldLevel) bool {
	if mct, ok := fl.Field().Interface().(MealComponentType); ok {
		return mct.IsValid()
	}
	return false
}

func (fc *FixedMealComponent) GetTotalCalories(count int) int {
	return (fc.Protein*4 + fc.Carbs*4 + fc.Fat*9) * count
}

func (vc *VariableMealComponent) GetTotalCalories(amountInGrams int) int {
	proteinCalories := (vc.ProteinPerHundredGram * amountInGrams / 100) * 4
	carbsCalories := (vc.CarbsPerHundredGram * amountInGrams / 100) * 4
	fatCalories := (vc.FatPerHundredGram * amountInGrams / 100) * 9

	totalCalories := proteinCalories + carbsCalories + fatCalories

	return int(totalCalories)
}

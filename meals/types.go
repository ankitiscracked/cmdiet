package meals

import "gorm.io/gorm"

type MealComponentType string

const (
	FixedMealComponentType    MealComponentType = "fixed"
	VariableMealComponentType MealComponentType = "variable"
)

type Meal struct {
	gorm.Model
	Name                  string
	FixedComponentData    []FixedMealComponentData
	VariableComponentData []VariableMealComponentData
}

type FixedMealComponent struct {
	gorm.Model
	Name    string `validate:"required"`
	Protein int    `validate:"required"`
	Carbs   int    `validate:"required"`
	Fat     int    `validate:"required"`
}

type VariableMealComponent struct {
	gorm.Model
	Name                  string `validate:"required"`
	ProteinPerHundredGram int    `validate:"required"`
	CarbsPerHundredGram   int    `validate:"required"`
	FatPerHundredGram     int    `validate:"required"`
}

type FixedMealComponentData struct {
	gorm.Model
	ComponentID int
	Component   FixedMealComponent `gorm:"foreignKey:ComponentID"`
	Amount      int
	MealID      uint
}

type VariableMealComponentData struct {
	gorm.Model
	ComponentID   int
	Component     VariableMealComponent `gorm:"foreignKey:ComponentID"`
	AmountInGrams int
	MealID        uint
}

type MealComponent struct {
	Name          string
	Type          MealComponentType `validate:"required,oneof=fixed variable"`
	Id            int               `validate:"required"`
	Count         int               `validate:"required_without=AmountInGrams"`
	AmountInGrams int               `validate:"required_without=Count"`
}

type MealComponentInput struct {
	Name          string
	Type          MealComponentType
	Id            int
	Amount        string
	TotalCalories int
	TotalProtein  int
	TotalCarbs    int
	TotalFat      int
}
type MealPayload struct {
	Name       string          `validate:"required"`
	Components []MealComponent `validate:"required,dive"`
}

type UpdateMealPayload struct {
	Name       string
	Components []MealComponent
}

type UpdateMealComponentPayload struct {
	ID      int
	Type    MealComponentType
	Name    string
	Protein int
	Carbs   int
	Fat     int
}

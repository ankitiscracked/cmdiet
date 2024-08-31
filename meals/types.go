package meals

import "gorm.io/gorm"

type MealComponentType int

const (
	FixedMealComponentType MealComponentType = iota
	VariableMealComponentType
)

type Meal struct {
	gorm.Model
	Id                    int
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
	Id        int
	Component FixedMealComponent
	Amount    int
}

type VariableMealComponentData struct {
	gorm.Model
	Id            int
	Component     VariableMealComponent
	AmountInGrams int
}

type MealComponent struct {
	Name          string
	Type          MealComponentType `validate:"required"`
	Id            int               `validate:"required"`
	Count         int               `validate:"required_without=amountInGrams"`
	AmountInGrams int               `validate:"required_without=count"`
}

type MealComponentInput struct {
	Name          string
	Type          MealComponentType
	Id            int
	Count         string
	AmountInGrams string
}
type MealPayload struct {
	Name       string          `validate:"required"`
	Components []MealComponent `validate:"required,dive"`
}

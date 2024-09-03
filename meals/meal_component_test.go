package meals

import (
	"cmdiet/constants"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestAddingMealComponent(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	mcs := MealComponentService{DB: db}
	constants.Validator = validator.New(validator.WithRequiredStructEnabled())
	t.Run("should add a fixed meal component", func(t *testing.T) {
		component := FixedMealComponent{
			Name:    "Test Fixed Component",
			Protein: 10,
			Carbs:   20,
			Fat:     5,
		}
		added, err := mcs.AddFixedMealComponent(component)
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
		added, err := mcs.AddVariableMealComponent(component)
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
		_, err := mcs.AddFixedMealComponent(component)
		assert.NotNil(t, err)
	})

	t.Run("should not add a variable meal component with invalid data", func(t *testing.T) {
		component := VariableMealComponent{
			Name:                  "",
			ProteinPerHundredGram: -15,
			CarbsPerHundredGram:   25,
			FatPerHundredGram:     8,
		}
		_, err := mcs.AddVariableMealComponent(component)
		assert.NotNil(t, err)
	})
}

func TestFetchingMealComponent(t *testing.T) {
	db, err := setupTestDB(t)
	if err != nil {
		t.Fatal(err)
	}

	mcs := MealComponentService{DB: db}

	t.Run("should retrieve a fixed meal component", func(t *testing.T) {
		component := FixedMealComponent{
			Name:    "Retrievable Fixed Component",
			Protein: 12,
			Carbs:   22,
			Fat:     7,
		}
		added, _ := mcs.AddFixedMealComponent(component)
		retrieved, err := mcs.GetFixedMealComponent(int(added.ID))
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
		added, _ := mcs.AddVariableMealComponent(component)
		retrieved, err := mcs.GetVariableMealComponent(int(added.ID))
		assert.Nil(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, added.ID, retrieved.ID)
		assert.Equal(t, component.Name, retrieved.Name)
	})

	t.Run("should return error when retrieving non-existent fixed component", func(t *testing.T) {
		_, err := mcs.GetFixedMealComponent(9999)
		assert.NotNil(t, err)
	})

	t.Run("should return error when retrieving non-existent variable component", func(t *testing.T) {
		_, err := mcs.GetVariableMealComponent(9999)
		assert.NotNil(t, err)
	})
}

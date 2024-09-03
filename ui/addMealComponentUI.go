package ui

import (
	"cmdiet/meals"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type EnteredMealComponent struct {
	name    string
	protein string
	carbs   string
	fat     string
}
type AddMealComponentModel struct {
	form                 *huh.Form
	err                  error
	done                 bool
	componentType        meals.MealComponentType
	enteredMealComponent EnteredMealComponent
}

func NewAddMealComponentModel(componentType meals.MealComponentType) *AddMealComponentModel {
	m := &AddMealComponentModel{componentType: componentType}
	m.initForm()
	return m
}

func (m *AddMealComponentModel) initForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the name of the meal component").
				Placeholder("e.g. Chicken Breast").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name cannot be empty")
					}
					return nil
				}).
				Value(&m.enteredMealComponent.name),
			huh.NewInput().
				Title(macroInputTitle(meals.Protein, m.componentType)).
				Placeholder("e.g. 25").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("protein content cannot be empty")
					}
					return nil
				}).
				Value(&m.enteredMealComponent.protein),
			huh.NewInput().
				Title(macroInputTitle(meals.Carbs, m.componentType)).
				Placeholder("e.g. 30").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("carbs content cannot be empty")
					}
					return nil
				}).
				Value(&m.enteredMealComponent.carbs),
			huh.NewInput().
				Title(macroInputTitle(meals.Fat, m.componentType)).
				Placeholder("e.g. 10").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("fat content cannot be empty")
					}
					return nil
				}).
				Value(&m.enteredMealComponent.fat),
		),
	)
}

func macroInputTitle(macroType meals.MacroType, componentType meals.MealComponentType) string {
	if componentType == meals.FixedMealComponentType {
		return fmt.Sprintf("%s Content", macroType.String())
	} else {
		return fmt.Sprintf("%s Content per 100 grams", macroType.String())
	}
}

func (m *AddMealComponentModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *AddMealComponentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		if m.form.State == huh.StateCompleted {
			m.done = true

			err := createMealComponent(m)
			if err != nil {
				m.err = err
				return m, nil
			}

			return m, tea.Quit
		}
	}
	return m, cmd
}

func createMealComponent(m *AddMealComponentModel) error {
	var err error
	if m.componentType == meals.FixedMealComponentType {
		payload := meals.FixedMealComponent{
			Name:    m.enteredMealComponent.name,
			Protein: atoiIgnoreError(m.enteredMealComponent.protein),
			Carbs:   atoiIgnoreError(m.enteredMealComponent.carbs),
			Fat:     atoiIgnoreError(m.enteredMealComponent.fat),
		}
		_, err = meals.MCS.AddFixedMealComponent(payload)
	} else {
		payload := meals.VariableMealComponent{
			Name:                  m.enteredMealComponent.name,
			ProteinPerHundredGram: atoiIgnoreError(m.enteredMealComponent.protein),
			CarbsPerHundredGram:   atoiIgnoreError(m.enteredMealComponent.carbs),
			FatPerHundredGram:     atoiIgnoreError(m.enteredMealComponent.fat),
		}
		_, err = meals.MCS.AddVariableMealComponent(payload)
	}
	return err
}

func (m *AddMealComponentModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	if m.done {
		return fmt.Sprintf("Added %s meal component: %s \n", m.componentType, m.enteredMealComponent.name)
	}
	return m.form.View()
}

package ui

import (
	"cmdiet/meals"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type errMsg error
type AddMealModel struct {
	form               *huh.Form
	err                error
	done               bool
	componentsOptions  []meals.MealComponentInput
	selectedComponents []meals.MealComponentInput
	mealName           string
}

func NewAddMealModel() *AddMealModel {
	m := &AddMealModel{}
	m.componentsOptions, _ = meals.MCS.GetAllMealComponents()
	m.selectedComponents = append(m.selectedComponents, meals.MealComponentInput{})
	m.createForm(m.selectedComponents)
	return m
}

func (m *AddMealModel) createForm(selectedComponents []meals.MealComponentInput) {
	nameInput := huh.NewInput().
		Title("Meal Name").
		Description("Enter the name of the meal").
		Placeholder("e.g. Chicken Salad").
		Validate(func(s string) error {
			if s == "" {
				return fmt.Errorf("meal name cannot be empty")
			}
			return nil
		}).
		Value(&m.mealName)

	m.form = huh.NewForm(m.createFormGroup(nameInput))
}

func (m *AddMealModel) createFormGroup(nameField *huh.Input) *huh.Group {
	var inputs []huh.Field

	for i := range m.selectedComponents {
		comp := &m.selectedComponents[i]
		componentOptions := m.getComponentOptions()

		inputs = append(
			inputs,
			huh.NewSelect[meals.MealComponentInput]().
				Title("Component").
				Options(componentOptions...).
				Value(comp),

			huh.NewInput().
				Title("Amount").
				Description(fmt.Sprintf("Enter the amount for %s", comp.Name)).
				Placeholder("e.g. 2 or 200 (gms)").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("count cannot be empty")
					}
					return nil
				}).
				Value(&comp.Amount))
	}

	return huh.NewGroup(append([]huh.Field{nameField}, inputs...)...)
}

func (m *AddMealModel) getComponentOptions() []huh.Option[meals.MealComponentInput] {
	options := make([]huh.Option[meals.MealComponentInput], 0, len(m.componentsOptions))
	for i := range m.componentsOptions {
		comp := &m.componentsOptions[i]
		options = append(options, huh.NewOption(comp.Name, *comp))
	}
	return options
}

func (m *AddMealModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *AddMealModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+n":
			m.selectedComponents = append(m.selectedComponents, meals.MealComponentInput{})
			m.createForm(m.selectedComponents)
			m.form.Init()
		}
	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		if m.form.State == huh.StateCompleted {
			m.done = true

			err := m.saveMeal()
			if err != nil {
				m.err = err
				return m, nil
			}

			return m, tea.Quit

		}
	}
	return m, cmd
}

func (m *AddMealModel) saveMeal() error {
	payload := meals.MealPayload{
		Name:       m.mealName,
		Components: make([]meals.MealComponent, 0),
	}

	log.Println("component opitons: ", m.componentsOptions)
	log.Println("selected components: ", m.selectedComponents)
	for _, comp := range m.selectedComponents {
		mealComponent := meals.MealComponent{
			Id:   comp.Id,
			Type: comp.Type,
		}

		switch comp.Type {
		case meals.FixedMealComponentType:
			mealComponent.Count = meals.AtoiIgnoreError(comp.Amount)
		case meals.VariableMealComponentType:
			mealComponent.AmountInGrams = meals.AtoiIgnoreError(comp.Amount)
		}

		payload.Components = append(payload.Components, mealComponent)
	}

	_, err := meals.MS.AddMealWithComponents(payload)
	if err != nil {
		return err
	}
	return nil
}

func (m *AddMealModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	if m.done {
		return fmt.Sprintf("Added meal: %s\n", m.mealName)
	}
	return m.form.View()
}

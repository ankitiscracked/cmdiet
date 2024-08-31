package ui

import (
	"cmdiet/meals"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type errMsg error
type AddMealModel struct {
	form               *huh.Form
	err                error
	done               bool
	result             meals.Meal
	components         []meals.MealComponentInput
	selectedComponents []meals.MealComponentInput
	componentAmount    float64
}

func NewAddMealModel() *AddMealModel {
	m := &AddMealModel{}
	m.initForm()
	return m
}

func (m *AddMealModel) createComponentFormGroups() []*huh.Group {
	var groups []*huh.Group

	for _, comp := range m.selectedComponents {
		selectedComp := &comp
		var amountField huh.Field

		switch comp.Type {
		case meals.FixedMealComponentType:
			amountField = huh.NewInput().
				Title("Count").
				Description(fmt.Sprintf("Enter the count for %s", comp.Name)).
				Placeholder("e.g. 2").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("count cannot be empty")
					}
					return nil
				}).
				Value(&comp.Count)
		case meals.VariableMealComponentType:
			amountField = huh.NewInput().
				Title("Amount in grams").
				Description(fmt.Sprintf("Enter the amount in grams for %s", comp.Name)).
				Placeholder("e.g. 100").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("amount cannot be empty")
					}
					return nil
				}).
				Value(&comp.AmountInGrams)
		}

		componentOptions := m.getComponentOptions()
		group := huh.NewGroup(
			huh.NewSelect[*meals.MealComponentInput]().
				Title("Component").
				Options(componentOptions...).
				Value(&selectedComp),
			amountField,
		).Title(comp.Name)

		groups = append(groups, group)
	}

	return groups
}

func (m *AddMealModel) getComponentOptions() []huh.Option[*meals.MealComponentInput] {
	options := []huh.Option[*meals.MealComponentInput]{}
	for i := range m.components {
		comp := &m.components[i]
		options = append(options, huh.NewOption(comp.Name, comp))
	}
	return options
}

func (m *AddMealModel) initForm() {
	m.form = huh.NewForm(m.createComponentFormGroups()...)
}

func (m *AddMealModel) Init() tea.Cmd {
	return tea.Batch(
		m.form.Init(),
		m.fetchComponents,
	)
}

func (m *AddMealModel) fetchComponents() tea.Msg {
	components, err := meals.MS.GetAllMealComponents()
	if err != nil {
		return errMsg(err)
	}
	return componentsMsg{components: components}
}

func (m *AddMealModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case componentsMsg:
		m.components = msg.components
		m.initForm()
	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		if m.form.State == huh.StateCompleted {
			m.done = true
			return m, tea.Sequence(m.saveMeal, tea.Quit)
		}
	}
	return m, cmd
}

func (m *AddMealModel) saveMeal() tea.Msg {
	payload := meals.MealPayload{
		Name:       m.result.Name,
		Components: make([]meals.MealComponent, 0),
	}

	for _, comp := range m.selectedComponents {
		payload.Components = append(payload.Components, meals.MealComponent{
			Id:            comp.Id,
			Type:          comp.Type,
			Count:         meals.AtoiIgnoreError(comp.Count),
			AmountInGrams: meals.AtoiIgnoreError(comp.AmountInGrams),
		})
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
		return fmt.Sprintf("Added meal: %s\n", m.result.Name)
	}
	return m.form.View()
}

type componentsMsg struct {
	components []meals.MealComponentInput
}

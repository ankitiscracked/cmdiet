package ui

import (
	"cmdiet/meals"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type AddMealComponentModel struct {
	form   *huh.Form
	err    error
	done   bool
	result meals.MealComponent
}

func NewAddMealComponentModel(componentType meals.MealComponentType) *AddMealComponentModel {
	m := &AddMealComponentModel{}
	m.initForm()
	return m
}

func (m *AddMealComponentModel) initForm() {
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Description("Enter the name of the meal component").
				Placeholder("e.g. Chicken Breast").
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("name cannot be empty")
					}
					return nil
				}).
				Value(&m.result.Name),
			huh.NewSelect[string]().
				Title("Component Type").
				Description("Choose the type of meal component").
				Options(
					huh.NewOption("Fixed", "fixed"),
					huh.NewOption("Variable", "variable"),
				),
		),
	)
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
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m *AddMealComponentModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	if m.done {
		return fmt.Sprintf("Added meal component: %s (%s)\n", m.result.Name, m.result.Type)
	}
	return m.form.View()
}

package ui

import (
	"cmdiet/diet"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle.Copy()

	blurredButton = fmt.Sprintf("[%s]", blurredStyle.Render("Log"))
	focusedButton = focusedStyle.Copy().Render("Log")
	formStyle     = lipgloss.NewStyle().Padding(1, 2)
)

var (
	meal     string
	calories string
	source   string
)

type (
	errMsg error
	model  struct {
		focusIndex int
		inputs     []textinput.Model
		err        error
		cursorMode cursor.Mode
		mealType   string
		form       *huh.Form
	}
)

func NewDietModel(mealType string) model {
	m := model{mealType: mealType}
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Enter the meal name").Value(&meal),
			huh.NewInput().Title("Enter the no. of calories").Value(&calories),
			huh.NewInput().Title("Enter the source").Prompt("Home/Out").Value(&source),
		),
	)

	return m
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		i, err := strconv.Atoi(calories)
		if err != nil {
			m.err = fmt.Errorf("calories must be a number")
		}
		diet.DefaultDietService.LogDiet(m.mealType, meal, i, source)
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m model) View() string {
	return formStyle.Render(m.form.View())
}

// func LogMealModel(mealType string) model {
// 	newModel := model{
// 		inputs:   make([]textinput.Model, 3),
// 		mealType: mealType,
// 	}
//
// 	for i := range newModel.inputs {
// 		input := textinput.New()
// 		input.Cursor.Style = cursorStyle
// 		input.CharLimit = 32
//
// 		switch i {
// 		case 0:
// 			input.Placeholder = "What did you eat?"
// 			input.Focus()
// 			input.PromptStyle = focusedStyle
// 			input.TextStyle = focusedStyle
// 		case 1:
// 			input.Placeholder = "Enter the calories. (Press enter to skip):"
// 			input.CharLimit = 5
// 		case 2:
// 			input.Placeholder = "Home or Out?"
// 			input.CharLimit = 5
// 		}
//
// 		newModel.inputs[i] = input
// 	}
//
// 	return newModel
// }

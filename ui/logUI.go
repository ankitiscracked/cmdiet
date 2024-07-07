package ui

import (
	"cmdiet/diet"
	"cmdiet/meals"
	"fmt"
	"log"
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
	mealId   int
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
		quitting   bool
	}
)

func NewDietModel(mealType string) model {
	m := model{mealType: mealType}
	m.form = createForm(
		huh.NewInput().Title("Enter the meal name").Value(&meal),
	)
	return m
}

func createForm(mealInput huh.Field) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			mealInput,
			huh.NewInput().Title("Enter the no. of calories").Value(&calories),
			huh.NewSelect[string]().
				Title("Where did you eat?").
				Options(
					huh.NewOption("Cooked", "cooked"),
					huh.NewOption("Orderd-in", "orderedin"),
					huh.NewOption("Went Out", "wentout"),
				).
				Value(&source),
		),
	)
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
		case "ctrl+s":
			meals, err := meals.MS.GetAllMeals()
			if err != nil {
				log.Fatal(err)
			}
			var options []huh.Option[int]
			for _, meal := range meals {
				options = append(options, huh.NewOption(meal.Name, meal.Id))
			}
			mealInput := huh.NewSelect[int]().Title("Select one of your meals").Options(options...).Height(10).Value(&mealId)
			m.form = createForm(mealInput)
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
		var logError error
		if mealId != 0 {
			logError = diet.DS.LogDietWithExistingMeal(int64(mealId), m.mealType, source)
		} else {
			logError = diet.DS.LogDietWithNewMeal(m.mealType, meal, i, source)
		}
		if logError != nil {
			log.Fatal(logError)
		}

		m.quitting = true
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return lipgloss.NewStyle().Render("I've logged your diet!")
	}
	return formStyle.Render(m.form.View())
}

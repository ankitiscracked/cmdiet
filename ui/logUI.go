package ui

import (
	"cmdiet/diet"
	"cmdiet/meals"
	"fmt"
	"log"
	"time"

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
	model struct {
		focusIndex int
		err        error
		mealType   diet.MealType
		form       *huh.Form
		quitting   bool
		logForDate time.Time
	}
)

func NewDietModel(mealType diet.MealType, logForDate time.Time) model {
	m := model{mealType: mealType, logForDate: logForDate}
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
				options = append(options, huh.NewOption(meal.Name, int(meal.ID)))
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
		var logError error
		if mealId != 0 {
			logError = diet.DS.LogDietWithExistingMeal(int64(mealId), m.mealType, source, m.logForDate)
		} else {
			logError = diet.DS.LogDietWithNewMeal(m.mealType, meal, source, m.logForDate)
		}
		if logError != nil {
			log.Fatal(logError)
		}

		m.quitting = true
		return m, tea.Quit
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return TerminalInfo("Keep logging, keep growing!")
	}

	return "\n" + lipgloss.NewStyle().
		MarginLeft(2).
		Padding(0, 1).
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#6C50FF")).
		Render(fmt.Sprintf("Let's log your %s for %s", m.mealType, m.logForDate.Format("02 Jan"))) +
		"\n\n" +
		formStyle.Render(m.form.View())
}

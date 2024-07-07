package ui

import (
	"cmdiet/meals"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listStyle        = lipgloss.NewStyle().Margin(1, 2)
	mealEditingStyle = lipgloss.NewStyle()
)

type listItem struct {
	title, desc string
	mealId      int
}

func (i listItem) Title() string {
	return i.title
}

func (i listItem) Description() string {
	return i.desc
}

func (i listItem) FilterValue() string {
	return i.title
}

type mealsModel struct {
	list          list.Model
	editing       bool
	currentDetail detailsModel
}

type detailsModel struct {
	textInputs []textinput.Model
	focusIndex int
}

func (m mealsModel) Init() tea.Cmd {
	return nil
}

func (m mealsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.editing {
				return m.handleMealDetailsUpdate(msg)
			} else {
				return m.setMealDetailsModel()
			}
		}

		if m.editing {
			return m.handleMealDetailsUpdate(msg)
		}
	case tea.WindowSizeMsg:
		x, y := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-x, msg.Height-y)
	}

	var cmd tea.Cmd
	if m.editing {
		cmd = m.updateInputs(msg)
	} else {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m mealsModel) updateInputs(msg tea.Msg) tea.Cmd {
	textInputs := m.currentDetail.textInputs
	cmds := make([]tea.Cmd, len(textInputs))
	for i := range textInputs {
		textInputs[i], cmds[i] = textInputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m mealsModel) setMealDetailsModel() (mealsModel, tea.Cmd) {
	selectedItem := m.list.SelectedItem().(listItem)
	inputs := make([]textinput.Model, 2)

	titleInput := textinput.New()
	titleInput.SetValue(selectedItem.title)
	titleInput.PromptStyle = focusedStyle
	titleInput.TextStyle = focusedStyle
	inputs[0] = titleInput

	caloriesInput := textinput.New()
	caloriesInput.SetValue(selectedItem.desc)
	caloriesInput.TextStyle = blurredStyle
	inputs[1] = caloriesInput

	details := detailsModel{textInputs: inputs, focusIndex: 0}
	m.currentDetail = details
	m.editing = true

	return m, m.currentDetail.textInputs[0].Focus()
}

func (m mealsModel) handleMealDetailsUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyTab, tea.KeyShiftTab, tea.KeyEnter, tea.KeyUp, tea.KeyDown:
		keyType := msg.Type

		if keyType == tea.KeyEnter && m.currentDetail.focusIndex == len(m.currentDetail.textInputs) {
			m.editing = false
			return m, nil
		}

		if keyType == tea.KeyUp || keyType == tea.KeyShiftTab {
			m.currentDetail.focusIndex--
		} else {
			m.currentDetail.focusIndex++
		}

		if m.currentDetail.focusIndex > len(m.currentDetail.textInputs) {
			m.currentDetail.focusIndex = 0
		} else if m.currentDetail.focusIndex < 0 {
			m.currentDetail.focusIndex = len(m.currentDetail.textInputs)
		}
	}

	cmds := make([]tea.Cmd, len(m.currentDetail.textInputs))
	for i := 0; i < len(m.currentDetail.textInputs); i++ {
		if i == m.currentDetail.focusIndex {
			cmds[i] = m.currentDetail.textInputs[i].Focus()
			m.currentDetail.textInputs[i].PromptStyle = focusedStyle
			m.currentDetail.textInputs[i].TextStyle = focusedStyle
			continue
		}

		m.currentDetail.textInputs[i].Blur()
		m.currentDetail.textInputs[i].PromptStyle = blurredStyle
		m.currentDetail.textInputs[i].TextStyle = blurredStyle
	}
	return m, tea.Batch(cmds...)
}

func (m mealsModel) View() string {
	if m.editing {
		var b strings.Builder

		for i := range m.currentDetail.textInputs {
			b.WriteString(m.currentDetail.textInputs[i].View())
			if i < len(m.currentDetail.textInputs)-1 {
				b.WriteRune('\n')
			}
		}

		button := &blurredButton
		if m.currentDetail.focusIndex == len(m.currentDetail.textInputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)

		return lipgloss.JoinHorizontal(lipgloss.Center, m.list.View(), b.String())
	} else {
		return listStyle.Render(m.list.View())
	}
}

func ViewMealsModel() mealsModel {
	allMeals, err := meals.MS.GetAllMeals()
	if err != nil {
		log.Fatal(err)
	}

	var items []list.Item
	for _, meal := range allMeals {
		items = append(items, listItem{title: meal.Name, desc: "Calories: " + strconv.Itoa(meal.Calories)})
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Your meals"

	model := mealsModel{list: list}
	return model
}

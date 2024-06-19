package ui

import (
	"cmdiet/diet"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tableModel struct {
	table table.Model
}

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		}
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func (m tableModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func ViewWeeklyDiet(weeklyDiet []diet.DayDiet) tableModel {
	column := []table.Column{
		{Title: "Day", Width: 16},
		{Title: "Breakfast", Width: 16},
		{Title: "Lunch", Width: 16},
		{Title: "Dinner", Width: 16},
		{Title: "Protein", Width: 16},
		{Title: "Carbs", Width: 16},
		{Title: "Fat", Width: 16},
		{Title: "Total calories", Width: 16},
	}

	rows := make([]table.Row, len(weeklyDiet))
	for i, diet := range weeklyDiet {
		data := table.Row{diet.Day, diet.Breakfast, diet.Lunch, diet.Dinner, strconv.Itoa(diet.Protein), strconv.Itoa(diet.Carbs), strconv.Itoa(diet.Fat), strconv.Itoa(diet.TotalCalories)}
		rows[i] = data
	}

	t := table.New(
		table.WithColumns(column),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	t.SetStyles(s)

	return tableModel{t}
}

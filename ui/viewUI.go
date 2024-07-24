package ui

import (
	"cmdiet/constants"
	"cmdiet/diet"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type summary struct {
	totlaCalories int
	totalProtein  int
	totalCarbs    int
	totalFat      int
}

type tableModel struct {
	table          table.Model
	summary        *summary
	startTimestamp int64
	endTimestamp   int64
	batchCount     int
	quitting       bool
	help           help.Model
}

var baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "h", "left":
			endTimestamp := m.startTimestamp
			startTimestamp := time.UnixMilli(m.startTimestamp).AddDate(0, 0, -m.batchCount).UnixMilli()
			updateTableAndBatchTimeRange(startTimestamp, endTimestamp, &m)
		case "l", "right":
			startTimestamp := m.endTimestamp
			endTimestamp := time.UnixMilli(m.endTimestamp).AddDate(0, 0, m.batchCount).UnixMilli()
			updateTableAndBatchTimeRange(startTimestamp, endTimestamp, &m)
		}
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func updateTableAndBatchTimeRange(startTimestamp int64, endTimestamp int64, m *tableModel) {
	table, err := createTableFromTimestampRange(startTimestamp, endTimestamp)
	if err != nil {
		log.Fatal(err)
	}
	if len(table.Rows()) != 0 {
		m.startTimestamp = startTimestamp
		m.endTimestamp = endTimestamp
		m.table = table
	}
}

func (m tableModel) View() string {
	if m.quitting {
		return ""
	}

	fmtString := "02 Jan"
	from := time.UnixMilli(m.startTimestamp).Format(fmtString)
	to := time.UnixMilli(m.endTimestamp).Format(fmtString)

	var builder strings.Builder
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#32de84")).
		Render(fmt.Sprintf("Your diet logs from %s to %s", from, to))
	builder.WriteString(title)
	builder.WriteString("\n\n" + baseStyle.Render(m.table.View()) + "\n")
	if m.summary != nil {
		fmt.Fprintf(&builder, "Total calories: %d kcal, Total protein: %d grams, Total carbs: %d grams, Total fat: %d grams\n", m.summary.totlaCalories, m.summary.totalProtein, m.summary.totalCarbs, m.summary.totalFat)
	}
	builder.WriteString("\n\n" + m.help.View(constants.ViewTableKeyMap))
	return builder.String()
}

func ViewWeeklyDiet(weeklyDiet []diet.DayDiet, offset int) tableModel {
	startTimestamp := time.Now().AddDate(0, 0, -offset).UnixMilli()
	endTimestamp := time.Now().UnixMilli()
	t := createTable(weeklyDiet)
	summary := BatchedDietsSummary(weeklyDiet)

	return tableModel{
		t, summary, startTimestamp, endTimestamp, offset, false, help.New(),
	}
}

func createTable(weeklyDiet []diet.DayDiet) table.Model {
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
		table.WithHeight(len(weeklyDiet)),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240")).BorderBottom(true).Bold(false)
	s.Selected = s.Selected.Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(false)
	t.SetStyles(s)
	return t
}

func createTableFromTimestampRange(startTimestamp int64, endTimestamp int64) (table.Model, error) {
	weeklyDiet, err := diet.DS.GetBatchedDayDiets(startTimestamp, endTimestamp)
	if err != nil {
		return table.Model{}, fmt.Errorf("couldn't fetch batched diet: %v", err)
	}
	return createTable(weeklyDiet), nil
}

func BatchedDietsSummary(weeklyDiet []diet.DayDiet) *summary {
	var totalCalories, totalProtein, totalCarbs, totalFat int
	for _, diet := range weeklyDiet {
		totalCalories += diet.TotalCalories
		totalProtein += diet.Protein
		totalCarbs += diet.Carbs
		totalFat += diet.Fat
	}

	summary := &summary{
		totlaCalories: totalCalories,
		totalProtein:  totalProtein,
		totalCarbs:    totalCarbs,
		totalFat:      totalFat,
	}

	return summary
}

package ui

import (
	"cmdiet/diet"
	"log"
	"time"

	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type evalModel struct {
	chart                 tslc.Model
	quitting              bool
	batchCount            int
	startTimestamp        int64
	endTimestamp          int64
	earliestDietTimestamp int64
}

func (m evalModel) Init() tea.Cmd {
	m.chart.DrawXYAxisAndLabel()
	return nil
}

func (m evalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "h", "left":
			endTimestamp := m.startTimestamp
			startTimestamp := time.UnixMilli(m.startTimestamp).AddDate(0, 0, -m.batchCount).UnixMilli()

			if time.UnixMilli(endTimestamp).Before(time.UnixMilli(m.earliestDietTimestamp)) {
				return m, nil
			}

			updateTimestamps(startTimestamp, endTimestamp, &m)
			updateChartAndBatchTimeRange(startTimestamp, endTimestamp, &m.chart)
		case "l", "right":
			startTimestamp := m.endTimestamp
			endTimestamp := time.UnixMilli(m.endTimestamp).AddDate(0, 0, m.batchCount).UnixMilli()

			if time.UnixMilli(startTimestamp).After(time.Now()) {
				return m, nil
			}
			updateTimestamps(startTimestamp, endTimestamp, &m)
			updateChartAndBatchTimeRange(startTimestamp, endTimestamp, &m.chart)
		}
	}

	m.chart.DrawBraille()
	return m, nil
}

func updateTimestamps(startTimestamp int64, endTimestamp int64, m *evalModel) {
	m.startTimestamp = startTimestamp
	m.endTimestamp = endTimestamp
}

func updateChartAndBatchTimeRange(startTimestamp, endTimestamp int64, chart *tslc.Model) {
	batchDiet, err := diet.DS.GetBatchedDayDiets(startTimestamp, endTimestamp)
	if err != nil {
		log.Fatal(err)
	}

	chart.SetTimeRange(time.UnixMilli(startTimestamp), time.UnixMilli(endTimestamp))
	chart.SetViewTimeRange(time.UnixMilli(startTimestamp), time.UnixMilli(endTimestamp))
	chart.ClearAllData()
	chart.Clear()
	chart.DrawXYAxisAndLabel()

	for _, diet := range batchDiet {
		date, err := time.Parse(time.DateOnly, diet.Day)
		if err != nil {
			log.Fatal(err)
		}
		chart.Push(tslc.TimePoint{date, float64(diet.TotalCalories)})
	}
}

func (m evalModel) View() string {
	if m.quitting {
		return ""
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#32de84")).
		Render("Here's your calories graph for the last week") + "\n" +
		lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Render(m.chart.View())
}

func NewEvalModel() evalModel {
	width := 80
	height := 20
	chart := tslc.New(width, height)

	// set default data set line color to red
	chart.SetStyle(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")), // red
	)

	offset := 7
	startTimestamp := time.Now().AddDate(0, 0, -offset).UnixMilli()
	endTimestamp := time.Now().UnixMilli()
	updateChartAndBatchTimeRange(startTimestamp, endTimestamp, &chart)
	earliestDietTimestamp := diet.DS.EarliestDietTimestamp()

	chart.Focus()
	return evalModel{
		chart:                 chart,
		startTimestamp:        startTimestamp,
		endTimestamp:          endTimestamp,
		batchCount:            offset,
		earliestDietTimestamp: earliestDietTimestamp,
	}
}

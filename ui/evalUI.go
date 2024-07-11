package ui

import (
	"cmdiet/diet"
	"fmt"
	"log"
	"time"

	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type evalModel struct {
	chart          tslc.Model
	quitting       bool
	batchCount     int
	startTimestamp int64
	endTimestamp   int64
}

func (m evalModel) Init() tea.Cmd {
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
			fmt.Println("pressed left")
			endTimestamp := m.startTimestamp
			startTimestamp := time.UnixMilli(m.startTimestamp).AddDate(0, 0, -m.batchCount).UnixMilli()
			updateChartAndBatchTimeRange(startTimestamp, endTimestamp, &m)
		case "l", "right":
			fmt.Println("pressed left")
			startTimestamp := m.endTimestamp
			endTimestamp := time.UnixMilli(m.endTimestamp).AddDate(0, 0, m.batchCount).UnixMilli()
			updateChartAndBatchTimeRange(startTimestamp, endTimestamp, &m)
		}
	}

	var cmd tea.Cmd
	m.chart, cmd = m.chart.Update(msg)
	m.chart.DrawBrailleAll()
	return m, cmd
}

func updateChartAndBatchTimeRange(startTimestamp, endTimestamp int64, m *evalModel) {
	batchDiet, err := diet.DS.GetBatchedDayDiets(startTimestamp, endTimestamp)
	if err != nil {
		log.Fatal(err)
	}
	if len(batchDiet) == 0 {
		m.startTimestamp = startTimestamp
		m.endTimestamp = endTimestamp
		pushDietsToChartDataset(&m.chart, batchDiet)
	}
}

func (m evalModel) View() string {
	if m.quitting {
		return ""
	}
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Render(m.chart.View())
}

func NewEvalModel() evalModel {
	width := 80
	height := 20
	chart := tslc.New(width, height)

	diets, err := diet.DS.GetDayDietsByOffset(17, 0)
	if err != nil {
		log.Fatal(err)
	}

	pushDietsToChartDataset(&chart, diets)

	// set default data set line color to red
	chart.SetStyle(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")), // red
	)

	offset := 7
	startTimestamp := time.Now().AddDate(0, 0, -offset).UnixMilli()
	endTimestamp := time.Now().UnixMilli()

	return evalModel{
		chart:          chart,
		startTimestamp: startTimestamp,
		endTimestamp:   endTimestamp,
		batchCount:     offset,
	}
}

func pushDietsToChartDataset(chart *tslc.Model, diets []diet.DayDiet) {
	chart.ClearAllData()
	for _, diet := range diets {
		date, err := time.Parse(time.DateOnly, diet.Day)
		if err != nil {
			log.Fatal(err)
		}
		chart.Push(tslc.TimePoint{date, float64(diet.TotalCalories)})
	}
}

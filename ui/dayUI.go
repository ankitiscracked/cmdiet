package ui

import (
	"cmdiet/constants"
	"cmdiet/diet"
	"cmdiet/meals"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	purple = lipgloss.Color("99")
	pink   = lipgloss.Color("205")
)

type totalMacros struct {
	totalCalories int
	totalProtein  int
	totalCarbs    int
	totalFats     int
}

type dayModel struct {
	diets       []diet.Diet
	help        help.Model
	currentTime time.Time
}

func (m dayModel) Init() tea.Cmd {
	return nil
}

func (m dayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h", "left":
			return NewDayModel(m.currentTime.AddDate(0, 0, -1)), nil
		case "l", "right":
			newTime := m.currentTime.AddDate(0, 0, 1)
			if newTime.Before(time.Now()) {
				return NewDayModel(newTime), nil
			}
		}
	}

	return m, nil
}

func (m dayModel) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		getDietsSummary(&m),
		lipgloss.NewStyle().MarginTop(2).Faint(true).Render(m.help.ShortHelpView(constants.DayViewKeyMap)),
	)
}

func NewDayModel(currentTime time.Time) dayModel {
	diets, err := diet.DS.GetDietsForDay(currentTime)
	if err != nil {
		log.Fatal(err)
	}

	return dayModel{currentTime: currentTime, diets: diets}
}

func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

func getDietsSummary(m *dayModel) string {
	var today string
	if isToday(m.currentTime) {
		today = "today"
	} else {
		today = m.currentTime.Format("02 Jan")
	}

	if len(m.diets) == 0 {
		return "\n" + TerminalInfo("No diet logs for "+today)
	}

	macros := getTotalMacros(m.diets)
	header := lipgloss.NewStyle().Foreground(pink).Bold(true).Render("Today's Diet Summary")
	return lipgloss.JoinVertical(lipgloss.Left, "\n"+header, "\n"+getDietSummary(m.diets), "\n"+getMacrosSummary(macros)+"\n", getMacrosDistribution(macros))
}

func getMacrosDistribution(macros totalMacros) string {
	totalProtein := macros.totalProtein
	totalCarbs := macros.totalCarbs
	totalFat := macros.totalFats

	proteinColor := lipgloss.Color("#fa84d7")
	carbsColor := lipgloss.Color("#84faab")
	fatsColor := lipgloss.Color("#6ba9e3")

	macroBarStyle := func(macro int, color lipgloss.Color) string {
		totalCalories := totalProtein*4 + totalCarbs*4 + totalFat*9
		if totalCalories == 0 {
			return ""
		}
		width := int(float64(macro) / float64(totalCalories) * 100)
		return lipgloss.NewStyle().Background(color).Width(width).Render()
	}

	return lipgloss.NewStyle().Foreground(pink).Bold(true).Render("Calories by macro") +
		"\n\n" +
		lipgloss.JoinHorizontal(
			lipgloss.Bottom,
			macroBarStyle(totalProtein*4, proteinColor),
			macroBarStyle(totalCarbs*4, carbsColor),
			macroBarStyle(totalFat*9, fatsColor),
		)
}

func getMacrosSummary(macros totalMacros) string {
	title := lipgloss.NewStyle().Foreground(purple)

	result := []string{
		lipgloss.NewStyle().Foreground(pink).Bold(true).Render("Macros breakdown\n"),
		title.Render("Total Calories: ") + strconv.Itoa(macros.totalCalories),
		title.Render("Total Protein: ") + strconv.Itoa(macros.totalProtein),
		title.Render("Total Carbs: ") + strconv.Itoa(macros.totalCarbs),
		title.Render("Total Fats: ") + strconv.Itoa(macros.totalFats),
	}

	return strings.Join(result, "\n")
}

func getTotalMacros(diets []diet.Diet) totalMacros {
	var totalCalories, totalProtein, totalCarbs, totalFat int
	for _, diet := range diets {
		meal, err := meals.MS.GetMeal(int(diet.MealId))
		if err != nil {
			log.Fatal(err)
		}

		totalCalories += int(meal.GetTotalCalories())
		totalProtein += int(meal.GetMealMacro(meals.Protein))
		totalCarbs += int(meal.GetMealMacro(meals.Carbs))
		totalFat += int(meal.GetMealMacro(meals.Fat))
	}
	return totalMacros{totalCalories, totalProtein, totalCarbs, totalFat}
}

func getMealRow(meal meals.Meal, mealType string) string {
	headers := []string{"Type", "Meal", "Calories", "Proteins", "Carbs", "Fats"}
	return table.New().
		Border(lipgloss.ThickBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style
			if row == 0 {
				style = lipgloss.NewStyle().Padding(0, 1).Foreground(purple).Bold(true).Align(lipgloss.Center)
			} else {
				return lipgloss.NewStyle().Padding(0, 1).Align(lipgloss.Center)
			}
			return style
		}).
		Headers(headers...).
		Row(mealType, meal.Name, fmt.Sprintf("%d", meal.GetTotalCalories()), fmt.Sprintf("%d", meal.GetMealMacro(meals.Protein)), fmt.Sprintf("%d", meal.GetMealMacro(meals.Carbs)), fmt.Sprintf("%d", meal.GetMealMacro(meals.Fat))).
		String()
}

func getDietSummary(diets []diet.Diet) string {
	var tables []string
	for _, diet := range diets {
		meal, err := meals.MS.GetMeal(int(diet.MealId))
		if err != nil {
			log.Fatal(err)
		}

		switch diet.MealType {
		case "breakfast":
			tables = append(tables, getMealRow(meal, "Breakast"))
		case "lunch":
			tables = append(tables, getMealRow(meal, "Lunch"))
		case "dinner":
			tables = append(tables, getMealRow(meal, "Dinner"))
		default:
		}

	}
	return lipgloss.JoinVertical(lipgloss.Left, tables...)
}

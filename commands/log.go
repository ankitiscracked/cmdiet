package commands

import (
	"cmdiet/diet"
	"cmdiet/ui"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var daysInPast int

var logCmd = &cobra.Command{
	Use:   "log [breakfast | lunch | snacks | dinner]",
	Short: "Log your diet for breakfast, lunch, or afternoon.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mealType := args[0]
		res, err := diet.GetMealType(mealType)
		if err != nil {
			log.Fatal("Invalid meal type. Please use breakfast, lunch, snacks, or dinner.")
		}

		logged, err := diet.DS.MealTypeLoggedForToday(res)
		if err != nil {
			log.Fatal(err)
		}

		if logged {
			log.Fatal("You have already logged your " + mealType + " for today.")
		}

		if daysInPast < 0 {
			log.Fatal("You can't log diets for future dates, duh!")
		}

		var logForDate time.Time
		if daysInPast != 0 {
			logForDate = time.Now().AddDate(0, 0, -daysInPast)
		} else {
			logForDate = time.Now()
		}

		program := tea.NewProgram(ui.NewDietModel(res, logForDate))
		if _, err := program.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	logCmd.Flags().IntVarP(&daysInPast, "pdays", "d", 0, "Number of days in the past to log")
	RootCmd.AddCommand(logCmd)
}

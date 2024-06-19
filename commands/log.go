package commands

import (
	"cmdiet/diet"
	"cmdiet/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func validMealType(mealType string) bool {
	mealTypes := []string{"breakfast", "lunch", "snacks", "dinner"}

	for _, meal := range mealTypes {
		if meal == mealType {
			return true
		}
	}
	return false
}

var logCmd = &cobra.Command{
	Use:   "log [meal type]",
	Short: "Log your diet for breakfast, lunch, or afternoon.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mealType := args[0]
		if !validMealType(mealType) {
			log.Fatal("Invalid meal type. Please use breakfast, lunch, snacks, or dinner.")
		}
		if diet.DefaultDietService.MealTypeLoggedForToday(mealType) {
			log.Fatal("You have already logged your " + mealType + " for today.")
		}
		program := tea.NewProgram(ui.NewDietModel(mealType))
		if _, err := program.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(logCmd)
}

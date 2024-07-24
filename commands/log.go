package commands

import (
	"cmdiet/diet"
	"cmdiet/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

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

		program := tea.NewProgram(ui.NewDietModel(mealType))
		if _, err := program.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(logCmd)
}

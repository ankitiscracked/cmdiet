package commands

import (
	"github.com/spf13/cobra"
)

var addMealCmd = &cobra.Command{
	Use:   "add-meal [meal type]",
	Short: "Add a meal to your catalog",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// program := tea.NewProgram(meals.MS.AddMeal(mealType))
		// if _, err := program.Run(); err != nil {
		// log.Fatal(err)
		// }
	},
}

func init() {
	RootCmd.AddCommand(addMealCmd)
}

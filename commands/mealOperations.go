package commands

import (
	"cmdiet/meals"
	"cmdiet/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your logged meals",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		program := tea.NewProgram(ui.ViewMealsModel())
		if _, err := program.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

var addMealCmd = &cobra.Command{
	Use:   "add-meal",
	Short: "Add a new meal",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		program := tea.NewProgram(ui.NewAddMealModel())
		if _, err := program.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

var addComponentCmd = &cobra.Command{
	Use:   "add-component",
	Short: "Add a new meal component",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var componentType meals.MealComponentType
		fixed, _ := cmd.Flags().GetBool("fixed")
		variable, _ := cmd.Flags().GetBool("variable")

		if fixed {
			componentType = meals.FixedMealComponentType
		} else if variable {
			componentType = meals.VariableMealComponentType
		}

		program := tea.NewProgram(ui.NewAddMealComponentModel(componentType))
		if _, err := program.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

var mealsCmd = &cobra.Command{
	Use:   "meals",
	Short: "List all your logged meals",
}

func init() {
	fixed := "fixed"
	variable := "variable"
	addComponentCmd.Flags().BoolP(fixed, "f", false, "Add a fixed meal component")
	addComponentCmd.Flags().BoolP(variable, "v", false, "Add a variable meal component")
	addComponentCmd.MarkFlagsOneRequired(fixed, variable)
	addComponentCmd.MarkFlagsMutuallyExclusive(fixed, variable)

	mealsCmd.AddCommand(listCmd)
	mealsCmd.AddCommand(addMealCmd)
	mealsCmd.AddCommand(addComponentCmd)
	RootCmd.AddCommand(mealsCmd)
}

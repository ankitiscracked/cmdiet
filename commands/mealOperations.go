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

var addCmd = &cobra.Command{
	Use:   "add [meal|component]",
	Short: "Add a new meal or meal component",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var program *tea.Program
		switch args[0] {
		case "meal":
			program = tea.NewProgram(ui.NewAddMealModel())
		case "component":
			var componentType meals.MealComponentType
			fixed, _ := cmd.Flags().GetBool("fixed")
			variable, _ := cmd.Flags().GetBool("variable")

			if fixed && variable {
				log.Fatal("Please specify either --fixed or --variable, not both.")
			} else if fixed {
				componentType = meals.FixedMealComponentType
			} else if variable {
				componentType = meals.VariableMealComponentType
			} else {
				log.Fatal("Please specify either --fixed or --variable flag.")
			}
			program = tea.NewProgram(ui.NewAddMealComponentModel(componentType))
		default:
			log.Fatal("Invalid argument. Use 'meal' or 'component'.")
		}
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
	mealsCmd.AddCommand(listCmd)
	mealsCmd.AddCommand(addCmd)
	RootCmd.AddCommand(mealsCmd)
}

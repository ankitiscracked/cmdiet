package commands

import (
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

var mealsCmd = &cobra.Command{
	Use:   "meals",
	Short: "List all your logged meals",
}

func init() {
	mealsCmd.AddCommand(listCmd)
	RootCmd.AddCommand(mealsCmd)
}

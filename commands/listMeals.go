package commands

import (
	"cmdiet/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var listMealsCmd = &cobra.Command{
	Use:   "list-meals",
	Short: "List all your logged meals",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		program := tea.NewProgram(ui.ViewMealsModel())
		if _, err := program.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(listMealsCmd)
}

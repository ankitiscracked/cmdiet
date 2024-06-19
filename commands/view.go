package commands

import (
	"cmdiet/diet"
	"cmdiet/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view [date]",
	Short: "View the diet breakdown for a given day.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		weekDiets := diet.DefaultDietService.GetLastWeekDiet()
		program := tea.NewProgram(ui.ViewWeeklyDiet(weekDiets))
		if _, err := program.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(viewCmd)
}

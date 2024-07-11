package commands

import (
	"cmdiet/ui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var evalCmd = &cobra.Command{
	Use:   "eval [calories | carbs | protein | fat]",
	Short: "Evaluate different aspects of your diet.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := tea.NewProgram(ui.NewEvalModel(), tea.WithAltScreen()).Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(evalCmd)
}

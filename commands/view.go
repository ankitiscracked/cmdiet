package commands

import (
	"cmdiet/diet"
	"cmdiet/ui"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	offset   int
	week     int // flag for --week
	month    int // flag for --month
	mealType int // flag for meal type
	macros   int // flag for macros
	today    bool
)

var viewCmd = &cobra.Command{
	Use:   "view [date]",
	Short: "View the diet breakdown for a given day.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 && args[0] == "today" {
			program := tea.NewProgram(ui.NewDayModel(time.Now()))
			if _, err := program.Run(); err != nil {
				log.Fatal(err)
			}
			return
		}

		var finalOffset int
		if offset != 0 {
			finalOffset = offset
		} else {
			finalOffset = 7
		}

		weekDiets, err := diet.DS.GetDayDietsByOffset(finalOffset, 0)
		if err != nil {
			log.Fatal(err)
		}
		program := tea.NewProgram(ui.ViewWeeklyDiet(weekDiets, finalOffset))
		if _, err := program.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	viewCmd.Flags().IntVarP(&offset, "offset", "o", 7, "Number of days to view")
	viewCmd.Flags().BoolVarP(&today, "today", "t", false, "View today's log summary")
	viewCmd.ValidArgs = []string{"today"}
	RootCmd.AddCommand(viewCmd)
}

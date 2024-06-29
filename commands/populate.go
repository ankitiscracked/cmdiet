package commands

import (
	"cmdiet/database"
	"fmt"

	"github.com/spf13/cobra"
)

var populateCmd = &cobra.Command{
	Use:   "populate [diets|meals]",
	Short: "Populate the fake database with fake data for diets and meals",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !fake {
			fmt.Println("You must provide the --fake flag to populate the database with fake data")
			return
		}

		if args[0] == "diets" {
			database.PopulateFakeDiets()
		} else if args[0] == "meals" {
			database.PopulateFakeMeals()
		}
	},
}

func init() {
	RootCmd.AddCommand(populateCmd)
}

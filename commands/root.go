package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "diet",
	Short: "Diet is a simple CLI to log your meals",
}

func Execute() error {
	return RootCmd.Execute()
}

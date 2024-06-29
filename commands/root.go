package commands

import (
	"github.com/spf13/cobra"
)

var fake bool

var RootCmd = &cobra.Command{
	Use:   "diet",
	Short: "Diet is a simple CLI to log your meals",
}

func init() {
	RootCmd.PersistentFlags().BoolVar(&fake, "fake", false, "Use this flag to initialize a fake database")
}

func Execute() error {
	return RootCmd.Execute()
}

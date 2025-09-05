package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "auth",
	Short: "Auth Service",
}

func init() {
	rootCmd.AddCommand(serverCommand)
	rootCmd.AddCommand(migrateCommand)
}

func Execute() error {
	return rootCmd.Execute()
}

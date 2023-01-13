package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd *cobra.Command = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of accicalc",
	Long:  `Print the version number of accicalc.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("accicalc %v (%v)\n", Version, Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

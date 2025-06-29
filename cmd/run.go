package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run ClipTree in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ClipTree run executed (placeholder)")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

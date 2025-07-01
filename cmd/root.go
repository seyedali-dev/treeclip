// Package cmd houses all CLI commands and flags.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "treeclip",
	Short: "treeclip – copy directory contents to clipboard and editor",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

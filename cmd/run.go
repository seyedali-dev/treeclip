package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

// runCmd concatenates the contents of all files in a given directory and writes them to a text file.
var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Traverse a folder and output all file contents into a .txt file",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine root path to walk
		var rootDir string
		if len(args) > 0 {
			var err error
			rootDir, err = filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("invalid path: %w", err)
			}
		} else {
			// Default to current directory
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			rootDir = cwd
		}

		// Create output file in CWD
		outputFilePath := "cliptree_output.txt"
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func(outputFile *os.File) {
			err := outputFile.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close temp file: %v\n", err)
			}
		}(outputFile)

		fmt.Printf("Writing concatenated contents to: %s\n", outputFilePath)

		// Walk directory recursively
		err = filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			// Calculate relative path to display
			relPath, _ := filepath.Rel(rootDir, path)
			fmt.Fprintf(outputFile, "==> ./%s\n", relPath)

			// Open file and copy content
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to close temp file: %v\n", err)
				}
			}(f)

			_, err = io.Copy(outputFile, f)
			if err != nil {
				return err
			}
			fmt.Fprintln(outputFile) // Add newline between files

			return nil
		})
		if err != nil {
			return fmt.Errorf("error while traversing: %w", err)
		}

		fmt.Println("✅ File contents written successfully.")
		return nil
	},
}

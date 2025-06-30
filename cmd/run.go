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

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Traverse a directory and output all file contents into a .txt file",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		tempFile, err := os.CreateTemp(rootDir, "cliptree_*.txt")
		if err != nil {
			return fmt.Errorf("failed to create temp file: %w", err)
		}
		defer func() {
			if err := tempFile.Close(); err != nil {
				safeFprintf(os.Stderr, "Warning: failed to close temp file: %v\n", err)
			}
		}()

		fmt.Printf("Creating output file at: %s\n", tempFile.Name())

		err = filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error accessing path %s: %w", path, err)
			}
			if d.IsDir() {
				return nil
			}

			// Write file path
			relPath, err := filepath.Rel(rootDir, path)
			if err != nil {
				return fmt.Errorf("error getting relative path for %s: %w", path, err)
			}

			safeFprintf(tempFile, "==> ./%s\n", relPath)

			// Copy file content
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open %s: %w", path, err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					safeFprintf(os.Stderr, "Warning: failed to close %s: %v\n", path, err)
				}
			}()

			_, err = io.Copy(tempFile, f)
			if err != nil {
				return err
			}
			safeFprintln(tempFile, "\n") // separate files

			return nil
		})

		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		fmt.Println("File contents written successfully.")
		return nil
	},
}

// safeFprintf is a wrapper around fmt.Fprintf that panics if the write fails.
func safeFprintf(file io.Writer, format string, a ...any) {
	_, err := fmt.Fprintf(file, format, a...)
	if err != nil {
		panic(fmt.Sprintf("Warning: failed to write to temp file: {%v}\n", err))
	}
}

// safeFprintln is a wrapper around fmt.Fprintln that panics if the write fails.
func safeFprintln(file io.Writer, a ...any) {
	_, err := fmt.Fprintln(file, a...)
	if err != nil {
		panic(fmt.Sprintf("Warning: failed to write to temp file: {%v}\n", err))
	}
}

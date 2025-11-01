// Package cmd. runCmd handles directory traversal and file concatenation.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	fileUtils "github.com/seyedali-dev/treeclip/pkg/utils"

	"github.com/seyedali-dev/treeclip/internal/clipboard"
	"github.com/seyedali-dev/treeclip/internal/editor"
	"github.com/seyedali-dev/treeclip/internal/exclude"
	"github.com/seyedali-dev/treeclip/internal/traversal"
	"github.com/spf13/cobra"
)

// outputFile is the temporary traversed file created via treeclip.
const outputFile = "treeclip_temp.txt"

var (
	excludePatterns    []string
	clipboardEnabled   bool
	showClipboardStats bool
	editorEnabled      bool
	deleteAfterEditor  bool
)

func init() {
	runCmd.Flags().StringSliceVarP(&excludePatterns, "exclude", "e", []string{}, "Exclude files/folders matching these patterns (can be used multiple times)")
	runCmd.Flags().BoolVarP(&clipboardEnabled, "clipboard", "c", true, "Copy output to clipboard")
	runCmd.Flags().BoolVar(&showClipboardStats, "stats", false, "Show clipboard content statistics")
	runCmd.Flags().BoolVarP(&editorEnabled, "editor", "o", false, "Open output file in the default text editor")
	runCmd.Flags().BoolVarP(&deleteAfterEditor, "delete", "d", true, "Delete the output file after editor is closed")

	rootCmd.AddCommand(runCmd)
}

// runCmd concatenates the contents of all files in a given directory and writes them to a text file.
var runCmd = &cobra.Command{
	Use:   "run [path | cwd if empty]",
	Short: "Traverse a folder and output all file contents into a .txt file",
	Long:  generateLongDescription(),
	Args:  cobra.MaximumNArgs(1),
	RunE:  registerRunCmd(),
}

func generateLongDescription() string {
	return `Traverse a folder recursively and output all file contents into a .txt file.
    
Examples:
  treeclip run                                     # Current directory, copy to clipboard
  treeclip run /path/to/dir                        # Specific directory
  treeclip run --exclude "*.log" --exclude "*.tmp" # Exclude patterns
  treeclip run -e "*.md" -e "folder1" -e "app.go"  # Multiple exclusions
  treeclip run --stats                             # Show content statistics
  treeclip run --editor                            # Open output file in the default text editor
  treeclip run --delete                            # Delete the output file after editor is closed`
}

// registerRunCmd handles the actual logic for treeclip dir traversal.
func registerRunCmd() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Determine root path
		rootDir, err := determineRootDir(args)
		if err != nil {
			return err
		}

		// Create output file
		outF, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		fileUtils.WriteDataLn(outF, "// 💡Paths are displayed in Unix-style format (forward slashes)")

		// Load exclusions
		ignoreFilePatterns, err := exclude.LoadIgnorePatterns(rootDir)
		if err != nil {
			return err
		}
		allEx := append(excludePatterns, ignoreFilePatterns...)
		allEx = append(allEx, exclude.DefaultExclusions...)

		// Traverse and write
		filesProcessed, filesSkipped, err := traversal.TraverseDir(rootDir, allEx, outF)
		if err != nil {
			return err
		}

		fileUtils.SafeCloseFile(outF)

		// Clipboard
		if err := clipboard.HandleClipboardCommandFlag(clipboardEnabled, showClipboardStats, outputFile); err != nil {
			return err
		}

		// Editor
		if err := editor.HandleEditorCommandFlag(editorEnabled, deleteAfterEditor, outputFile); err != nil {
			return err
		}

		fmt.Printf("\n------------ (●'◡'●) ------------\n")
		fmt.Printf("🎉  Process completed! ＼(＾▽＾)／\n")
		fmt.Printf("📊  Files processed: %d (•̀ᴗ•́)و\n", filesProcessed)
		fmt.Printf("🚫  Files/folders skipped: %d (；一_一)\n", filesSkipped)
		fmt.Printf("📄  Output file: %s (ᵔ◡ᵔ)\n", outputFile)
		fmt.Println("\n  totoro!  ㄟ( ▔, ▔ )ㄏ")
		return nil
	}
}

// determineRootDir determines the root directory to traverse to.
func determineRootDir(args []string) (string, error) {
	rootDir := "."
	if len(args) > 0 {
		rootPathArg, err := filepath.Abs(args[0])
		if err != nil {
			return "", fmt.Errorf("invalid path: %w", err)
		}
		rootDir = rootPathArg
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get cwd: %w", err)
		}
		rootDir = cwd
	}
	return rootDir, nil
}

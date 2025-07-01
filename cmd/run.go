package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// excludePatterns holds the patterns to exclude during directory traversal
var excludePatterns []string

func init() {
	rootCmd.AddCommand(runCmd)

	// Add the --exclude flag that can be used multiple times
	runCmd.Flags().StringSliceVarP(&excludePatterns, "exclude", "e", []string{},
		"Exclude files/folders matching these patterns (can be used multiple times)")
}

// runCmd concatenates the contents of all files in a given directory and writes them to a text file.
var runCmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Traverse a folder and output all file contents into a .txt file",
	Long: `Traverse a folder recursively and output all file contents into a .txt file.
    
Examples:
  treeclip run                                     # Current directory
  treeclip run /path/to/dir                        # Specific directory
  treeclip run --exclude "*.log" --exclude "*.tmp" # Exclude patterns
  treeclip run -e "*.md" -e "folder1" -e "app.go"  # Multiple exclusions`,
	Args: cobra.MaximumNArgs(1),
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

		// Validate that the root directory exists
		if _, err := os.Stat(rootDir); os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", rootDir)
		}

		// Create output file in CWD
		outputFilePath := "treeclip_output.txt"
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func(outputFile *os.File) {
			err := outputFile.Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close output file: %v\n", err)
			}
		}(outputFile)

		// Add default exclusions to prevent infinite loops and common unwanted files
		defaultExclusions := []string{
			"treeclip_output.txt", // Our own output file. Prevent recursion!
			"*.tmp",
			"*.temp",
			"*.exe",
			"*.sh",
			".git",
			".idea",
			".DS_Store",
			"Thumbs.db",
		}

		// Combine user patterns with default exclusions
		allExcludePatterns := append(excludePatterns, defaultExclusions...)

		fmt.Printf("🔍 Scanning directory: %s\n", rootDir)
		if len(excludePatterns) > 0 {
			fmt.Printf("🚫 User exclusions: %v\n", excludePatterns)
		}
		fmt.Printf("🛡️  Default exclusions: %v\n", defaultExclusions)
		fmt.Printf("📄 Writing concatenated contents to: %s\n", outputFilePath)

		var filesProcessed int
		var filesSkipped int

		// Walk directory recursively
		err = filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Calculate relative path for pattern matching and display
			relPath, err := filepath.Rel(rootDir, path)
			if err != nil {
				return err
			}

			// Check if current path should be excluded (using combined patterns)
			if shouldExclude(relPath, d.Name(), d.IsDir(), allExcludePatterns) {
				filesSkipped++
				if d.IsDir() {
					fmt.Printf("⏭️  Skipping directory: %s\n", relPath)
					return filepath.SkipDir // Skip entire directory
				}
				fmt.Printf("⏭️  Skipping file: %s\n", relPath)
				return nil
			}

			// Skip directories (we only want to process files)
			if d.IsDir() {
				return nil
			}

			filesProcessed++
			fmt.Printf("📖 Processing: %s\n", relPath)

			// Write file header with relative path
			fmt.Fprintf(outputFile, "==> ./%s\n", relPath)

			// Open file and copy content
			f, err := os.Open(path)
			if err != nil {
				fmt.Printf("⚠️  Warning: failed to open %s: %v\n", relPath, err)
				fmt.Fprintf(outputFile, "[ERROR: Could not read file - %v]\n\n", err)
				return nil // Continue processing other files
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to close file %s: %v\n", path, err)
				}
			}(f)

			// Copy file content to output
			_, err = io.Copy(outputFile, f)
			if err != nil {
				fmt.Printf("⚠️  Warning: failed to copy content from %s: %v\n", relPath, err)
				fmt.Fprintf(outputFile, "[ERROR: Could not copy file content - %v]\n", err)
			}

			// Add separator between files
			fmt.Fprintln(outputFile)
			fmt.Fprintln(outputFile)

			return nil
		})

		if err != nil {
			return fmt.Errorf("error while traversing directory: %w", err)
		}

		fmt.Printf("\n✅ File contents written successfully!\n")
		fmt.Printf("📊 Files processed: %d\n", filesProcessed)
		fmt.Printf("🚫 Files/folders skipped: %d\n", filesSkipped)
		return nil
	},
}

// shouldExclude checks if a file or directory should be excluded based on the exclude patterns.
// It supports:
// - Exact filename/dirname matches (e.g., "app.go", "folder1")
// - Wildcard patterns (e.g., "*.log", "*.md")
// - Relative path matches (e.g., "src/*.go")
func shouldExclude(relPath, name string, isDir bool, patterns []string) bool {
	for _, pattern := range patterns {
		// Clean the pattern to handle different input formats
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		// Check exact name match (for both files and directories)
		if name == pattern {
			return true
		}

		// Check exact relative path match
		if relPath == pattern {
			return true
		}

		// Check wildcard pattern match against filename
		if matched, err := filepath.Match(pattern, name); err == nil && matched {
			return true
		}

		// Check wildcard pattern match against relative path
		if matched, err := filepath.Match(pattern, relPath); err == nil && matched {
			return true
		}

		// For directories, also check if the pattern matches any parent directory
		if isDir {
			// Check if any part of the relative path matches the pattern
			pathParts := strings.Split(relPath, string(filepath.Separator))
			for _, part := range pathParts {
				if part == pattern {
					return true
				}
				if matched, err := filepath.Match(pattern, part); err == nil && matched {
					return true
				}
			}
		}

		// Handle patterns that might include path separators
		// e.g., "src/*.go" should match files in src directory
		if strings.Contains(pattern, string(filepath.Separator)) {
			if matched, err := filepath.Match(pattern, relPath); err == nil && matched {
				return true
			}
		}
	}

	return false
}

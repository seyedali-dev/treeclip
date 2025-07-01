// TreeClip is a CLI tool that recursively traverses directories, concatenates file contents,
// and makes them available in your clipboard and text editor.
//
// Key Features:
//
// - Recursive directory traversal with exclusion patterns
//
// - Combined output with file headers
//
// - Clipboard integration
//
// - Temporary file handling with auto-cleanup
//
// - Configurable via command flags and .treeclipignore
//
// Note: some features are to be implemented for future and new features will be added,
// as the project continues.
package main

import "github.com/seyedali-dev/treeclip/cmd"

func main() {
	cmd.Execute()
}

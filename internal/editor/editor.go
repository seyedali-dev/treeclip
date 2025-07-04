// Package editor - editor provides function for opening and handling the traversed data in OS's default editor.
package editor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// HandleEditorCommandFlag handles clipboardEnabled command flag.
func HandleEditorCommandFlag(editorFlag, deleteAfterEditorFlag bool, outputFilePath string) error {
	if editorFlag {
		fmt.Println("\n📝  Opening file in default text editor... (◠‿◠)✎")
		if deleteAfterEditorFlag {
			fmt.Println("⚠️  Warning! Will delete the temporary file after editor closes (×_×)⌒☆")
		}

		err := openInEditor(outputFilePath)
		if err != nil {
			return fmt.Errorf("⚠️  Warning: failed to open editor: %w\n", err)
		} else {
			fmt.Println("✅  Editor closed. Proceeding... (ﾉ´ヮ`)ﾉ*: ･ﾟ")

			if deleteAfterEditorFlag {
				fmt.Println()
				fmt.Println("\n🗑️  Attempting to delete the temp file (⋟﹏⋞)")
				err := os.Remove(outputFilePath)
				if err != nil {
					fmt.Printf("⚠️  Warning: failed to delete file: %v\n", err)
				} else {
					fmt.Printf("🧽  Output temp file deleted: %s (￣ω￣)\n", outputFilePath)
				}
			}
		}
	}
	return nil
}

// openInEditor opens the given file in the system's default text editor and waits for it to close.
func openInEditor(filePath string) error {
	time.Sleep(100 * time.Millisecond)
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-W", filePath) //TODO: test me!
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", "/WAIT", filePath)
	default: // Linux and others
		cmd = exec.Command("xdg-open", filePath) //TODO: test me!
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

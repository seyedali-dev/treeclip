// Package utils. file_utils provides utility functions for working with file.
package utils

import (
	"fmt"
	"os"
)

// SafeCloseFile closes the file safely by panicking on error.
func SafeCloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		panic(fmt.Sprintf("⚠️  Warning! failed to close opened file: %v ╰（‵□′）╯", err))
	}
}

// WriteData writes the provided data to the end of file.
func WriteData(file *os.File, data string) {
	if _, err := fmt.Fprintf(file, data); err != nil {
		panic(fmt.Sprintf("❌🪲  [ERROR] failed to write data to file %s: %v (╯°□°）╯︵ ┻━┻", file.Name(), err))
	}
}

// WriteDataLn writes the provided data to the end of file and creates a new line.
func WriteDataLn(file *os.File, data string) {
	WriteData(file, data+"\n")
}

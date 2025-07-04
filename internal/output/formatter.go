// Package output. formatter provides file header and separator formatting.
package output

import (
	"fmt"
	"io"
)

// WriteHeader adds the file header to writer.
func WriteHeader(file io.Writer, relPath string) {
	if _, err := fmt.Fprintf(file, "==> %s\n", relPath); err != nil {
		panic(fmt.Sprintf("❌🪲  [ERROR] failed to write to file: %v", err))
	}
}

// WriteSeparator writes an empty line as separator.
func WriteSeparator(file io.Writer) {
	if _, err := fmt.Fprintln(file); err != nil {
		panic(fmt.Sprintf("❌🪲  [ERROR] failed to write to file: %v", err))
	}
}

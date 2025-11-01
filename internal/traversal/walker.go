// Package traversal. walker handles walking the directory tree.
package traversal

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/seyedali-dev/treeclip/internal/exclude"
	"github.com/seyedali-dev/treeclip/internal/output"
	"github.com/seyedali-dev/treeclip/pkg/utils"
)

// TraverseDir walks root, writes each file via formatter, returns counts.
func TraverseDir(root string, folderPatterns []string, outputFile io.Writer) (processed, skipped int, err error) {
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		rel, _ := filepath.Rel(root, path)

		if exclude.ShouldExclude(rel, d.Name(), d.IsDir(), folderPatterns) {
			skipped++
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}

		processed++
		output.WriteHeader(outputFile, rel)

		openedFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("❌🪲  [ERROR] error opening file %v: %v", path, err)
		}
		defer utils.SafeCloseFile(openedFile)

		_, _ = io.Copy(outputFile, openedFile)
		output.WriteSeparator(outputFile)
		return nil
	})
	return
}

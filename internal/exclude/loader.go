// Package exclude. loader loads .treeclipignore file and handles it's data.
package exclude

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadIgnorePatterns reads .treeclipignore from the given root path and returns a slice of patterns.
func LoadIgnorePatterns(rootPath string) ([]string, error) {
	var patterns []string

	ignoreFilePath := filepath.Join(rootPath, ".treeclipignore")

	content, err := os.ReadFile(ignoreFilePath)
	if err != nil {
		// File does not exist — not an error
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read .treeclipignore: %w (ノಠ益ಠ)ノ", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Normalize slashes for cross-platform support
		line = filepath.ToSlash(line)
		patterns = append(patterns, line)
	}

	return patterns, nil
}

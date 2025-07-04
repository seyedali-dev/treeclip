// Package exclude. matcher implements matching logic for ignore patterns.
package exclude

import (
	"path/filepath"
	"strings"
)

var DefaultExclusions = []string{
	"treeclip_output.txt",
	"*.tmp", "*.temp", "*.exe", "*.sh",
	".git", ".idea", ".DS_Store", "Thumbs.db",
}

// ShouldExclude checks if a file or directory should be excluded based on the exclude patterns. It supports exact, wildcard, relative file/folder name/path(s).
func ShouldExclude(relPath, name string, isDir bool, patterns []string) bool {
	// Normalize the relative path to use forward slashes
	normalizedRelPath := filepath.ToSlash(relPath)

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		// Normalize the pattern to use forward slashes
		normalizedPattern := filepath.ToSlash(pattern)

		// Check exact name match
		if name == pattern || name == filepath.Base(normalizedPattern) {
			return true
		}

		// Check exact relative path match (both normalized)
		if normalizedRelPath == normalizedPattern {
			return true
		}

		// Check wildcard pattern match against filename
		if matched, _ := filepath.Match(normalizedPattern, name); matched {
			return true
		}

		// Check wildcard pattern match against relative path
		if matched, _ := filepath.Match(normalizedPattern, normalizedRelPath); matched {
			return true
		}

		// For directories, check parent directories
		if isDir {
			pathParts := strings.Split(normalizedRelPath, "/")
			for _, part := range pathParts {
				if matched, _ := filepath.Match(normalizedPattern, part); matched {
					return true
				}
			}
		}
	}
	return false
}

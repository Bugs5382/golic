package pkg

/*
[... Keep your Apache License Header ...]
*/

import (
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// IsMatch checks if a file path matches a given glob pattern (including **).
func IsMatch(path string, pattern string) bool {
	// Keep your existing logic for basic file extensions
	if strings.HasPrefix(pattern, ".") {
		return filepath.Ext(path) == pattern
	}

	cleanPath := filepath.ToSlash(path)
	cleanPattern := filepath.ToSlash(pattern)

	matched, err := doublestar.Match(cleanPattern, cleanPath)
	if err != nil {
		return false
	}

	return matched
}

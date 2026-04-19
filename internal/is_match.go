package internal

/*
Apache License 2.0

Copyright 2026 Shane & Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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

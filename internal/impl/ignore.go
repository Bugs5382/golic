package impl

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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/denormal/go-gitignore"
	"github.com/rs/zerolog/log"
)

// Scope decides which walked paths golic may touch (#46). It layers the
// .licignore file over the ignore: list from the config (the built-in
// defaults plus any local entries):
//
//  1. A path the .licignore ignores is out of scope.
//  2. A path a config ignore pattern matches is out of scope, unless a
//     specific negation line in .licignore matches it too.
//  3. Everything else is in scope.
//
// A negation is specific when it is more than a whitelist catch-all such as
// "!*/" or "!*.go" (see isCatchAll). Whitelist-style .licignore files lean on
// those lines, and letting them lift the defaults would pull every vendored
// and generated file back in.
type scope struct {
	licignore gitignore.GitIgnore
	patterns  []string
	lifts     []lift
}

// lift is one specific negation line from .licignore, compiled on its own so
// it can be tested against a path without the lines around it.
type lift struct {
	line    string
	matcher gitignore.GitIgnore
}

// newScope reads the .licignore at path and pairs it with the config ignores.
func newScope(path string, patterns []string) (*scope, error) {
	log.Debug().Str("licignore", path).Int("ignorePatterns", len(patterns)).Msg("building the file scope")
	gi, err := gitignore.NewFromFile(path)
	if err != nil {
		log.Debug().Err(err).Str("licignore", path).Msg("reading .licignore failed")
		return nil, err
	}
	raw, err := os.ReadFile(path) // #nosec G304 -- the .licignore the user pointed golic at
	if err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	base := filepath.Dir(abs)

	s := &scope{licignore: gi, patterns: patterns}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimRight(line, " \t\r")
		if !strings.HasPrefix(line, "!") {
			continue
		}
		pattern := strings.TrimPrefix(line, "!")
		if isCatchAll(pattern, patterns) {
			log.Trace().Str("line", line).Msg(".licignore negation is a catch-all; it does not lift built-in ignores")
			continue
		}
		log.Debug().Str("line", line).Msg(".licignore negation lifts built-in ignores for the paths it matches")
		s.lifts = append(s.lifts, lift{line: line, matcher: gitignore.New(strings.NewReader(pattern), base, nil)})
	}
	return s, nil
}

// skipDir reports whether the walk can skip a whole directory. .git is never
// walked, at any depth and whatever .licignore says. Other directories are
// skipped only when a config ignore covers everything under them and no
// specific .licignore negation exists that could lift part of it.
func (s *scope) skipDir(path string) (string, bool) {
	if filepath.Base(path) == ".git" {
		return ".git", true
	}
	if len(s.lifts) > 0 {
		return "", false
	}
	return matchIgnore(s.patterns, filepath.ToSlash(path)+"/golic-probe")
}

// excluded reports whether a file is out of scope and, when a config ignore
// is the reason, which pattern matched.
func (s *scope) excluded(path string) (reason string, out bool) {
	if m := s.licignore.Match(path); m != nil && m.Ignore() {
		return ".licignore " + m.String(), true
	}
	pattern, ok := matchIgnore(s.patterns, filepath.ToSlash(path))
	if !ok {
		return "", false
	}
	for _, l := range s.lifts {
		if l.matcher.Match(path) != nil {
			log.Debug().Str("path", path).Str("ignore", pattern).Str("licignore", l.line).Msg("built-in ignore lifted by .licignore")
			return "", false
		}
	}
	return pattern, true
}

// matchIgnore returns the first pattern that matches path (slash separated,
// relative to the walk root).
func matchIgnore(patterns []string, path string) (string, bool) {
	for _, p := range patterns {
		if ok, err := doublestar.Match(p, path); err == nil && ok {
			return p, true
		}
	}
	return "", false
}

// singleExt matches a bare one-extension glob such as "*.go".
var singleExt = regexp.MustCompile(`^\*\.[^*?\[/.]+$`)

// isCatchAll reports whether a .licignore negation (without its "!") is a
// whitelist catch-all: "*", "**", "*/", or a bare one-extension glob such as
// "*.go", with or without a leading "/" or "**/". Repeating a config ignore
// pattern word for word ("!**/*.lock") is never a catch-all, so every default
// can be lifted on purpose.
func isCatchAll(pattern string, patterns []string) bool {
	for _, p := range patterns {
		if pattern == p {
			return false
		}
	}
	p := strings.TrimPrefix(pattern, "/")
	p = strings.TrimPrefix(p, "**/")
	p = strings.TrimSuffix(p, "/")
	return p == "*" || p == "**" || singleExt.MatchString(p)
}

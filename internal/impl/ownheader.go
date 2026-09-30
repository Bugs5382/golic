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
	"regexp"
	"sort"
	"strings"
	"sync"

	golog "github.com/Bugs5382/go-log"
	"github.com/Bugs5382/golic/internal"
	"github.com/Bugs5382/golic/internal/logging"
)

// copyrightMark stands in for {{copyright}} while a template is rendered into
// a search pattern. It cannot occur in a license template or a rule.
const copyrightMark = "\x00golic-copyright\x00"

// ownHeader is a header golic rendered, located in a file.
type ownHeader struct {
	ok         bool
	template   string
	text       string
	start, end int
}

var headerPatterns sync.Map

// headerPattern compiles the rendering of a template for a rule into a regexp
// that matches the exact text, with the copyright free to be anything on its
// line.
func headerPattern(rendered string) (*regexp.Regexp, error) {
	if re, ok := headerPatterns.Load(rendered); ok {
		return re.(*regexp.Regexp), nil
	}
	parts := strings.Split(rendered, copyrightMark)
	for i, p := range parts {
		parts[i] = regexp.QuoteMeta(p)
	}
	re, err := regexp.Compile(strings.Join(parts, `[^\n]*`))
	if err != nil {
		return nil, err
	}
	headerPatterns.Store(rendered, re)
	return re, nil
}

// findOwnHeader looks for a header golic could have written into path: any
// configured license rendered for the file's rule, with any copyright. The
// earliest match wins, and the longest one at that offset, so a template whose
// text sits inside another's is never picked over the one that wrote the file.
func findOwnHeader(source, path string, config *Config) (ownHeader, error) {
	keys := make([]string, 0, len(config.Golic.Licenses))
	for k := range config.Golic.Licenses {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	best := ownHeader{}
	for _, k := range keys {
		rendered, err := getCommentedLicense(config, internal.Options{Template: k, Copyright: copyrightMark}, path)
		if err != nil {
			return ownHeader{}, err
		}
		re, err := headerPattern(rendered)
		if err != nil {
			return ownHeader{}, err
		}
		loc := re.FindStringIndex(source)
		if loc == nil {
			continue
		}
		logging.L().Trace("golic header found", golog.F("path", path), golog.F("template", k), golog.F("offset", loc[0]))
		if !best.ok || loc[0] < best.start || (loc[0] == best.start && loc[1] > best.end) {
			best = ownHeader{ok: true, template: k, text: source[loc[0]:loc[1]], start: loc[0], end: loc[1]}
		}
	}
	return best, nil
}

// licenseWords mark a leading comment as a license header.
var licenseWords = []string{"copyright", "license", "spdx-license-identifier", "all rights reserved"}

// foreignHeader reports whether text opens with a comment in the rule's style
// that reads like a license header. replace leaves such a file alone rather
// than guessing where someone else's header ends.
func foreignHeader(text string, r Rule) bool {
	open := strings.TrimSpace(strings.TrimLeft(r.Prefix, "\n"))
	if open == "" {
		return false
	}
	rest := strings.TrimLeft(text, " \t\r\n")
	if !strings.HasPrefix(rest, open) {
		return false
	}

	var block string
	if r.Suffix != "" {
		end := strings.Index(rest, strings.TrimSpace(r.Suffix))
		if end == -1 {
			return false
		}
		block = rest[:end]
	} else {
		var lines []string
		for _, l := range strings.Split(rest, "\n") {
			if !strings.HasPrefix(strings.TrimSpace(l), open) {
				break
			}
			lines = append(lines, l)
		}
		block = strings.Join(lines, "\n")
	}

	lower := strings.ToLower(block)
	for _, w := range licenseWords {
		if strings.Contains(lower, w) {
			return true
		}
	}
	return false
}

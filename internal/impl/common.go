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
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	golog "github.com/Bugs5382/go-log"
	"github.com/Bugs5382/golic/internal/logging"

	"github.com/Bugs5382/golic/internal"
	"github.com/enescakir/emoji"
	"github.com/goccy/go-yaml"
	"github.com/logrusorgru/aurora"
)

func injectFile(path string, o internal.Options, config *Config) (rule string, skip bool, err error) {
	source, err := read(path)
	if err != nil {
		return "", false, err
	}
	rule = getRule(config, path)
	license, err := getCommentedLicense(config, o, rule)
	if err != nil {
		return "", false, err
	}
	// license is injected, continue
	if strings.Contains(source, license) {
		return rule, true, nil
	}
	source = insertLicense(source, license, config.Golic.Rules[rule])

	if !o.Dry {
		data := []byte(source)
		err = os.WriteFile(path, data, os.ModeExclusive)
	}
	return
}

// insertLicense puts license below the lines the rule keeps first.
func insertLicense(source, license string, r Rule) string {
	header, footer := splitSource(source, r.Under)
	if header != "" {
		header = header + "\n"
	}
	return fmt.Sprintf("%s%s%s", header, license, footer)
}

// removeFile Remove text from file overall function
func removeFile(path string, o internal.Options, config *Config) (rule string, skip bool, err error) {
	source, err := read(path)
	if err != nil {
		return "", false, err
	}
	rule = getRule(config, path)
	license, err := getCommentedLicense(config, o, rule)
	if err != nil {
		return rule, false, err
	}
	if strings.Contains(source, license) {
		return "", false, RemoveFromFile(path, o, source, license, err)
	}
	return rule, true, nil
}

// replaceFile swaps a header golic wrote for the configured one. It looks for
// the exact rendering of any configured license for the file's rule, with any
// copyright value, and replaces that text in place, so the bytes around the
// header (comments, blank lines, a shebang) are kept. A file with no header is
// stamped the way inject would stamp it. A file whose leading comment is some
// other license header is left alone and reported.
func replaceFile(path string, o internal.Options, config *Config) (rule string, skip bool, err error) {
	source, err := read(path)
	if err != nil {
		return "", false, err
	}
	rule = getRule(config, path)
	license, err := getCommentedLicense(config, o, rule)
	if err != nil {
		return rule, false, err
	}

	found, err := findOwnHeader(source, path, config)
	if err != nil {
		return rule, false, err
	}

	var updated string
	switch {
	case found.ok && found.text == license:
		logging.L().Trace("already carries the target header", golog.F("path", path), golog.F("template", o.Template))
		return rule, true, nil
	case found.ok:
		logging.L().Debug("swapping golic header in place", golog.F("path", path), golog.F("from", found.template), golog.F("to", o.Template), golog.F("offset", found.start))
		updated = source[:found.start] + license + source[found.end:]
	case strings.Contains(source, license):
		return rule, true, nil
	default:
		r := config.Golic.Rules[rule]
		_, footer := splitSource(source, r.Under)
		if foreignHeader(footer, r) {
			logging.L().Warn("file has a license header golic did not write; left unchanged", golog.F("path", path))
			return rule, true, nil
		}
		logging.L().Debug("no header found; injecting", golog.F("path", path))
		updated = insertLicense(source, license, r)
	}

	if updated == source {
		return rule, true, nil
	}

	if !o.Dry {
		err = os.WriteFile(path, []byte(updated), os.ModeExclusive)
	}
	return rule, false, err
}

// readCommonConfig Read the commong/master config
func (u *Process) readCommonConfig() (*Config, error) {
	c := &Config{}
	rawYaml := defaultConfig

	if err := yaml.Unmarshal([]byte(rawYaml), c); err != nil {
		return nil, fmt.Errorf("failed to parse master config: %w", err)
	}

	return c, nil
}

// readLocalConfig Read the local config.
func (u *Process) readLocalConfig() (*Config, error) {
	var rc = &Config{}

	if rc.Golic.Licenses == nil {
		rc.Golic.Licenses = make(map[string]string)
	}

	if u.cfgBase.Golic.Licenses != nil {
		for k, v := range u.cfgBase.Golic.Licenses {
			rc.Golic.Licenses[k] = v
		}
	}

	if rc.Golic.Rules == nil {
		rc.Golic.Rules = make(map[string]Rule)
	}

	if u.cfgBase.Golic.Rules != nil {
		for k, v := range u.cfgBase.Golic.Rules {
			rc.Golic.Rules[k] = v
		}
	}

	rc.Golic.Ignore = append([]string(nil), u.cfgBase.Golic.Ignore...)

	// If the path is empty or file doesn't exist, we return the base copy immediately
	if u.Opts.ConfigPath == "" {
		return rc, nil
	}

	yamlFile, err := os.ReadFile(u.Opts.ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return rc, nil
		}
		return nil, fmt.Errorf("failed to read local config at %s: %w", u.Opts.ConfigPath, err)
	}

	var localCfg = &Config{}
	localCfg.Golic.MergeRules = true // Set default expectation

	if err := yaml.Unmarshal(yamlFile, localCfg); err != nil {
		return nil, fmt.Errorf("failed to parse local config: %w", err)
	}

	// If localCfg has licenses, they overwrite/append to the base in rc
	for k, v := range localCfg.Golic.Licenses {
		rc.Golic.Licenses[k] = v
	}

	// Local ignore patterns add to the built-in list; they never replace it.
	if len(localCfg.Golic.Ignore) > 0 {
		logging.L().Debug("adding local ignore patterns to the built-in list", golog.F("ignore", localCfg.Golic.Ignore))
		rc.Golic.Ignore = append(rc.Golic.Ignore, localCfg.Golic.Ignore...)
	}

	if localCfg.Golic.MergeRules {
		// Append or overwrite individual rules
		for k, v := range localCfg.Golic.Rules {
			rc.Golic.Rules[k] = v
		}
	} else if len(localCfg.Golic.Rules) > 0 {
		// Replace the entire ruleset if MergeRules is explicitly false
		rc.Golic.Rules = localCfg.Golic.Rules
	}

	return rc, nil
}

// traverseFiles walks the tree from ./ and processes every file in scope:
// not ignored by .licignore or the config ignore list, covered by a rule, and
// not generated.
func (u *Process) traverseFiles() error {
	skipped := 0
	visited := 0
	outOfScope := 0
	generated := 0
	p := func(path string, o internal.Options, config *Config) (err error) {
		if reason, out := u.scope.excluded(path); out {
			logging.L().Trace("out of scope", golog.F("path", path), golog.F("reason", reason))
			outOfScope++
			return nil
		}
		ruleName := getRule(config, path)
		if ruleName == "" {
			logging.L().Trace("no rule matches; skipping", golog.F("path", path))
			return nil
		}
		marker, gen, err := isGenerated(path)
		if err != nil {
			return err
		}
		if gen {
			logging.L().Debug("generated file; skipping", golog.F("path", path), golog.F("marker", marker))
			generated++
			return nil
		}

		var rule string
		var skip bool
		symbol := ""
		prefix := ""
		cp := aurora.BrightYellow(path)

		visited++

		if rule, skip, err = processUpdate(path, o, config); err != nil {
			return err
		} else if skip {
			symbol = "-> skip"
			cp = aurora.Magenta(path)
			skipped++
		}

		if u.Opts.Dry {
			prefix = aurora.Bold(aurora.Yellow(fmt.Sprintf("%s DRY RUN: ", emoji.TestTube))).String()
		}

		if logging.DebugEnabled() {
			logging.L().Info(fmt.Sprintf("%s %s  %s %s %s",
				prefix,
				emoji.Minus,
				cp,
				aurora.Bold(aurora.BrightMagenta(symbol)),
				aurora.Gray(12, fmt.Sprintf("[%s]", rule)),
			))
		} else {
			logging.L().Info(fmt.Sprintf("%s %s  %s %s",
				prefix,
				emoji.Minus,
				cp,
				aurora.Bold(aurora.BrightMagenta(symbol)),
			))
		}
		return nil
	}

	err := filepath.Walk("./",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if path == "." {
					return nil
				}
				if pattern, skip := u.scope.skipDir(path); skip {
					logging.L().Debug("not walking an ignored directory", golog.F("dir", path), golog.F("ignore", pattern))
					return filepath.SkipDir
				}
				return nil
			}
			return p(path, u.Opts, u.cfg)
		})

	if err != nil {
		return err
	}

	u.modified = visited - skipped
	logging.L().Debug("walk finished", golog.F("visited", visited), golog.F("unchanged", skipped), golog.F("outOfScope", outOfScope), golog.F("generated", generated))
	displaySummary(skipped, visited)

	return nil
}

// processUpdate Update the file, but how?
func processUpdate(path string, o internal.Options, config *Config) (rule string, skip bool, err error) {
	switch o.Type {
	case internal.LicenseInject:
		return injectFile(path, o, config)
	case internal.LicenseRemove:
		return removeFile(path, o, config)
	case internal.LicenseReplace:
		return replaceFile(path, o, config)
	}
	return "", true, fmt.Errorf("invalid license type")
}

func displaySummary(skipped, visited int) {
	if skipped == visited {
		logging.L().Info(fmt.Sprintf("%s %v/%v %s", emoji.Ice, aurora.BrightCyan(visited-skipped), aurora.BrightWhite(visited), aurora.BrightCyan("changed")))
		return
	}
	logging.L().Info(fmt.Sprintf("%s %v/%v %s", emoji.Fire, aurora.BrightYellow(visited-skipped), aurora.BrightWhite(visited), aurora.BrightYellow("changed")))
}

// read File
func read(f string) (s string, err error) {
	logging.L().Trace("reading source file", golog.F("path", f))
	// The path comes from golic's own walk of the working tree, filtered by
	// .licignore. Reading the files the user asked golic to process is the
	// point of the tool, so this is not untrusted input (#34).
	content, err := os.ReadFile(f) // #nosec G304 -- path comes from the working-tree walk
	if err != nil {
		logging.L().Debug("reading source file failed", golog.F("path", f), golog.F("error", err.Error()))
		return
	}
	// Convert []byte to string and print to screen
	return string(content), nil
}

// RemoveFromFile
func RemoveFromFile(path string, o internal.Options, source string, license string, err error) error {
	if !o.Dry {
		source = strings.Replace(source, license, "", 1)
		err = os.WriteFile(path, []byte(source), os.ModeExclusive)
	}
	return err
}

// matchRule Match Rule
func matchRule(config *Config, path string) (rule string, ok bool) {
	rule = getRule(config, path)
	_, ok = config.Golic.Rules[rule]
	return rule, ok
}

// getCommentedLicense Get Commented License File
func getCommentedLicense(config *Config, o internal.Options, file string) (string, error) {
	var ok bool
	var template string
	var rule string
	if template, ok = config.Golic.Licenses[o.Template]; !ok {
		return "", fmt.Errorf("no license found for %s, check configuration (.golic.yaml)", o.Template)
	}

	//if _, ok =  config.Golic.Rules[rule]; !ok {
	if rule, ok = matchRule(config, file); !ok {
		return "", fmt.Errorf("no rule found for %s, check configuration (.golic.yaml)", rule)
	}
	template = strings.ReplaceAll(template, "{{copyright}}", o.Copyright)
	if config.IsWrapped(rule) {
		return fmt.Sprintf("%s\n%s%s\n",
				config.Golic.Rules[rule].Prefix,
				template,
				config.Golic.Rules[rule].Suffix),
			nil
	}
	// `\r\n` -> `\r\n #`, `\n` -> `\n #`
	content := strings.ReplaceAll(template, "\n", fmt.Sprintf("\n%s", config.Golic.Rules[rule].Prefix))
	content = strings.TrimSuffix(content, config.Golic.Rules[rule].Prefix)
	content = config.Golic.Rules[rule].Prefix + content
	// "# \n" -> "#\n" // "# \r\n" -> "#\r\n"; some environments automatically remove spaces in empty lines. This makes problems in license PR's
	cleanedPrefix := strings.TrimSuffix(config.Golic.Rules[rule].Prefix, " ")
	content = strings.ReplaceAll(content, fmt.Sprintf("%s \n", cleanedPrefix), fmt.Sprintf("%s\n", cleanedPrefix))
	content = strings.ReplaceAll(content, fmt.Sprintf("%s \r\n", cleanedPrefix), fmt.Sprintf("%s\r\n", cleanedPrefix))
	return content, nil
}

// splitSource splits source into the lines that must stay above the license
// (the header) and the rest (the footer), using a rule's "under" patterns.
// Patterns are tried in order; the first one that matches a line ends the
// header at that line. Lines right after it that match any of the patterns
// join the header too, so a run of directives (Dockerfile "# syntax=" and
// "# escape=", in either order) stays together above the license.
func splitSource(source string, rules []string) (header, footer string) {
	lines := strings.Split(source, "\n")
	if len(rules) == 0 {
		return "", source
	}
	for _, r := range rules {
		header, footer = findHeaderAndFooter(lines, r, rules)
		if header != "" {
			logging.L().Trace("license goes below a matched line", golog.F("under", r), golog.F("headerLines", strings.Count(header, "\n")+1))
			return
		}
	}
	logging.L().Trace("no under pattern matched; license goes at the top", golog.F("under", rules))
	return
}

func findHeaderAndFooter(lines []string, match string, all []string) (header, footer string) {
	for i, l := range lines {
		if underMatch(i, l, match) {
			end := i + 1
			for end < len(lines) && anyUnderMatch(end, lines[end], all) {
				end++
			}
			header = strings.Join(lines[0:end], "\n")
			footer = strings.Join(lines[end:], "\n")
			return
		}
	}
	return "", strings.Join(lines, "\n")
}

// underMatch reports whether line number i matches an "under" pattern. A
// shebang only means something on the first line, so patterns starting with
// "#!" never match further down (a heredoc or template literal that happens
// to hold one).
func underMatch(i int, line, pattern string) bool {
	if i > 0 && strings.HasPrefix(pattern, "#!") {
		return false
	}
	return internal.IsMatch(line, pattern)
}

func anyUnderMatch(i int, line string, patterns []string) bool {
	for _, p := range patterns {
		if underMatch(i, line, p) {
			return true
		}
	}
	return false
}

// getRule Get Rule for Match
func getRule(config *Config, path string) (rule string) {
	keys := make([]string, 0, len(config.Golic.Rules))
	for k := range config.Golic.Rules {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})

	cleanPath := filepath.ToSlash(path)
	for _, k := range keys {
		if internal.IsMatch(cleanPath, k) {
			return k
		}
	}

	ext := filepath.Ext(path)
	if _, ok := config.Golic.Rules[ext]; ok {
		return ext
	}

	return ""
}

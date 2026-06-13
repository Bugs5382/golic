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

	"github.com/Bugs5382/golic"
	"github.com/Bugs5382/golic/internal"
	"github.com/denormal/go-gitignore"
	"github.com/enescakir/emoji"
	"github.com/goccy/go-yaml"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog/log"
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
	// split file to header and footer and extend with license
	header, footer := splitSource(source, config.Golic.Rules[rule].Under)
	if header != "" {
		header = header + "\n"
	}
	source = fmt.Sprintf("%s%s%s", header, license, footer)

	if !o.Dry {
		data := []byte(source)
		err = os.WriteFile(path, data, os.ModeExclusive)
	}
	return
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

// replaceFile Remove any existing license header and inject the current one in a
// single pass. Unlike inject, replace strips whatever license-style comment block
// already sits in the header region so the copyright string or license type can be
// changed, not just refreshed in place. When the file already carries the exact
// license, it is left untouched and reported as skipped so a no-op replace does not
// count as a change.
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

	// Already stamped with the current license; nothing to replace.
	if strings.Contains(source, license) {
		return rule, true, nil
	}

	r := config.Golic.Rules[rule]

	// Split the file into the part that precedes the license region and the rest.
	header, footer := splitSource(source, r.Under)
	if header != "" {
		header = header + "\n"
	}

	// Drop an existing license block from the front of the license region before
	// injecting the new one. Normalize the leading newlines of what remains so the
	// spacing matches a fresh inject regardless of how the old block was laid out.
	footer = stripLicenseBlock(footer, r)
	if header != "" && footer != "" {
		// The license sits below an "under" line (e.g. package). Inject leaves a
		// single newline between the license and the body; mirror that here.
		footer = "\n" + strings.TrimLeft(footer, "\n")
	} else if header == "" {
		// The license sits at the very top of the file; drop any blank lines the
		// removed block left behind so the body starts cleanly after the license.
		footer = strings.TrimLeft(footer, "\n")
	}

	updated := fmt.Sprintf("%s%s%s", header, license, footer)

	// Nothing actually changed (e.g. an unstamped file already matches output).
	if updated == source {
		return rule, true, nil
	}

	if !o.Dry {
		err = os.WriteFile(path, []byte(updated), os.ModeExclusive)
	}
	return rule, false, err
}

// stripLicenseBlock removes a leading license-style comment block from text using
// the comment delimiters of the given rule. Wrapped rules (with a suffix) are
// matched from their opening delimiter to the first closing delimiter; line-comment
// rules drop the contiguous run of prefixed lines. Leading blank lines are tolerated
// so the detection survives the blank line inject leaves between header and license.
func stripLicenseBlock(text string, r Rule) string {
	prefix := strings.TrimLeft(r.Prefix, "\n")

	// Preserve any leading blank lines, then inspect the first non-blank line.
	rest := text
	lead := ""
	for {
		nl := strings.IndexByte(rest, '\n')
		if nl == -1 {
			break
		}
		if strings.TrimSpace(rest[:nl]) != "" {
			break
		}
		lead += rest[:nl+1]
		rest = rest[nl+1:]
	}

	trimmedPrefix := strings.TrimSpace(prefix)
	if trimmedPrefix == "" || !strings.HasPrefix(strings.TrimSpace(rest), trimmedPrefix) {
		// No recognizable license block at the front; leave the text untouched.
		return text
	}

	if r.Suffix != "" {
		// Wrapped block: cut from the opener through the first closing delimiter.
		suffix := strings.TrimSpace(r.Suffix)
		idx := strings.Index(rest, suffix)
		if idx == -1 {
			return text
		}
		after := rest[idx+len(suffix):]
		after = strings.TrimPrefix(after, "\n")
		return lead + after
	}

	// Line-comment block: drop the contiguous run of prefixed lines.
	lines := strings.Split(rest, "\n")
	cut := 0
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), trimmedPrefix) {
			cut++
			continue
		}
		break
	}
	remaining := strings.Join(lines[cut:], "\n")
	remaining = strings.TrimPrefix(remaining, "\n")
	return lead + remaining
}

// readCommonConfig Read the commong/master config
func (u *Process) readCommonConfig() (*Config, error) {
	c := &Config{}
	rawYaml := golic.DefaultConfig

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

// traverseFiles Go through all files in paths and process. Will ignore files and folders that match GitIgnore.
func (u *Process) traverseFiles() error {
	skipped := 0
	visited := 0
	p := func(path string, i gitignore.GitIgnore, o internal.Options, config *Config) (err error) {
		if !i.Ignore(path) {
			ruleName := getRule(config, path)
			if ruleName == "" {
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

			if log.Debug().Enabled() {
				log.Info().Msgf("%s %s  %s %s %s",
					prefix,
					emoji.Minus,
					cp,
					aurora.Bold(aurora.BrightMagenta(symbol)),
					aurora.Gray(12, fmt.Sprintf("[%s]", rule)),
				)
			} else {
				log.Info().Msgf("%s %s  %s %s",
					prefix,
					emoji.Minus,
					cp,
					aurora.Bold(aurora.BrightMagenta(symbol)),
				)
			}
		}
		return nil
	}

	err := filepath.Walk("./",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				return p(path, u.ignore, u.Opts, u.cfg)
			}
			return nil
		})

	if err != nil {
		return err
	}

	u.modified = visited - skipped
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
		log.Info().Msgf("%s %v/%v %s", emoji.Ice, aurora.BrightCyan(visited-skipped), aurora.BrightWhite(visited), aurora.BrightCyan("changed"))
		return
	}
	log.Info().Msgf("%s %v/%v %s", emoji.Fire, aurora.BrightYellow(visited-skipped), aurora.BrightWhite(visited), aurora.BrightYellow("changed"))
}

// read File
func read(f string) (s string, err error) {
	content, err := os.ReadFile(f)
	if err != nil {
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

// splitSource Split Source
func splitSource(source string, rules []string) (header, footer string) {
	lines := strings.Split(source, "\n")
	if len(rules) == 0 {
		return "", source
	}
	for _, r := range rules {
		header, footer = findHeaderAndFooter(lines, r)
		if header != "" {
			return
		}
	}
	return
}

func findHeaderAndFooter(lines []string, match string) (header, footer string) {
	for i, l := range lines {
		if internal.IsMatch(l, match) {
			header = strings.Join(lines[0:i+1], "\n")
			footer = strings.Join(lines[i+1:], "\n")
			return
		}
	}
	return "", strings.Join(lines, "\n")
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

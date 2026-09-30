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

import "fmt"

// Rule says how a license header is written for the files a rule key matches.
type Rule struct {
	// Prefix starts every header line, or opens the block when Suffix is set.
	// golic finds its own header by the exact rendered text, so changing the
	// prefix or suffix of a rule already in use stamps every file again.
	Prefix string `yaml:"prefix"`
	// Suffix closes a block comment. Empty means a line-comment rule.
	Suffix string `yaml:"suffix"`
	// Under lists line patterns the header must go below (a shebang, an XML
	// declaration, a Go package clause). The first pattern that matches a
	// line wins, and the lines right after it that match any pattern stay
	// above the header too. Patterns starting with "#!" only match line 1.
	Under []string `yaml:"under"`
}

type GolicConfig struct {
	Licenses map[string]string `yaml:"licenses"`
	// MergeRules controls how the rules in a local .golic.yaml combine with
	// the built-in ones. It defaults to true: each local rule is added, or
	// replaces the built-in rule with the same key as a whole (fields are not
	// merged). When false and the local file lists at least one rule, the
	// local rules replace the built-in set entirely, and any file type not
	// listed there has no rule and is skipped without a warning. When false
	// with no local rules, the built-in rules are kept. Licenses always merge.
	MergeRules bool            `yaml:"mergeRules"`
	Rules      map[string]Rule `yaml:"rules"`
	// Ignore lists doublestar globs, relative to the directory golic runs
	// in, for files that are never processed (#46). The built-in list covers
	// VCS metadata, vendored and built trees, lock files and generated code.
	// A local .golic.yaml adds to it. A specific "!" line in .licignore can
	// lift an entry for the paths it names.
	Ignore []string `yaml:"ignore"`
	// LicenseFile turns on writing the full licence text: true for LICENSE,
	// or a path. The --license-file flag wins over it.
	LicenseFile LicenseFileSetting `yaml:"licenseFile"`
}

// LicenseFileSetting is the licenseFile config value, a bool or a path.
type LicenseFileSetting struct {
	Path string
}

// UnmarshalYAML accepts true (LICENSE), false (off) or a path.
func (l *LicenseFileSetting) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v interface{}
	if err := unmarshal(&v); err != nil {
		return err
	}
	switch x := v.(type) {
	case nil:
		l.Path = ""
	case bool:
		l.Path = ""
		if x {
			l.Path = defaultLicenseFile
		}
	case string:
		l.Path = x
	default:
		return fmt.Errorf("licenseFile must be true, false or a path, got %v", v)
	}
	return nil
}

// MarshalYAML writes the setting back as the path, for the debug config dump.
func (l LicenseFileSetting) MarshalYAML() (interface{}, error) {
	return l.Path, nil
}

type Config struct {
	Golic GolicConfig `yaml:"golic"`
}

func (c *Config) IsWrapped(key string) bool {
	return c.Golic.Rules[key].Suffix != ""
}

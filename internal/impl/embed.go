package impl

/*
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
	_ "embed"

	"github.com/goccy/go-yaml"
)

// defaultConfig is the built-in ruleset and license catalogue. A local
// .golic.yaml is merged on top of it at run time.
//
//go:embed default.golic.yaml
var defaultConfig string

// supersededConfig holds the earlier texts of the built-in licenses.
//
//go:embed superseded.golic.yaml
var supersededConfig string

// supersededLicenses maps a built-in license key to its earlier texts, newest
// first. Headers rendered from them are still golic's own.
var supersededLicenses = mustParseSuperseded(supersededConfig)

func mustParseSuperseded(raw string) map[string][]string {
	var doc struct {
		Superseded map[string][]string `yaml:"superseded"`
	}
	if err := yaml.Unmarshal([]byte(raw), &doc); err != nil {
		panic("golic: embedded superseded.golic.yaml does not parse: " + err.Error())
	}
	return doc.Superseded
}

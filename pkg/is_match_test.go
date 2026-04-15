package pkg_test

/*
Apache License 2.0

Copyright 2026 Shane

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
	"testing"

	"github.com/Bugs5382/golic/impl"
	"github.com/Bugs5382/golic/pkg"
)

func TestMatchRule(t *testing.T) {
	config := &impl.Config{
		Golic: impl.GolicConfig{
			Licenses:   make(map[string]string),
			MergeRules: false,
			Rules: map[string]impl.Rule{
				"Makefile": {
					Prefix: "#",
				},
				"**/templates/**/*.yaml": {
					Prefix: "{{/*",
					Suffix: "*/}}",
				},
				"*.go": {
					Prefix: "/*",
					Suffix: "*/",
				},
				"cmd/server/main.go": {
					Prefix: "//",
				},
			},
		},
	}

	tests := []struct {
		name     string
		fileName string
		wantRule string
		wantOk   bool
	}{
		{name: "Direct Match (Exact)", fileName: "Makefile", wantRule: "Makefile", wantOk: true},
		{name: "Recursive Match (Deep Path)", fileName: "charts/technitium/templates/service.yaml", wantRule: "**/templates/**/*.yaml", wantOk: true},
		{name: "Wildcard Match (Shallow)", fileName: "main.go", wantRule: "*.go", wantOk: true},
		{name: "Precedence Match", fileName: "cmd/server/main.go", wantRule: "cmd/server/main.go", wantOk: true},
		{name: "No Match Case", fileName: "README.md", wantRule: "", wantOk: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// FIX: We loop through the rules manually to find which pattern matches
			var gotRule string
			var gotOk bool

			for pattern := range config.Golic.Rules {
				if pkg.IsMatch(tt.fileName, pattern) {
					// Precedence Logic: If we find multiple matches,
					// we usually want the most specific (longest) pattern.
					if len(pattern) > len(gotRule) {
						gotRule = pattern
						gotOk = true
					}
				}
			}

			if gotOk != tt.wantOk {
				t.Fatalf("IsMatch(%s) ok = %v, want %v", tt.fileName, gotOk, tt.wantOk)
			}

			if gotRule != tt.wantRule {
				t.Errorf("IsMatch(%s) rule = %v, want %v", tt.fileName, gotRule, tt.wantRule)
			}
		})
	}
}

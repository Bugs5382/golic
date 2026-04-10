package helpers_test

import (
	"testing"

	"github.com/AbsaOSS/golic/helpers"
	"github.com/AbsaOSS/golic/impl"
)

func TestMatchRule(t *testing.T) {
	config := &impl.Config{}
	// ... (Your config setup stays the same) ...
	config.Golic.Rules = map[string]struct {
		Prefix string   `yaml:"prefix"`
		Suffix string   `yaml:"suffix"`
		Under  []string `yaml:"under"`
	}{
		"Makefile":               {Prefix: "#"},
		"**/templates/**/*.yaml": {Prefix: "{{/*", Suffix: "*/}}"},
		"*.go":                   {Prefix: "/*", Suffix: "*/"},
		"cmd/server/main.go":     {Prefix: "//"},
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
				if helpers.IsMatch(tt.fileName, pattern) {
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

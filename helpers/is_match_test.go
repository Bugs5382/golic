package helpers

import (
	"testing"
)

func TestIsMatch(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		pattern string
		want    bool
	}{
		// Extension matches
		{"ExtMatch", "values.yaml", ".yaml", true},
		{"ExtMismatch", "values.yml", ".yaml", false},

		// Standard wildcards
		{"SingleStarMatch", "service.yaml", "*.yaml", true},
		{"SingleStarDeepMismatch", "templates/service.yaml", "*.yaml", false},

		// Recursive wildcards
		{"RecursiveTemplatesMatch", "charts/templates/service.yaml", "**/templates/**/*.yaml", true},
		{"RecursiveDeepMatch", "charts/sub/templates/nested/pod.yaml", "**/templates/**/*.yaml", true},
		{"RecursiveRootMatch", "templates/pod.yaml", "**/templates/**/*.yaml", true},

		// Question mark wildcards
		{"QuestionMarkMatch", "templates/service.yaml", "template?/service.yaml", true},
		{"QuestionMarkMismatch", "template/service.yaml", "template?/service.yaml", false},
		{"QuestionMarkMulti", "v1.yaml", "v?.yaml", true},

		// General Glob behavior
		{"GlobalYaml", "some/deep/path/config.yaml", "**/*.yaml", true},
		{"GlobalMismatch", "some/deep/path/config.json", "**/*.yaml", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMatch(tt.path, tt.pattern)
			if got != tt.want {
				t.Errorf("isMatch(%q, %q) = %v; want %v", tt.path, tt.pattern, got, tt.want)
			}
		})
	}
}

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
	"strings"
	"testing"

	"github.com/Bugs5382/golic/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// embeddedConfig loads the bundled default ruleset the same way Process.Run
// does, without reading any local .golic.yaml override.
func embeddedConfig(t *testing.T) *Config {
	t.Helper()
	u := &Process{Opts: internal.Options{}}
	cfg, err := u.readCommonConfig()
	require.NoError(t, err)
	return cfg
}

// TestBuiltinTypeScriptRules verifies golic ships block-comment rules for the
// TypeScript family (and the additional JS variants) out of the box, so a
// project only needs a .licignore and no local .golic.yaml override.
func TestBuiltinTypeScriptRules(t *testing.T) {
	cfg := embeddedConfig(t)

	for _, ext := range []string{".ts", ".tsx", ".mts", ".cts", ".js", ".jsx", ".mjs", ".cjs"} {
		t.Run(ext, func(t *testing.T) {
			rule, ok := cfg.Golic.Rules[ext]
			require.True(t, ok, "embedded config should define a rule for %s", ext)
			assert.Equal(t, "/*", rule.Prefix, "%s should use the block-comment prefix", ext)
			assert.Equal(t, "*/", rule.Suffix, "%s should use the block-comment suffix", ext)
			assert.True(t, cfg.IsWrapped(ext), "%s should be treated as a wrapped (block) comment", ext)
		})
	}
}

// TestInjectRemoveTypeScriptRoundTrip proves the embedded TypeScript rules
// drive a full inject -> remove round-trip on a real .ts/.tsx/.mts/.cts file
// without any local config. Fixtures live in t.TempDir(), so they are never
// tracked and never seen by golic's own license gate.
func TestInjectRemoveTypeScriptRoundTrip(t *testing.T) {
	cfg := embeddedConfig(t)

	const body = "export const value = 1;\n"

	for _, ext := range []string{".ts", ".tsx", ".mts", ".cts"} {
		t.Run(ext, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "fixture"+ext)
			require.NoError(t, os.WriteFile(path, []byte(body), 0o600))

			injectOpts := internal.Options{
				Type:      internal.LicenseInject,
				Template:  "mit",
				Copyright: "2026 Shane & Contributors",
			}

			rule, skip, err := injectFile(path, injectOpts, cfg)
			require.NoError(t, err)
			assert.Equal(t, ext, rule, "inject should match the %s rule", ext)
			assert.False(t, skip, "first inject should modify the file")

			injected, err := os.ReadFile(path)
			require.NoError(t, err)
			got := string(injected)
			assert.True(t, strings.HasPrefix(got, "/*\n"), "%s header should open with a block comment, got: %q", ext, got)
			assert.Contains(t, got, "MIT License")
			assert.Contains(t, got, "Copyright (c) 2026 Shane & Contributors")
			assert.Contains(t, got, "*/")
			assert.True(t, strings.HasSuffix(got, body), "original source should be preserved after the header")

			// Re-injecting is a no-op (header already present).
			_, skip, err = injectFile(path, injectOpts, cfg)
			require.NoError(t, err)
			assert.True(t, skip, "second inject should skip an already-licensed %s file", ext)

			removeOpts := internal.Options{
				Type:      internal.LicenseRemove,
				Template:  "mit",
				Copyright: "2026 Shane & Contributors",
			}

			_, skip, err = removeFile(path, removeOpts, cfg)
			require.NoError(t, err)
			assert.False(t, skip, "remove should strip the header from a licensed %s file", ext)

			removed, err := os.ReadFile(path)
			require.NoError(t, err)
			assert.Equal(t, body, string(removed), "remove should restore the original %s source", ext)
		})
	}
}

// TestLocalConfigOverridesBuiltinTypeScriptRule confirms a user-supplied
// .golic.yaml can still override the bundled TypeScript rule (the merge path
// is not broken by shipping the built-in).
func TestLocalConfigOverridesBuiltinTypeScriptRule(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".golic.yaml")
	const local = `golic:
  rules:
    .ts:
      prefix: "//"
`
	require.NoError(t, os.WriteFile(cfgPath, []byte(local), 0o600))

	u := &Process{Opts: internal.Options{ConfigPath: cfgPath}}
	base, err := u.readCommonConfig()
	require.NoError(t, err)
	u.cfgBase = base

	merged, err := u.readLocalConfig()
	require.NoError(t, err)

	assert.Equal(t, "//", merged.Golic.Rules[".ts"].Prefix, "local config should override the built-in .ts prefix")
	assert.Equal(t, "", merged.Golic.Rules[".ts"].Suffix, "override should drop the built-in block suffix")
	// Untouched built-in TS rules remain available.
	assert.Equal(t, "/*", merged.Golic.Rules[".tsx"].Prefix, "unoverridden built-in .tsx rule should survive the merge")
}

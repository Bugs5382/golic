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
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bugs5382/golic/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCopyright = "2026 Shane & Contributors"

// runProcess writes files into a fresh directory, runs a full golic pass over it
// the way the CLI does (walk from ./, embedded rules, an empty .licignore kept
// outside the tree) and returns every file's content afterwards together with
// the number of files the pass changed.
func runProcess(t *testing.T, kind internal.LicenseCommandType, template string, files map[string]string) (map[string]string, int) {
	t.Helper()

	ignore := filepath.Join(t.TempDir(), ".licignore")
	require.NoError(t, os.WriteFile(ignore, nil, 0o600))

	dir := t.TempDir()
	for name, content := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o750))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o600))
	}
	t.Chdir(dir)

	u := ProcessFile(context.Background(), internal.Options{
		Type:      kind,
		Template:  template,
		Copyright: testCopyright,
		LicIgnore: ignore,
	})
	require.NoError(t, u.Run())

	out := make(map[string]string, len(files))
	for name := range files {
		b, err := os.ReadFile(filepath.FromSlash(name))
		require.NoError(t, err)
		out[name] = string(b)
	}
	return out, u.modified
}

// header renders the license block golic writes for path with the embedded rules.
func header(t *testing.T, path, template string) string {
	t.Helper()
	h, err := getCommentedLicense(embeddedConfig(t), internal.Options{Template: template, Copyright: testCopyright}, path)
	require.NoError(t, err)
	return h
}

// TestInjectKeepsShebangFirst checks the shebang stays on line 1 for every
// script-capable rule, whatever interpreter path it names.
func TestInjectKeepsShebangFirst(t *testing.T) {
	const body = "echo hi\n"
	shebangs := []string{"#!/usr/bin/env node", "#!/usr/local/bin/bash", "#!/bin/sh"}
	exts := []string{".sh", ".bash", ".zsh", ".py", ".js", ".mjs", ".cjs", ".ts", ".mts", ".cts"}

	for _, ext := range exts {
		for i, sb := range shebangs {
			name := "script" + ext
			t.Run(ext+"/"+sb, func(t *testing.T) {
				got, modified := runProcess(t, internal.LicenseInject, "apache2", map[string]string{name: sb + "\n" + body})
				assert.Equal(t, 1, modified, "case %d", i)
				want := sb + "\n" + header(t, name, "apache2") + body
				assert.Equal(t, want, got[name])
			})
		}
	}
}

// TestInjectShebangMatchesV100Header pins the exact bytes against what the
// v1.0.0 binary wrote for the same file without a shebang: the header text is
// unchanged and only moves below the shebang.
func TestInjectShebangMatchesV100Header(t *testing.T) {
	stamped := stampedByV100["apache2"]
	cases := []struct{ name, shebang, plain, body string }{
		{"local.sh", "#!/usr/local/bin/bash", "plain.sh", "echo hi\n"},
		{"cli.js", "#!/usr/bin/env node", "a.js", "const a = 1;\n"},
		{"cli.ts", "#!/usr/bin/env node", "a.ts", "export const a = 1;\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := runProcess(t, internal.LicenseInject, "apache2", map[string]string{c.name: c.shebang + "\n" + c.body})
			assert.Equal(t, c.shebang+"\n"+stamped[c.plain], got[c.name])
		})
	}
}

// TestShebangOnlyCountsOnLineOne makes sure a "#!" line further down (a
// heredoc, a template literal) is not taken for the file's shebang.
func TestShebangOnlyCountsOnLineOne(t *testing.T) {
	files := map[string]string{
		"gen.sh": "cat <<EOF > out.sh\n#!/bin/sh\necho inner\nEOF\n",
		"gen.ts": "export const script = `\n#!/usr/bin/env node\nconsole.log(1)\n`;\n",
	}
	got, _ := runProcess(t, internal.LicenseInject, "apache2", files)
	for name, body := range files {
		assert.Equal(t, header(t, name, "apache2")+body, got[name], name)
	}
}

// TestInjectKeepsXMLDeclarationFirst checks the XML declaration stays first.
func TestInjectKeepsXMLDeclarationFirst(t *testing.T) {
	const decl = `<?xml version="1.0" encoding="UTF-8"?>`
	files := map[string]string{
		"pom.xml":   decl + "\n<project/>\n",
		"icon.svg":  decl + "\n<svg xmlns=\"http://www.w3.org/2000/svg\"/>\n",
		"plain.svg": "<svg xmlns=\"http://www.w3.org/2000/svg\"/>\n",
	}
	got, modified := runProcess(t, internal.LicenseInject, "apache2", files)
	assert.Equal(t, 3, modified)

	assert.Equal(t, decl+"\n"+header(t, "pom.xml", "apache2")+"<project/>\n", got["pom.xml"])
	assert.Equal(t, decl+"\n"+header(t, "icon.svg", "apache2")+"<svg xmlns=\"http://www.w3.org/2000/svg\"/>\n", got["icon.svg"])
	assert.Equal(t, header(t, "plain.svg", "apache2")+files["plain.svg"], got["plain.svg"])
	assert.True(t, strings.HasPrefix(got["pom.xml"], "<?xml"), "the declaration must stay first")
	// Same header bytes v1.0.0 used for .xml, now below the declaration.
	assert.Equal(t, decl+"\n"+stampedByV100["apache2"]["plain.xml"], strings.Replace(got["pom.xml"], "<project/>", "<a/>", 1))
}

// TestInjectKeepsDockerfileDirectivesFirst checks BuildKit parser directives
// stay above the header, in any order, for Dockerfiles and Containerfiles at
// any depth.
func TestInjectKeepsDockerfileDirectivesFirst(t *testing.T) {
	const rest = "FROM alpine\nRUN true\n"
	directives := map[string]string{
		"syntax":   "# syntax=docker/dockerfile:1\n",
		"labs":     "# syntax=docker.io/docker/dockerfile:1.7-labs\n",
		"all":      "# syntax=docker/dockerfile:1\n# escape=`\n# check=skip=JSONArgsRecommended;error=true\n",
		"reversed": "# check=error=true\n# escape=\\\n# syntax=docker/dockerfile:1\n",
		"none":     "",
	}
	paths := []string{"Dockerfile", "Dockerfile.dev", "build/Dockerfile", "images/app/Dockerfile.prod", "Containerfile", "deploy/Containerfile"}

	for label, d := range directives {
		for _, p := range paths {
			t.Run(label+"/"+p, func(t *testing.T) {
				got, modified := runProcess(t, internal.LicenseInject, "apache2", map[string]string{p: d + rest})
				assert.Equal(t, 1, modified)
				want := header(t, p, "apache2") + rest
				if d != "" {
					want = d + header(t, p, "apache2") + rest
				}
				assert.Equal(t, want, got[p])
			})
		}
	}
}

// TestNestedBuildFilesMatch checks Dockerfiles, Containerfiles and Makefiles
// match at any depth and that the root patterns still render the same "# "
// header v1.0.0 used.
func TestNestedBuildFilesMatch(t *testing.T) {
	cfg := embeddedConfig(t)
	for _, p := range []string{
		"Dockerfile", "Dockerfile.dev", "build/Dockerfile", "a/b/c/Dockerfile.ci",
		"Containerfile", "deploy/Containerfile",
		"Makefile", "sub/Makefile", "a/b/Makefile", "mk/rules.mk", "common.mk",
	} {
		t.Run(p, func(t *testing.T) {
			rule := getRule(cfg, p)
			require.NotEmpty(t, rule, "%s should match a rule", p)
			assert.Equal(t, "# ", cfg.Golic.Rules[rule].Prefix)
			assert.Empty(t, cfg.Golic.Rules[rule].Suffix)
		})
	}

	got, modified := runProcess(t, internal.LicenseInject, "apache2", map[string]string{
		"build/Dockerfile": "FROM alpine\n",
		"sub/Makefile":     "all:\n\ttrue\n",
	})
	assert.Equal(t, 2, modified)
	assert.Equal(t, stampedByV100["apache2"]["Dockerfile"], got["build/Dockerfile"])
	assert.Equal(t, stampedByV100["apache2"]["Makefile"], got["sub/Makefile"])
}

// TestBuiltinRulesForNewFileTypes checks each added file type renders with the
// expected comment style, header first.
func TestBuiltinRulesForNewFileTypes(t *testing.T) {
	cases := []struct {
		name, prefix, suffix string
	}{
		{"app.py", "# ", ""},
		{"api.proto", "/*", "*/"},
		{"schema.sql", "-- ", ""},
		{"site.scss", "/*", "*/"},
		{"site.less", "/*", "*/"},
		{"Cargo.toml", "# ", ""},
		{"main.hcl", "# ", ""},
		{"prod.tfvars", "# ", ""},
		{"schema.graphql", "# ", ""},
		{"schema.graphqls", "# ", ""},
		{".dockerignore", "# ", ""},
		{"chart/.helmignore", "# ", ""},
		{"run.bash", "# ", ""},
		{"env.zsh", "# ", ""},
		{"lib.rs", "/*", "*/"},
		{"main.c", "/*", "*/"},
		{"main.h", "/*", "*/"},
		{"main.cpp", "/*", "*/"},
		{"Main.kt", "/*", "*/"},
		{"App.swift", "/*", "*/"},
	}
	cfg := embeddedConfig(t)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rule := getRule(cfg, c.name)
			require.NotEmpty(t, rule, "%s should have a built-in rule", c.name)
			assert.Equal(t, c.prefix, cfg.Golic.Rules[rule].Prefix)
			assert.Equal(t, c.suffix, cfg.Golic.Rules[rule].Suffix)

			const body = "body\n"
			got, modified := runProcess(t, internal.LicenseInject, "mit", map[string]string{c.name: body})
			assert.Equal(t, 1, modified)
			h := header(t, c.name, "mit")
			assert.Equal(t, h+body, got[c.name])
			if c.suffix != "" {
				assert.True(t, strings.HasPrefix(h, c.prefix+"\nMIT License\n"), "got %q", h)
				assert.True(t, strings.HasSuffix(h, c.suffix+"\n"), "got %q", h)
			} else {
				assert.True(t, strings.HasPrefix(h, c.prefix+"MIT License\n"), "got %q", h)
			}
		})
	}
}

// TestDataAndDocFormatsStayUnmatched guards the deliberate gaps: JSON has no
// comments and Markdown/MDX headers would show up in rendered docs.
func TestDataAndDocFormatsStayUnmatched(t *testing.T) {
	cfg := embeddedConfig(t)
	for _, p := range []string{"package.json", "README.md", "docs/page.mdx"} {
		assert.Empty(t, getRule(cfg, p), "%s should not match a built-in rule", p)
	}
}

// TestNoChurnOnFilesStampedByV100 replays files the v1.0.0 binary stamped,
// one per built-in rule, including the ones it placed wrongly (above a
// shebang, an XML declaration or a Dockerfile directive). The current rules
// must see every one as done: no file is rewritten and no header is doubled.
func TestNoChurnOnFilesStampedByV100(t *testing.T) {
	cfg := embeddedConfig(t)
	for template, files := range stampedByV100 {
		for name := range files {
			require.NotEmpty(t, getRule(cfg, name), "%s lost its rule", name)
		}
		for _, kind := range []internal.LicenseCommandType{internal.LicenseInject, internal.LicenseReplace} {
			t.Run(template, func(t *testing.T) {
				got, modified := runProcess(t, kind, template, files)
				assert.Equal(t, 0, modified, "no file should change")
				for name, want := range files {
					assert.Equal(t, want, got[name], "%s was rewritten", name)
				}
			})
		}
	}
}

// TestMergeRulesFalseReplacesBuiltins documents mergeRules: false. The local
// rules become the whole ruleset, so a type the local file does not list has no
// rule and is skipped.
func TestMergeRulesFalseReplacesBuiltins(t *testing.T) {
	load := func(t *testing.T, local string) *Config {
		t.Helper()
		p := filepath.Join(t.TempDir(), ".golic.yaml")
		require.NoError(t, os.WriteFile(p, []byte(local), 0o600))
		u := &Process{Opts: internal.Options{ConfigPath: p}}
		base, err := u.readCommonConfig()
		require.NoError(t, err)
		u.cfgBase = base
		cfg, err := u.readLocalConfig()
		require.NoError(t, err)
		return cfg
	}

	t.Run("false with rules replaces the built-in set", func(t *testing.T) {
		cfg := load(t, "golic:\n  mergeRules: false\n  rules:\n    .go:\n      prefix: \"//\"\n")
		assert.Len(t, cfg.Golic.Rules, 1)
		assert.Equal(t, ".go", getRule(cfg, "main.go"))
		assert.Empty(t, getRule(cfg, "run.sh"), ".sh is not listed, so it has no rule")
		assert.NotEmpty(t, cfg.Golic.Licenses["apache2"], "licenses still merge")
	})

	t.Run("false without rules keeps the built-in set", func(t *testing.T) {
		cfg := load(t, "golic:\n  mergeRules: false\n")
		assert.Equal(t, len(embeddedConfig(t).Golic.Rules), len(cfg.Golic.Rules))
	})

	t.Run("default merges and a local key replaces the whole built-in rule", func(t *testing.T) {
		cfg := load(t, "golic:\n  rules:\n    .sh:\n      prefix: \"## \"\n")
		assert.Equal(t, "## ", cfg.Golic.Rules[".sh"].Prefix)
		assert.Empty(t, cfg.Golic.Rules[".sh"].Under, "the built-in under list is not carried over")
		assert.Equal(t, "/*", cfg.Golic.Rules[".ts"].Prefix)
	})
}

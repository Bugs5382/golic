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

const (
	oldOwner = "2026 Old Owner"
	newOwner = "2026 New Owner"
)

// runWith is runProcess with a chosen copyright and an optional local config.
func runWith(t *testing.T, kind internal.LicenseCommandType, template, copyright, localConfig string, files map[string]string) (map[string]string, int) {
	t.Helper()

	ignore := filepath.Join(t.TempDir(), ".licignore")
	require.NoError(t, os.WriteFile(ignore, nil, 0o600))

	cfgPath := ""
	if localConfig != "" {
		cfgPath = filepath.Join(t.TempDir(), ".golic.yaml")
		require.NoError(t, os.WriteFile(cfgPath, []byte(localConfig), 0o600))
	}

	dir := t.TempDir()
	for name, content := range files {
		p := filepath.Join(dir, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o750))
		require.NoError(t, os.WriteFile(p, []byte(content), 0o600))
	}
	t.Chdir(dir)

	u := ProcessFile(context.Background(), internal.Options{
		Type:       kind,
		Template:   template,
		Copyright:  copyright,
		LicIgnore:  ignore,
		ConfigPath: cfgPath,
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

// render is the header golic writes for path with the embedded rules.
func render(t *testing.T, path, template, copyright string) string {
	t.Helper()
	h, err := getCommentedLicense(embeddedConfig(t), internal.Options{Template: template, Copyright: copyright}, path)
	require.NoError(t, err)
	return h
}

// TestReplaceSwapsOwnHeaderInPlace checks that replace changes only the
// header golic wrote and keeps every byte around it: comments right below the
// header, the blank line (or its absence) after it, a shebang, an XML
// declaration and a Go package clause.
func TestReplaceSwapsOwnHeaderInPlace(t *testing.T) {
	type layout struct{ name, before, after string }
	cases := []layout{
		{"values.yaml", "", "# Default values for the chart.\n# Override anything below.\nkey: value\n"},
		{"blank.yaml", "", "\n# A section comment.\nkey: value\n"},
		{"doc.go", "package doc\n", "\n// Package doc is documented here.\n"},
		{"tight.go", "package tight\n", "func F() {}\n"},
		{"spaced.go", "package spaced\n", "\n\nimport \"fmt\"\n\nvar _ = fmt.Sprint\n"},
		{"App.tsx", "", "import React from \"react\";\n"},
		{"spaced.ts", "", "\nimport { a } from \"./a\";\n"},
		{"run.sh", "#!/usr/bin/env bash\n", "# Real comment about the script.\necho hi\n"},
		{"doc.xml", "<?xml version=\"1.0\"?>\n", "<!-- keep me -->\n<root/>\n"},
		{"_helpers.tpl", "", "{{/* Expand the name of the chart. */}}\n"},
		{"q.sql", "", "-- migration note\nSELECT 1;\n"},
	}

	files := map[string]string{}
	for _, c := range cases {
		files[c.name] = c.before + render(t, c.name, "apache2", oldOwner) + c.after
	}

	got, modified := runWith(t, internal.LicenseReplace, "apache2", newOwner, "", files)
	assert.Equal(t, len(cases), modified)
	for _, c := range cases {
		want := c.before + render(t, c.name, "apache2", newOwner) + c.after
		assert.Equal(t, want, got[c.name], c.name)
	}

	t.Run("second run changes nothing", func(t *testing.T) {
		again, modified := runWith(t, internal.LicenseReplace, "apache2", newOwner, "", got)
		assert.Equal(t, 0, modified)
		assert.Equal(t, got, again)

		_, modified = runWith(t, internal.LicenseInject, "apache2", newOwner, "", got)
		assert.Equal(t, 0, modified, "inject must see every file as stamped")
	})
}

// TestReplaceSwitchesLicenseType checks replace still moves a file from one
// configured license to another, without touching the comment below it.
func TestReplaceSwitchesLicenseType(t *testing.T) {
	body := "# Keep this comment.\nkey: value\n"
	files := map[string]string{"a.yaml": render(t, "a.yaml", "mit", oldOwner) + body}

	got, modified := runWith(t, internal.LicenseReplace, "apache2", newOwner, "", files)
	assert.Equal(t, 1, modified)
	assert.Equal(t, render(t, "a.yaml", "apache2", newOwner)+body, got["a.yaml"])
}

// TestReplaceFindsLocalLicense checks a license defined only in the local
// .golic.yaml is recognised as golic's own header too.
func TestReplaceFindsLocalLicense(t *testing.T) {
	local := "golic:\n  licenses:\n    acme: |\n      Copyright {{copyright}}\n\n      Proprietary to Acme.\n"
	body := "// Package a.\n"
	old := "package a\n\n/*\nCopyright " + oldOwner + "\n\nProprietary to Acme.\n*/\n"
	files := map[string]string{"a.go": old + body}

	got, modified := runWith(t, internal.LicenseReplace, "apache2", newOwner, local, files)
	assert.Equal(t, 1, modified)
	assert.Equal(t, "package a\n"+render(t, "a.go", "apache2", newOwner)+body, got["a.go"])
}

// TestReplaceLeavesOtherHeadersAlone checks a leading comment golic did not
// write is never stripped: a near miss of a golic header, another license
// header, and a plain comment. A file with a foreign license header is left
// as it is; a file with only an ordinary comment gets a header like inject.
func TestReplaceLeavesOtherHeadersAlone(t *testing.T) {
	nearMiss := strings.Replace(render(t, "near.yaml", "apache2", oldOwner), "WITHOUT WARRANTIES", "WITHOUT ANY WARRANTIES", 1) + "key: value\n"
	foreign := "package f\n\n/*\nCopyright 2020 Someone Else. All rights reserved.\nSPDX-License-Identifier: BSD-3-Clause\n*/\n\nfunc F() {}\n"
	files := map[string]string{
		"near.yaml": nearMiss,
		"f.go":      foreign,
	}

	got, modified := runWith(t, internal.LicenseReplace, "apache2", newOwner, "", files)
	assert.Equal(t, 0, modified)
	assert.Equal(t, nearMiss, got["near.yaml"])
	assert.Equal(t, foreign, got["f.go"])

	t.Run("ordinary comment gets a header like inject", func(t *testing.T) {
		plain := map[string]string{"plain.yaml": "# Default values for the chart.\nkey: value\n"}
		replaced, modified := runWith(t, internal.LicenseReplace, "apache2", newOwner, "", plain)
		assert.Equal(t, 1, modified)
		injected, _ := runWith(t, internal.LicenseInject, "apache2", newOwner, "", plain)
		assert.Equal(t, injected["plain.yaml"], replaced["plain.yaml"])
	})
}

// TestReplaceInjectsWhenNoHeader checks a file without any header comes out
// exactly as inject would write it.
func TestReplaceInjectsWhenNoHeader(t *testing.T) {
	files := map[string]string{
		"a.go":   "package a\n\nfunc A() {}\n",
		"run.sh": "#!/bin/sh\necho hi\n",
		"a.ts":   "export const a = 1;\n",
	}
	replaced, modified := runWith(t, internal.LicenseReplace, "apache2", newOwner, "", files)
	assert.Equal(t, len(files), modified)
	injected, _ := runWith(t, internal.LicenseInject, "apache2", newOwner, "", files)
	assert.Equal(t, injected, replaced)
}

// TestReplaceDryRun checks --dry reports the change without writing it.
func TestReplaceDryRun(t *testing.T) {
	files := map[string]string{"a.yaml": render(t, "a.yaml", "apache2", oldOwner) + "# note\nk: v\n"}

	ignore := filepath.Join(t.TempDir(), ".licignore")
	require.NoError(t, os.WriteFile(ignore, nil, 0o600))
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.yaml"), []byte(files["a.yaml"]), 0o600))
	t.Chdir(dir)

	u := ProcessFile(context.Background(), internal.Options{
		Type: internal.LicenseReplace, Template: "apache2", Copyright: newOwner,
		LicIgnore: ignore, Dry: true, ModifiedExitStatus: true,
	})
	require.NoError(t, u.Run())
	require.Error(t, u.Changes())
	assert.Equal(t, 1, u.modified)
	b, rerr := os.ReadFile("a.yaml")
	require.NoError(t, rerr)
	assert.Equal(t, files["a.yaml"], string(b))
}

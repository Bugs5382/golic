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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bugs5382/golic/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuiltinTemplatesAreOfficial pins every built-in license to the official
// text kept in testdata/official, byte for byte. default.golic.yaml names the
// source of each one.
func TestBuiltinTemplatesAreOfficial(t *testing.T) {
	cfg := embeddedConfig(t)
	files, err := filepath.Glob(filepath.Join("testdata", "official", "*.txt"))
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, f := range files {
		key := strings.TrimSuffix(filepath.Base(f), ".txt")
		want, err := os.ReadFile(f)
		require.NoError(t, err)
		assert.Equal(t, string(want), cfg.Golic.Licenses[key], key)
	}
}

// renderText renders an arbitrary license text for path with the embedded rules.
func renderText(t *testing.T, path, text, copyright string) string {
	t.Helper()
	cfg := embeddedConfig(t)
	cfg.Golic.Licenses["under-test"] = text
	h, err := getCommentedLicense(cfg, internal.Options{Template: "under-test", Copyright: copyright}, path)
	require.NoError(t, err)
	return h
}

// TestSupersededTextsMatchV100Output checks the recorded earlier texts are the
// ones the v1.0.0 binary really wrote.
func TestSupersededTextsMatchV100Output(t *testing.T) {
	for template, files := range stampedByV100 {
		require.NotEmpty(t, supersededLicenses[template], template)
		old := supersededLicenses[template][0]
		for name, stamped := range files {
			assert.Contains(t, stamped, renderText(t, name, old, testCopyright), "%s/%s", template, name)
		}
	}
}

// TestUpgradeFromV100Stamps replays the v1.0.0 fixtures through the current
// binary: inject leaves them alone, replace moves each to the current text in
// place, and after that both inject and replace report nothing to do.
func TestUpgradeFromV100Stamps(t *testing.T) {
	for template, files := range stampedByV100 {
		t.Run(template, func(t *testing.T) {
			injected, modified := runProcess(t, internal.LicenseInject, template, files)
			assert.Equal(t, 0, modified, "inject must not stamp a second header")
			assert.Equal(t, files, injected)

			want := map[string]string{}
			for name, stamped := range files {
				old := renderText(t, name, supersededLicenses[template][0], testCopyright)
				want[name] = strings.Replace(stamped, old, header(t, name, template), 1)
			}
			replaced, modified := runProcess(t, internal.LicenseReplace, template, files)
			assert.Equal(t, len(files), modified)
			assert.Equal(t, want, replaced)

			_, modified = runProcess(t, internal.LicenseInject, template, replaced)
			assert.Equal(t, 0, modified, "inject after replace")
			_, modified = runProcess(t, internal.LicenseReplace, template, replaced)
			assert.Equal(t, 0, modified, "second replace")
		})
	}
}

// upgradeSamples is one file per comment style, with content around the
// header that must survive an upgrade untouched.
var upgradeSamples = []struct{ name, before, after string }{
	{"a.go", "package a\n", "\nfunc A() {}\n"},
	{"a.ts", "", "import { b } from \"./b\";\n"},
	{"values.yaml", "", "# Default values.\nkey: value\n"},
	{"run.sh", "#!/usr/bin/env bash\n", "# about the script\necho hi\n"},
	{"a.xml", "<?xml version=\"1.0\"?>\n", "<root/>\n"},
	{"a.html", "", "<p>hi</p>\n"},
	{"_helpers.tpl", "", "{{- define \"x\" -}}{{- end -}}\n"},
	{"q.sql", "", "SELECT 1;\n"},
}

// TestUpgradeEveryChangedTemplate stamps each sample with every earlier text
// of every built-in license and checks inject, remove and replace handle it.
func TestUpgradeEveryChangedTemplate(t *testing.T) {
	for template, olds := range supersededLicenses {
		for i, old := range olds {
			t.Run(template, func(t *testing.T) {
				stamped := map[string]string{}
				migrated := map[string]string{}
				bare := map[string]string{}
				for _, s := range upgradeSamples {
					stamped[s.name] = s.before + renderText(t, s.name, old, testCopyright) + s.after
					migrated[s.name] = s.before + header(t, s.name, template) + s.after
					bare[s.name] = s.before + s.after
				}

				got, modified := runProcess(t, internal.LicenseInject, template, stamped)
				assert.Equal(t, 0, modified, "inject, text %d", i)
				assert.Equal(t, stamped, got)

				got, modified = runProcess(t, internal.LicenseReplace, template, stamped)
				assert.Equal(t, len(stamped), modified, "replace, text %d", i)
				assert.Equal(t, migrated, got)

				_, modified = runProcess(t, internal.LicenseReplace, template, got)
				assert.Equal(t, 0, modified, "second replace, text %d", i)

				got, modified = runProcess(t, internal.LicenseRemove, template, stamped)
				assert.Equal(t, len(stamped), modified, "remove, text %d", i)
				assert.Equal(t, bare, got)
			})
		}
	}
}

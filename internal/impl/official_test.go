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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	golog "github.com/Bugs5382/go-log"
	"github.com/Bugs5382/golic/internal"
	"github.com/Bugs5382/golic/internal/logging"
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

// captureLogs sends golic's log output to a buffer as JSON lines, at debug
// level when debug is set and at info otherwise.
func captureLogs(t *testing.T, debug bool) *bytes.Buffer {
	t.Helper()
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("LOG_FORMAT", "")
	level := golog.LevelInfo
	if debug {
		level = golog.LevelDebug
	}
	buf := new(bytes.Buffer)
	logging.Use(golog.NewLoggerWithOptions("golic",
		golog.WithOutput(buf),
		golog.WithDefaultFormat(golog.FormatJSON),
		golog.WithDefaultLevel(level),
	), debug)
	t.Cleanup(func() { logging.Init(false) })
	return buf
}

// logLines returns the captured lines at level whose message contains text.
func logLines(t *testing.T, buf *bytes.Buffer, level, text string) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, raw := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if raw == "" {
			continue
		}
		var line map[string]any
		require.NoError(t, json.Unmarshal([]byte(raw), &line), raw)
		if line["level"] == level && strings.Contains(fmt.Sprint(line["message"]), text) {
			out = append(out, line)
		}
	}
	return out
}

// TestInjectSummarisesEarlierTexts checks inject reports files with an earlier
// text in one warning at the end of the run, naming the replace command for
// the run's template and copyright, and keeps the per-file detail at debug.
func TestInjectSummarisesEarlierTexts(t *testing.T) {
	const template = "apache2"
	old := supersededLicenses[template][0]
	stamped := map[string]string{}
	for _, s := range upgradeSamples {
		stamped[s.name] = s.before + renderText(t, s.name, old, testCopyright) + s.after
	}
	summary := fmt.Sprintf("%d file(s) carry an earlier golic text of the %s license; run golic replace -t %s -c \"%s\" to update them",
		len(stamped), template, template, testCopyright)
	detail := "file carries an earlier golic text of this license"

	t.Run("one summary line at warn", func(t *testing.T) {
		buf := captureLogs(t, false)
		_, modified := runProcess(t, internal.LicenseInject, template, stamped)
		assert.Equal(t, 0, modified)
		lines := logLines(t, buf, "warn", "earlier golic text")
		require.Len(t, lines, 1, buf.String())
		assert.Equal(t, summary, lines[0]["message"])
		assert.Empty(t, logLines(t, buf, "debug", detail), "per-file detail must stay below info")
	})

	t.Run("per-file detail at debug", func(t *testing.T) {
		buf := captureLogs(t, true)
		runProcess(t, internal.LicenseInject, template, stamped)
		lines := logLines(t, buf, "debug", detail)
		require.Len(t, lines, len(stamped), buf.String())
		for _, l := range lines {
			assert.Equal(t, template, l["template"])
			assert.Contains(t, stamped, l["path"])
		}
		assert.Empty(t, logLines(t, buf, "warn", detail), "per-file detail must not be a warning")
		assert.Len(t, logLines(t, buf, "warn", summary), 1)
	})

	t.Run("dry run with -x still counts", func(t *testing.T) {
		buf := captureLogs(t, false)
		ignore := filepath.Join(t.TempDir(), ".licignore")
		require.NoError(t, os.WriteFile(ignore, nil, 0o600))
		dir := t.TempDir()
		for name, content := range stamped {
			require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
		}
		t.Chdir(dir)
		u := ProcessFile(context.Background(), internal.Options{
			Type: internal.LicenseInject, Template: template, Copyright: testCopyright,
			LicIgnore: ignore, Dry: true, ModifiedExitStatus: true,
		})
		require.NoError(t, u.Run())
		assert.NoError(t, u.Changes(), "earlier texts must not fail CI")
		assert.Len(t, logLines(t, buf, "warn", summary), 1)
	})

	t.Run("no summary when every file is current", func(t *testing.T) {
		current := map[string]string{}
		for _, s := range upgradeSamples {
			current[s.name] = s.before + header(t, s.name, template) + s.after
		}
		buf := captureLogs(t, true)
		runProcess(t, internal.LicenseInject, template, current)
		assert.Empty(t, logLines(t, buf, "warn", "earlier golic text"))
		assert.Empty(t, logLines(t, buf, "debug", detail))
	})
}

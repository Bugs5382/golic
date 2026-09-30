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
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bugs5382/golic/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// officialSHA256 pins each embedded LICENSE text to the file downloaded from
// its source (licenseFileSources in licensefile.go), so an edit shows up.
var officialSHA256 = map[string]string{
	"apache2":   "cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30",
	"gpl2":      "edaef632cbb643e4e7a221717a6c441a4c1a7c918e6e4d56debc3d8739b233f6",
	"gpl3":      "3972dc9744f6499f0f9b2dbf76696f2ae7ad8af9b23dde66d6af86c9dfb36986",
	"lgpl3":     "e3a994d82e644b03a792a930f574002658412f62407f5fee083f2555c5f23118",
	"agpl3":     "0d96a4ff68ad6d4b6f1f30f713b18d5184912ba8dd389f86aa7710db079abcb0",
	"mpl2":      "3f3d9e0024b1921b067d6f7f88deb4a60cbe7a78e76c64e3f1d7fc3b779b9d04",
	"epl2":      "0becf16567beb77fa252b7664631dd177c8f9a1889e48995b45379c7130e5303",
	"unlicense": "b5065838cbac452dfc855ba6e6e031481ad2c68406f70d21ead9321374653e6c",
	"mit":       "b05785f9f18e6716bab63424b11454513b9943a222595b70411009202fc592b5",
	"bsd2":      "f32fb3b417a194167cfad068223fc975ba96c5960513a10f66a3c28720aec1df",
	"bsd3":      "5a93d5831e1297ab10fe643e1a631e83be392896da14ee2951285a79012df69d",
	"isc":       "521c6f0ed8e64736684f89843d16a4df6b4b7449acbcfc6ca6630a3037ca8c53",
}

type lfRun struct {
	template, licenseFile, localConfig string
	dry, modifiedExit                  bool
	files                              map[string]string
}

// runLicenseFile runs one inject pass with the LICENSE options and returns
// the tree afterwards, the change count and the -x result.
func runLicenseFile(t *testing.T, r lfRun) (map[string]string, int, error) {
	t.Helper()
	ignore := filepath.Join(t.TempDir(), ".licignore")
	require.NoError(t, os.WriteFile(ignore, nil, 0o600))
	cfgPath := ""
	if r.localConfig != "" {
		cfgPath = filepath.Join(t.TempDir(), ".golic.yaml")
		require.NoError(t, os.WriteFile(cfgPath, []byte(r.localConfig), 0o600))
	}
	dir := t.TempDir()
	for name, content := range r.files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
	}
	t.Chdir(dir)

	u := ProcessFile(context.Background(), internal.Options{
		Type: internal.LicenseInject, Template: r.template, Copyright: testCopyright,
		LicIgnore: ignore, ConfigPath: cfgPath, LicenseFile: r.licenseFile,
		Dry: r.dry, ModifiedExitStatus: r.modifiedExit,
	})
	if err := u.Run(); err != nil {
		return nil, 0, err
	}
	out := map[string]string{}
	entries, err := os.ReadDir(".")
	require.NoError(t, err)
	for _, e := range entries {
		b, err := os.ReadFile(e.Name())
		require.NoError(t, err)
		out[e.Name()] = string(b)
	}
	return out, u.modified, u.Changes()
}

func TestEmbeddedLicenseTextsAreOfficial(t *testing.T) {
	for key, want := range officialSHA256 {
		raw, err := licenseTexts.ReadFile("licenses/" + key + ".txt")
		require.NoError(t, err, key)
		sum := sha256.Sum256(raw)
		assert.Equal(t, want, hex.EncodeToString(sum[:]), key)
	}
	assert.Len(t, officialSHA256, len(licenseFileSources), "every built-in license with an official text is pinned")
}

// TestLicenseFileWrittenForEachTemplate checks the LICENSE is the official text
// with the copyright filled only where that text has a copyright line.
func TestLicenseFileWrittenForEachTemplate(t *testing.T) {
	filled := map[string]string{
		"mit":  "Copyright (c) " + testCopyright + "\n",
		"bsd2": "Copyright (c) " + testCopyright + "\n",
		"bsd3": "Copyright (c) " + testCopyright + "\n",
		"isc":  "Copyright " + testCopyright + "\n",
	}
	for key := range officialSHA256 {
		t.Run(key, func(t *testing.T) {
			raw, err := licenseTexts.ReadFile("licenses/" + key + ".txt")
			require.NoError(t, err)

			got, modified, err := runLicenseFile(t, lfRun{template: key, licenseFile: "LICENSE"})
			require.NoError(t, err)
			assert.Equal(t, 1, modified)

			if line, ok := filled[key]; ok {
				assert.Contains(t, got["LICENSE"], line)
				assert.NotContains(t, got["LICENSE"], "<year>")
				assert.Equal(t, len(strings.Split(string(raw), "\n")), len(strings.Split(got["LICENSE"], "\n")), "only the copyright line changes")
			} else {
				assert.Equal(t, string(raw), got["LICENSE"], "verbatim, placeholders untouched")
			}
		})
	}

	t.Run("apache keeps the appendix placeholder", func(t *testing.T) {
		got, _, err := runLicenseFile(t, lfRun{template: "apache2", licenseFile: "LICENSE"})
		require.NoError(t, err)
		assert.Contains(t, got["LICENSE"], "Copyright [yyyy] [name of copyright owner]")
		assert.NotContains(t, got["LICENSE"], testCopyright)
	})
}

func TestLicenseFileDryAndModifiedExit(t *testing.T) {
	got, modified, changes := runLicenseFile(t, lfRun{template: "mit", licenseFile: "LICENSE", dry: true, modifiedExit: true})
	assert.Equal(t, 1, modified)
	assert.Error(t, changes, "a missing LICENSE is a change for -x")
	_, exists := got["LICENSE"]
	assert.False(t, exists, "dry run writes nothing")
}

func TestLicenseFileAlreadyCorrect(t *testing.T) {
	first, _, err := runLicenseFile(t, lfRun{template: "apache2", licenseFile: "LICENSE"})
	require.NoError(t, err)

	_, modified, changes := runLicenseFile(t, lfRun{template: "apache2", licenseFile: "LICENSE", modifiedExit: true, dry: true,
		files: map[string]string{"LICENSE": first["LICENSE"]}})
	assert.Equal(t, 0, modified)
	assert.NoError(t, changes)
}

func TestLicenseFileRewordedIsFixed(t *testing.T) {
	correct, _, err := runLicenseFile(t, lfRun{template: "apache2", licenseFile: "LICENSE"})
	require.NoError(t, err)
	reworded := strings.Replace(correct["LICENSE"], "Version 2.0, January 2004", "Version 2.0", 1)

	_, modified, changes := runLicenseFile(t, lfRun{template: "apache2", licenseFile: "LICENSE", dry: true, modifiedExit: true,
		files: map[string]string{"LICENSE": reworded}})
	assert.Equal(t, 1, modified)
	assert.Error(t, changes, "-x catches a reworded LICENSE")

	got, modified, err := runLicenseFile(t, lfRun{template: "apache2", licenseFile: "LICENSE",
		files: map[string]string{"LICENSE": reworded}})
	require.NoError(t, err)
	assert.Equal(t, 1, modified)
	assert.Equal(t, correct["LICENSE"], got["LICENSE"])
}

func TestLicenseFileOffByDefault(t *testing.T) {
	got, modified, err := runLicenseFile(t, lfRun{template: "mit"})
	require.NoError(t, err)
	assert.Equal(t, 0, modified)
	_, exists := got["LICENSE"]
	assert.False(t, exists)
}

func TestLicenseFileFromConfig(t *testing.T) {
	t.Run("true writes LICENSE", func(t *testing.T) {
		got, modified, err := runLicenseFile(t, lfRun{template: "mit", localConfig: "golic:\n  licenseFile: true\n"})
		require.NoError(t, err)
		assert.Equal(t, 1, modified)
		assert.Contains(t, got["LICENSE"], "MIT License")
	})
	t.Run("a path", func(t *testing.T) {
		got, _, err := runLicenseFile(t, lfRun{template: "mit", localConfig: "golic:\n  licenseFile: COPYING\n"})
		require.NoError(t, err)
		assert.Contains(t, got["COPYING"], "MIT License")
	})
	t.Run("false stays off", func(t *testing.T) {
		got, _, err := runLicenseFile(t, lfRun{template: "mit", localConfig: "golic:\n  licenseFile: false\n"})
		require.NoError(t, err)
		_, exists := got["LICENSE"]
		assert.False(t, exists)
	})
	t.Run("the flag wins", func(t *testing.T) {
		got, _, err := runLicenseFile(t, lfRun{template: "mit", licenseFile: "LICENSE.txt", localConfig: "golic:\n  licenseFile: COPYING\n"})
		require.NoError(t, err)
		_, copying := got["COPYING"]
		assert.False(t, copying)
		assert.Contains(t, got["LICENSE.txt"], "MIT License")
	})
}

func TestLicenseFileNeedsAnOfficialText(t *testing.T) {
	_, _, err := runLicenseFile(t, lfRun{template: "copyright", licenseFile: "LICENSE"})
	assert.ErrorContains(t, err, "no official LICENSE text")
}

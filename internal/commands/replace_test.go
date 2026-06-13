package commands

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
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Bugs5382/golic/internal"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplace(t *testing.T) {

	_ = os.Chdir(internal.GetProjectRoot())

	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Parallel()

	t.Run("replace -- template missing", func(t *testing.T) {
		cmd := RootCmd()

		b := new(bytes.Buffer)

		cmd.SetOut(b)
		cmd.SetErr(b)

		cmd.SetArgs([]string{
			"replace",
		})

		err := cmd.Execute()

		assert.ErrorContains(t, err, "license template not provided")
	})

	t.Run("replace -- custom config file not found", func(t *testing.T) {
		cmd := RootCmd()

		b := new(bytes.Buffer)

		cmd.SetOut(b)
		cmd.SetErr(b)

		cmd.SetArgs([]string{
			"replace",
			"-p",
			".golic-test.yaml",
			"-d", // safety
		})

		err := cmd.Execute()

		assert.ErrorContains(t, err, "custom config file not found: .golic-test.yaml")
	})

	t.Run("replace -- custom lic ignore not found", func(t *testing.T) {
		cmd := RootCmd()

		b := new(bytes.Buffer)

		cmd.SetOut(b)
		cmd.SetErr(b)

		cmd.SetArgs([]string{
			"replace",
			"-l",
			".licignoreNotFound",
			"-d", // safety
		})

		err := cmd.Execute()

		assert.ErrorContains(t, err, "custom ignore file not found: .licignoreNotFound")
	})

	t.Run("replace -- template mit (no error)", func(t *testing.T) {
		cmd := RootCmd()

		b := new(bytes.Buffer)

		cmd.SetOut(b)
		cmd.SetErr(b)

		cmd.SetArgs([]string{
			"replace",
			"-t",
			"mit",
			"-d", // safety
		})

		err := cmd.Execute()

		assert.NoError(t, err)
	})

}

const replaceLicIgnore = `*
!*.go
!*/
`

const replaceGolicConfig = `golic:
  licenses:
    mit: |
      MIT License

      Copyright (c) {{copyright}}
    apache2: |
      Apache License 2.0

      Copyright {{copyright}}
  rules:
    .go:
      prefix: "\n/*"
      suffix: "*/"
      under:
        - "package *"
`

// setupReplaceWorkspace builds an isolated working directory containing a config,
// a .licignore and a single Go source file with an existing MIT header, then
// chdirs into it. The previous working directory is restored on cleanup.
func setupReplaceWorkspace(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()

	require.NoError(t, os.WriteFile(filepath.Join(dir, ".golic.yaml"), []byte(replaceGolicConfig), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".licignore"), []byte(replaceLicIgnore), 0o600))

	src := "package sample\n\nfunc Sample() {}\n"
	require.NoError(t, os.WriteFile(filepath.Join(dir, "sample.go"), []byte(src), 0o600))

	wd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() {
		_ = os.Chdir(wd)
	})

	return dir
}

func runReplace(args ...string) error {
	cmd := RootCmd()
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.SetArgs(append([]string{"replace"}, args...))
	return cmd.Execute()
}

// TestReplaceFiles exercises the replace command against real files. It does not
// run in parallel because it changes the process working directory.
func TestReplaceFiles(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	t.Run("replace swaps an existing header for the configured one", func(t *testing.T) {
		dir := setupReplaceWorkspace(t)
		sample := filepath.Join(dir, "sample.go")

		// Inject an MIT header first so the file has an existing license to replace.
		require.NoError(t, runReplace("-t", "mit", "-c", "2025 Old Owner"))
		seeded, err := os.ReadFile(sample)
		require.NoError(t, err)
		require.Contains(t, string(seeded), "MIT License")

		// Replace MIT with Apache-2.0 and a new copyright string.
		require.NoError(t, runReplace("-t", "apache2", "-c", "2026 New Owner"))

		got, err := os.ReadFile(sample)
		require.NoError(t, err)
		assert.Contains(t, string(got), "Apache License 2.0")
		assert.Contains(t, string(got), "Copyright 2026 New Owner")
		assert.NotContains(t, string(got), "MIT License")
		assert.True(t, strings.HasSuffix(string(got), "func Sample() {}\n"))
	})

	t.Run("dry run makes no change", func(t *testing.T) {
		dir := setupReplaceWorkspace(t)
		sample := filepath.Join(dir, "sample.go")

		before, err := os.ReadFile(sample)
		require.NoError(t, err)

		require.NoError(t, runReplace("-t", "apache2", "-c", "2026 New Owner", "-d"))

		after, err := os.ReadFile(sample)
		require.NoError(t, err)
		assert.Equal(t, string(before), string(after))
	})

	t.Run("-x exits non-zero when a change is needed", func(t *testing.T) {
		setupReplaceWorkspace(t)

		err := runReplace("-t", "apache2", "-c", "2026 New Owner", "-d", "-x")
		assert.Error(t, err)
	})

	t.Run("no-op when the header is already correct", func(t *testing.T) {
		dir := setupReplaceWorkspace(t)
		sample := filepath.Join(dir, "sample.go")

		// Stamp the file with the target header.
		require.NoError(t, runReplace("-t", "apache2", "-c", "2026 New Owner"))
		before, err := os.ReadFile(sample)
		require.NoError(t, err)

		// Replaying the same replace must be a no-op and must not trip -x.
		require.NoError(t, runReplace("-t", "apache2", "-c", "2026 New Owner", "-x"))

		after, err := os.ReadFile(sample)
		require.NoError(t, err)
		assert.Equal(t, string(before), string(after))
	})
}

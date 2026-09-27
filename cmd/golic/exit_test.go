package main

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
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const exitLicIgnore = "*\n!*.go\n!*/\n"

// workspace creates a temp project with one Go file and moves into it. With
// licensed set, the file already carries the MIT header for "2026 Tester".
func workspace(t *testing.T, licensed bool) string {
	t.Helper()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, ".licignore"), []byte(exitLicIgnore), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a\n"), 0o600))

	t.Chdir(dir)

	if licensed {
		var out, errOut bytes.Buffer
		require.Equal(t, exitClean, run([]string{"inject", "-t", "mit", "-c", "2026 Tester"}, &out, &errOut), errOut.String())
	}
	return dir
}

func TestExitCodes(t *testing.T) {
	tests := []struct {
		name     string
		licensed bool
		args     []string
		want     int
		stderr   string
	}{
		{
			name:     "clean tree with -x exits 0",
			licensed: true,
			args:     []string{"inject", "-t", "mit", "-c", "2026 Tester", "--dry", "-x"},
			want:     exitClean,
		},
		{
			name: "missing header without -x exits 0",
			args: []string{"inject", "-t", "mit", "-c", "2026 Tester", "--dry"},
			want: exitClean,
		},
		{
			name:   "missing header with --dry -x exits 1",
			args:   []string{"inject", "-t", "mit", "-c", "2026 Tester", "--dry", "-x"},
			want:   exitChanges,
			stderr: "1 file(s) need a license header",
		},
		{
			name:   "inject that modifies a file with -x exits 1",
			args:   []string{"inject", "-t", "mit", "-c", "2026 Tester", "-x"},
			want:   exitChanges,
			stderr: "1 file(s) modified",
		},
		{
			name:     "replace dry run that needs a change exits 1",
			licensed: true,
			args:     []string{"replace", "-t", "mit", "-c", "2027 Tester", "--dry", "-x"},
			want:     exitChanges,
			stderr:   "1 file(s) need the license header replaced",
		},
		{
			name:     "remove dry run that would strip a header exits 1",
			licensed: true,
			args:     []string{"remove", "-t", "mit", "-c", "2026 Tester", "--dry", "-x"},
			want:     exitChanges,
			stderr:   "1 file(s) have a license header to remove",
		},
		{
			name:   "missing template exits 2",
			args:   []string{"inject", "--dry", "-x"},
			want:   exitError,
			stderr: "license template not provided",
		},
		{
			name:   "unknown license key exits 2",
			args:   []string{"inject", "-t", "nope", "--dry", "-x"},
			want:   exitError,
			stderr: "no license found for nope",
		},
		{
			name:   "unknown flag exits 2",
			args:   []string{"inject", "--nope"},
			want:   exitError,
			stderr: "unknown flag",
		},
		{
			name: "no command prints help and exits 0",
			args: []string{},
			want: exitClean,
		},
		{
			name:   "unknown command exits 2",
			args:   []string{"nope"},
			want:   exitError,
			stderr: "unknown command",
		},
		{
			name: "version exits 0",
			args: []string{"version"},
			want: exitClean,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workspace(t, tt.licensed)

			var out, errOut bytes.Buffer
			got := run(tt.args, &out, &errOut)

			assert.Equal(t, tt.want, got, "stderr: %s", errOut.String())
			if tt.stderr != "" {
				assert.Contains(t, errOut.String(), tt.stderr)
			}
			assert.NotContains(t, errOut.String(), "something went wrong")
		})
	}
}

// TestExitCodesBinary checks the real process exit status, since that is what
// CI and the hub's license job see.
func TestExitCodesBinary(t *testing.T) {
	if testing.Short() {
		t.Skip("builds the binary")
	}

	bin := buildBinary(t, "")

	cases := []struct {
		name string
		args []string
		want int
	}{
		{"needs headers", []string{"inject", "-t", "mit", "-c", "2026 Tester", "--dry", "-x"}, exitChanges},
		{"real error", []string{"inject", "-t", "nope", "--dry", "-x"}, exitError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			workspace(t, false)

			cmd := exec.Command(bin, tc.args...) // #nosec G204 -- binary built by the test
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err := cmd.Run()

			var exitErr *exec.ExitError
			require.True(t, errors.As(err, &exitErr), "expected a non-zero exit, got %v", err)
			assert.Equal(t, tc.want, exitErr.ExitCode(), stderr.String())
			wantErrors := 0
			if tc.want == exitError {
				wantErrors = 1
			}
			assert.Equal(t, wantErrors, strings.Count(stderr.String(), "Error:"), "stderr: %s", stderr.String())
		})
	}

	t.Run("clean", func(t *testing.T) {
		workspace(t, true)
		cmd := exec.Command(bin, "inject", "-t", "mit", "-c", "2026 Tester", "--dry", "-x") // #nosec G204 -- binary built by the test
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, string(out))
	})
}

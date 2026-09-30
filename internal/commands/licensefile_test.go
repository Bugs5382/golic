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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLicenseFileFlag checks --license-file on its own names LICENSE, takes a
// path with =, and exists on inject and replace but not remove.
func TestLicenseFileFlag(t *testing.T) {
	run := func(t *testing.T, args ...string) error {
		t.Helper()
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, ".licignore"), nil, 0o600))
		t.Chdir(dir)
		cmd := RootCmd()
		b := new(bytes.Buffer)
		cmd.SetOut(b)
		cmd.SetErr(b)
		cmd.SetArgs(args)
		return cmd.Execute()
	}

	for _, command := range []string{"inject", "replace"} {
		t.Run(command+" default path", func(t *testing.T) {
			require.NoError(t, run(t, command, "-t", "mit", "-c", "2026 Someone", "--license-file"))
			b, err := os.ReadFile("LICENSE")
			require.NoError(t, err)
			assert.Contains(t, string(b), "Copyright (c) 2026 Someone")
		})
		t.Run(command+" custom path", func(t *testing.T) {
			require.NoError(t, run(t, command, "-t", "apache2", "--license-file=COPYING"))
			_, err := os.Stat("COPYING")
			assert.NoError(t, err)
			_, err = os.Stat("LICENSE")
			assert.True(t, os.IsNotExist(err))
		})
	}

	t.Run("remove has no such flag", func(t *testing.T) {
		assert.ErrorContains(t, run(t, "remove", "-t", "mit", "--license-file"), "unknown flag")
	})
}

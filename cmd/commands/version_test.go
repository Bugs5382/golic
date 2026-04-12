package commands

/*
Apache License 2.0

Copyright 2006 Shane

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
	"testing"

	"github.com/AbsaOSS/golic/internal/buildinfo"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	root, out := SetupTest()

	t.Run("should display the correct semantic version", func(t *testing.T) {
		out.Reset() // Clear buffer for this specific subtest
		root.SetArgs([]string{"version"})

		err := root.Execute()

		assert.NoError(t, err)
		assert.Contains(t, out.String(), buildinfo.Version)
	})

	t.Run("should fail if any flags are passed", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"version", "--invalid-flag"})

		err := root.Execute()

		assert.Error(t, err)
	})
}

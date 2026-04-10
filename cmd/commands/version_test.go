package commands

import (
	"testing"

	"github.com/AbsaOSS/golic/helpers"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	root, out := SetupTest()

	t.Run("should display the correct semantic version", func(t *testing.T) {
		out.Reset() // Clear buffer for this specific subtest
		root.SetArgs([]string{"version"})

		err := root.Execute()

		assert.NoError(t, err)
		assert.Contains(t, out.String(), helpers.Version)
	})

	t.Run("should fail if any flags are passed", func(t *testing.T) {
		out.Reset()
		root.SetArgs([]string{"version", "--invalid-flag"})

		err := root.Execute()

		assert.Error(t, err)
	})
}
